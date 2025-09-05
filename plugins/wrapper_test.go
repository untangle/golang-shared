package plugins

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

type wrapperTest struct {
	WrapperHelper
	matches bool
}

// decorator is a test double for a plugin that 'decorates' another
// plugin. Decoration plugins are used in policy manager to inject the
// proper settings.
type decorator struct {
	decorated         map[string]SettingsInjectablePlugin
	newPluginCallback PluginGeneratorCallback
}

func (d *decorator) Startup() error {
	return nil
}

// SettingsKey returns the settings key
func (d *decorator) SettingsKey() string {
	return "decorator"
}
func (d *decorator) Name() string {
	return "decorator"
}

func (d *decorator) Shutdown() error {
	return nil
}

func (d *decorator) NotifyNewPolicy(pol string, fakeSettings map[string]any) {
	// This mimicks what happens on policy change -- we look
	// through our map, and instantiate new plugins.
	var policyPlugin SettingsInjectablePlugin
	if _, found := d.decorated[pol]; !found {
		policyPlugin = d.newPluginCallback().(SettingsInjectablePlugin)
		d.decorated[pol] = policyPlugin
		if err := d.decorated[pol].Startup(); err != nil {
			logger.Warn("Failed to startup plugin %s with error: %s\n", d.Name(), err.Error())
		}
	} else {
		policyPlugin = d.decorated[pol]
	}

	// Here is where we would do settings injection.  The
	// GetNewSettings method returns a new, blank object that we
	// can unmarshall into. Here, we use mapstructure, we could
	// also use the regular JSON unmarshaller to unmarshal into
	// this object.  Later, we call SetSettings, which tells the
	// object about the new settings. Real plugins would lock
	// their settings and set their internal settings object they
	// are using to this new one. See how the MockPlugin
	// SetSettings() works.
	settings := d.decorated[pol].GetNewSettings()
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  settings,
	}
	decoder, _ := mapstructure.NewDecoder(config)
	if err := decoder.Decode(fakeSettings); err != nil {
		logger.Warn("Failed to decode the given raw interface: %s\n", err.Error())

	}
	policyPlugin.SetSettings(settings)
}

func newWrapperTest(d *decorator) *wrapperTest {
	w := &wrapperTest{matches: true}
	w.SetConstructorReturn(
		ConstructorWrapperPluginFactory(func(gen PluginGeneratorCallback, _ ...any) Plugin {
			d.newPluginCallback = gen
			return d
		}))
	return w
}

func (w *wrapperTest) Matches(val PluginConstructor, metadata ...any) bool {
	// ideally this should examine 'val' to decide if it's the
	// type of plugin we want to wrap.
	return w.matches
}

func NewMockPlugin(config *Config) *MockPlugin {
	m := &MockPlugin{config: config}
	m.On("Startup").Maybe().Return(nil)
	m.On("Shutdown").Maybe().Return(nil)
	return m
}

func TestWrapperRegistrationOrder(t *testing.T) {
	testLogic := func(t *testing.T, register func(controller *PluginControl, decorator *decorator)) {
		controller := NewPluginControl()
		decorator := &decorator{decorated: map[string]SettingsInjectablePlugin{}}

		register(controller, decorator)

		err := controller.Provide(func() *Config {
			return &Config{}
		})
		assert.NoError(t, err)

		controller.Startup()

		// NotifyNewPolicy will create a new instance of the wrapped
		// plugin, and inject the settings into it.
		decorator.NotifyNewPolicy("policy1", map[string]any{
			"name": "policy1settings",
			"id":   "myIDpolicy1plugin",
		})
		decorator.NotifyNewPolicy("policy2", map[string]any{
			"name": "policy2settings",
			"id":   "myIDpolicy2plugin",
		})
		assert.NotNil(t, decorator.decorated["policy1"])
		assert.NotNil(t, decorator.decorated["policy2"])

		assert.Equal(t, decorator.decorated["policy1"].(*MockPlugin).config.Name,
			"policy1settings")
		assert.Equal(t,
			decorator.decorated["policy1"].(*MockPlugin).config.ID,
			"myIDpolicy1plugin")

		assert.Equal(t, decorator.decorated["policy2"].(*MockPlugin).config.Name,
			"policy2settings")
		assert.Equal(t, decorator.decorated["policy2"].(*MockPlugin).config.ID,
			"myIDpolicy2plugin")

		decorator.NotifyNewPolicy("policy2", map[string]any{
			"name": "policy2settingsSecondTime",
			"id":   "myIDpolicy2plugin",
		})
		assert.Equal(t, decorator.decorated["policy2"].(*MockPlugin).config.Name,
			"policy2settingsSecondTime")
	}

	t.Run("WrapperRegisteredFirst", func(t *testing.T) {
		testLogic(t, func(controller *PluginControl, decorator *decorator) {
			controller.RegisterConstructorWrapper(func() ConstructorWrapper {
				return newWrapperTest(decorator)
			})
			controller.RegisterPlugin(NewMockPlugin)
		})
	})

	t.Run("PluginRegisteredFirst", func(t *testing.T) {
		testLogic(t, func(controller *PluginControl, decorator *decorator) {
			controller.RegisterPlugin(NewMockPlugin)
			controller.RegisterConstructorWrapper(func() ConstructorWrapper {
				return newWrapperTest(decorator)
			})
		})
	})
}

func TestWrapperNoMatch(t *testing.T) {
	controller := NewPluginControl()
	decorator := &decorator{decorated: map[string]SettingsInjectablePlugin{}}
	wrapper := newWrapperTest(decorator)
	wrapper.matches = false

	controller.RegisterConstructorWrapper(func() ConstructorWrapper {
		return wrapper
	})
	err := controller.Provide(func() *Config {
		return &Config{Name: "unwrapped"}
	})
	assert.NoError(t, err)

	var createdPlugin *MockPlugin
	controller.RegisterPlugin(func(cfg *Config) *MockPlugin {
		p := NewMockPlugin(cfg)
		createdPlugin = p
		return p
	})
	controller.Startup()

	// wrapper was not used, so decorator is untouched
	assert.Nil(t, decorator.newPluginCallback)

	// original plugin was created
	assert.NotNil(t, createdPlugin)
	assert.Equal(t, "unwrapped", createdPlugin.config.Name)

	// ensure it was started
	createdPlugin.AssertCalled(t, "Startup")

	// and it's in the plugin list
	assert.Len(t, controller.pluginInfo, 1)
	assert.Same(t, createdPlugin, controller.pluginInfo[0].plugin)
}

func TestWrapperWithProvidedPlugin(t *testing.T) {
	controller := NewPluginControl()
	decorator := &decorator{decorated: map[string]SettingsInjectablePlugin{}}
	controller.RegisterConstructorWrapper(func() ConstructorWrapper {
		return newWrapperTest(decorator)
	})
	err := controller.Provide(func() *Config {
		return &Config{}
	})
	assert.NoError(t, err)

	controller.RegisterAndProvidePlugin(NewMockPlugin)
	controller.Startup()

	// The decorator should be the only plugin in the list.
	assert.Len(t, controller.pluginInfo, 1)
	assert.Same(t, decorator, controller.pluginInfo[0].plugin)

	// Assert wrapper was used
	assert.NotNil(t, decorator.newPluginCallback)
	// Now check if decoration works
	decorator.NotifyNewPolicy("policy1", map[string]any{
		"name": "policy1settings",
		"id":   "myIDpolicy1plugin",
	})
	assert.NotNil(t, decorator.decorated["policy1"])
	assert.Equal(t, "policy1settings", decorator.decorated["policy1"].(*MockPlugin).config.Name)
}

func TestWrapperFactoryReturnsWrongType(t *testing.T) {
	controller := NewPluginControl()
	assert.PanicsWithValue(t, "from RegisterConstructorWrapper: Unable to convert provided wrapper to PluginWrapper", func() {
		controller.RegisterConstructorWrapper(func() any {
			return struct{}{}
		})
	})
}

func TestWrapperFactoryWithMissingDependency(t *testing.T) {
	controller := NewPluginControl()
	type SomeDependency struct{}
	assert.Panics(t, func() {
		controller.RegisterConstructorWrapper(func(*SomeDependency) ConstructorWrapper {
			return newWrapperTest(nil)
		})
	})
}

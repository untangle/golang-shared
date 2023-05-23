package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type wrapperTest struct {
	WrapperHelper
}

type decorator struct {
	decorated         map[string]SettingsInjectablePlugin
	newPluginCallback func() Plugin
}

func (d *decorator) Startup() error {
	return nil
}

func (d *decorator) Name() string {
	return "decorator"
}

func (d *decorator) Shutdown() error {
	return nil
}

func (d *decorator) NotifyNewPolicy(pol string) {
	// This mimicks what happens on policy change -- we look
	// through our map, and instantiate new plugins.
	if old, found := d.decorated[pol]; found {
		old.Shutdown()
	}
	policyPlugin := d.newPluginCallback().(SettingsInjectablePlugin)
	d.decorated[pol] = policyPlugin
	d.decorated[pol].Startup()
	settings := d.decorated[pol].GetNewSettings()
	conf := settings.(*Config)
	conf.Name = pol
	policyPlugin.SetSettings(conf)
}

func newWrapperTest(d *decorator) *wrapperTest {
	w := &wrapperTest{}
	w.SetConstructorReturn(d)
	d.newPluginCallback = w.GenNewPlugin
	return w
}

func (w *wrapperTest) Matches(val PluginConstructor) bool {
	// ideally this should examine 'val' to decide if it's the
	// type of plugin we want to wrap.
	return true
}

func TestWrapper(t *testing.T) {
	controller := NewPluginControl()
	decorator := &decorator{decorated: map[string]SettingsInjectablePlugin{}}
	wrapper := newWrapperTest(decorator)
	controller.RegisterWrapper(wrapper)
	controller.Provide(func() *Config {
		return &Config{}
	})
	controller.RegisterPlugin(NewPlugin)
	controller.Startup()
	decorator.NotifyNewPolicy("policy1")
	decorator.NotifyNewPolicy("policy2")
	assert.NotNil(t, decorator.decorated["policy1"])
	assert.NotNil(t, decorator.decorated["policy2"])
	assert.Equal(t, decorator.decorated["policy1"].(*MockPlugin).config.Name, "policy1")
}

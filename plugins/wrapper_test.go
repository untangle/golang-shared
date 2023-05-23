package plugins

import (
	"testing"
)

type wrapperTest struct {
	WrapperHelper
}

type decorator struct {
	decorated         map[string]Plugin
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
	if old, found := d.decorated[pol]; found {
		old.Shutdown()
	}
	d.decorated[pol] = d.newPluginCallback()
	d.decorated[pol].Startup()
}

func newWrapperTest(d *decorator) *wrapperTest {
	w := &wrapperTest{}
	w.SetConstructorReturn(d)
	d.newPluginCallback = w.GenNewPlugin
	return w
}

func (w *wrapperTest) Matches(val PluginConstructor) bool {
	return true
}

func TestWrapper(t *testing.T) {
	controller := NewPluginControl()
	decorator := &decorator{decorated: map[string]Plugin{}}
	wrapper := newWrapperTest(decorator)
	controller.RegisterWrapper(wrapper)
	controller.Provide(func() *Config {
		return &Config{}
	})
	controller.RegisterPlugin(NewPlugin)
	controller.Startup()
	decorator.NotifyNewPolicy("policy1")
	decorator.NotifyNewPolicy("policy2")

}

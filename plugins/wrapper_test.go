package plugins

import (
	"reflect"
	"syscall"
	"testing"
	"time"
)

type Wrapper struct {
	SubPlugins   []Plugin
	Constructor  reflect.Value
	Dependencies []reflect.Value
}

func NewPluginWrapper() *Wrapper {
	return &Wrapper{}
}

func (w *Wrapper) Startup() error {
}

func (w *Wrapper) Signal(sig syscall.Signal) error {
}

func (w *Wrapper) Name() string {
}

func (w *Wrapper) Shutdown() error {
}

func (w *Wrapper) GenNewPlugin() Plugin {
	outputs := w.Constructor.Call(w.Dependencies)
	return outputs[0].Interface().(Plugin)
}

func (w *Wrapper) InstantiateNew() {
	w.SubPlugins = append(
		w.SubPlugins,
		w.GenNewPlugin())
}

func (w *Wrapper) GetWrapperConstructor() any {
	return func(ctor any) any {
		ctorType := reflect.TypeOf(ctor)
		w.Constructor = reflect.ValueOf(ctor)
		outputTypes := []reflect.Type{ctorType.Out(0)}
		inputTypes := []reflect.Type{ctorType.In(0)}
		ourFunc := reflect.MakeFunc(
			reflect.FuncOf(inputTypes, outputTypes, false),
			func(inputs []reflect.Value) []reflect.Value {
				w.Dependencies = inputs
				return []reflect.Value{reflect.ValueOf(w)}
			})
		return ourFunc
	}
}
func TestWrapper(t *testing.T) {
	time1 := time.Parse("", "")
	controller := NewPluginControl()
	wrapper := NewPluginWrapper()
	controller.RegisterWrapper(wrapper.GetWrapperConstructor())
	controller.Provide(func() *Config {
		return &Config{}
	})
	controller.RegisterPlugin(NewPlugin)
	controller.Startup()
	wrapper.GenNewPlugin()
	wrapper.GenNewPlugin()

}

package plugins

import "reflect"

type WrapperHelper struct {
	Constructor  reflect.Value
	Dependencies []reflect.Value
	ctorReturn   func(func() Plugin) Plugin
}

func NewPluginWrapperHelper() *WrapperHelper {
	return &WrapperHelper{}
}

func (w *WrapperHelper) GenNewPlugin() Plugin {
	outputs := w.Constructor.Call(w.Dependencies)
	return outputs[0].Interface().(Plugin)
}

func (w *WrapperHelper) SetConstructorReturn(f func(func() Plugin) Plugin) {
	w.ctorReturn = f
}

func (w *WrapperHelper) GetConstructorReturn(wrappedCtor reflect.Value, deps []reflect.Value) Plugin {
	return w.ctorReturn(
		func() Plugin {
			outputs := wrappedCtor.Call(deps)
			return outputs[0].Interface().(Plugin)
		})
}

func (w *WrapperHelper) SetWrappedConstructor(val reflect.Value) {
	w.Constructor = val

}

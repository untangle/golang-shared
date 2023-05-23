package plugins

import "reflect"

type WrapperHelper struct {
	Constructor  reflect.Value
	Dependencies []reflect.Value
	ctorReturn   any
}

func NewPluginWrapperHelper() *WrapperHelper {
	return &WrapperHelper{}
}

func (w *WrapperHelper) SetConstructorDependencies(deps []reflect.Value) {
	w.Dependencies = deps
}

func (w *WrapperHelper) GenNewPlugin() Plugin {
	outputs := w.Constructor.Call(w.Dependencies)
	return outputs[0].Interface().(Plugin)
}

func (w *WrapperHelper) SetConstructorReturn(val any) {
	w.ctorReturn = val
}

func (w *WrapperHelper) GetConstructorReturn() any {
	return w.ctorReturn

}

func (w *WrapperHelper) SetWrappedConstructor(val reflect.Value) {
	w.Constructor = val

}

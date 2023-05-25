package plugins

import "reflect"

// PluginGeneratorCallback is a callback function that generates a
// plugin. It is used to hide the complexity of the reflect package
// from users of WrapperHelper.
type PluginGeneratorCallback func() Plugin

// ConstructorWrapperPluginFactory is a function that takes a
// PluginGeneratorCallback and returns a Plugin -- the purpose is to
// allow a constructor wrapper to get a function that would generate
// the plugin it is wrapping, and for it to return the plugin it would
// like to substitute. It's second variadic argument is the plugin
// constructor metadata, which was passed when the plugin was registered.
type ConstructorWrapperPluginFactory func(PluginGeneratorCallback, ...any) Plugin

// WrapperHelper is a helper object for ConstructorWrapper, embed it
// in some wrapper you'd like to make.
type WrapperHelper struct {
	ctorReturn ConstructorWrapperPluginFactory
}

// NewPluginWrapperHelper returns a WrapperHelper.
func NewPluginWrapperHelper() *WrapperHelper {
	return &WrapperHelper{}
}

// SetConstructorReturn -- call this from your child class or client
// that would like to wrap plugin constructors. Give it a
// ConstructorWrapperPluginFactory -- that factory f will get passed
// as an argument a function that will return when called new
// instances of the wrapped plugin. f will then return instead what it
// would like to substitute for that plugin that would have gotten
// instantiated.
func (w *WrapperHelper) SetConstructorReturn(f ConstructorWrapperPluginFactory) {
	w.ctorReturn = f
}

// GetConstructorReturn -- this is what clients don't have to worry
// about. It takes care of reflection for you and is called by the
// plugin control when RegisterConstructorWrapper is called.
func (w *WrapperHelper) GetConstructorReturn(wrappedCtor reflect.Value, deps []reflect.Value, metadata ...any) Plugin {
	return w.ctorReturn(
		// This is the 'generator' -- we store the
		// reflect.Value of the constructor and its
		// dependencies in this closure, so you don't have to
		// mess with reflect.
		func() Plugin {
			outputs := wrappedCtor.Call(deps)
			return outputs[0].Interface().(Plugin)
		},
		metadata...)
}

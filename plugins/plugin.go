package plugins

import (
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/dig"

	logService "github.com/untangle/golang-shared/services/logger"
)

var logger = logService.GetLoggerInstance()

// Plugin is an interface for (right now only nfqueue) plugins.
type Plugin interface {
	Startup() error
	Name() string
	Shutdown() error
}

// PluginConstructor is a function that generates a plugin.
type PluginConstructor interface{}

// PluginConsumer is a function that consumes a particular interface
// that a plugin may fulfill that it is interested in.
type PluginConsumer interface{}

type consumer struct {
	consumedType reflect.Type
	consumerFunc reflect.Value
}

// PluginControl controls plugins. It controls construction of plugins
// with the dig DI container and also keeps a list of running plugins
// to call their Signal() and Shutdown() methods.  When you register a
// plugin, the PluginControl adds it to the DI container so that it
// can wire dependencies when Startup() is called. After that it just
// keeps track of a list of plugins to send method calls to.
type PluginControl struct {
	dig.Container
	wrapper            ConstructorWrapper
	plugins            []Plugin
	saverFuncs         []reflect.Value
	consumers          []consumer
	enableStartupPanic bool
}

// NewPluginControl creates an empty PluginControl
func NewPluginControl() *PluginControl {
	container := dig.New()
	return &PluginControl{
		Container:          *container,
		enableStartupPanic: true}
}

var pluginControl *PluginControl

// GlobalPluginControl returns the plugin control singleton.
func GlobalPluginControl() *PluginControl {
	if pluginControl != nil {
		return pluginControl
	}
	pluginControl = NewPluginControl()
	return pluginControl
}

type ConstructorWrapper interface {
	// Matches returns true if we'd like to wrap this plugin.
	Matches(PluginConstructor, ...any) bool

	// GetConstructorReturn returns the plugin you'd like to use
	// _instead_ of the plugin that the wrappedConstructor
	// argument would return.
	GetConstructorReturn(wrappedConstructor reflect.Value, deps []reflect.Value, metadata ...any) Plugin
}

// ConstructorWrapperFactory is a function that takes any number of
// arguments of some type, and returns a wrapper.
type ConstructorWrapperFactory any

func makeWrapperConstructor(
	wrapper ConstructorWrapper, ctor any, metadata []any) reflect.Value {
	ctorType := reflect.TypeOf(ctor)
	inputTypes := make([]reflect.Type, ctorType.NumIn(), ctorType.NumIn())
	for t := 0; t < ctorType.NumIn(); t++ {
		inputTypes[t] = ctorType.In(t)
	}

	ourFunc := reflect.MakeFunc(
		reflect.FuncOf(inputTypes, []reflect.Type{reflect.TypeOf((*Plugin)(nil)).Elem()}, false),
		func(inputs []reflect.Value) []reflect.Value {
			returnvalue := reflect.ValueOf(wrapper.GetConstructorReturn(
				reflect.ValueOf(ctor), inputs, metadata...))
			return []reflect.Value{returnvalue}

		})
	return ourFunc
}

// RegisterConstructorWrapper registers the 'constructor
// wrapper'. wrapper must be a function returning a
// ConstructorWrapper. It is instantiated using the objects available
// to the DI currently (i.e. the constructor passed can take
// Provide()-ed objects).  Once registered, all new registered plugins
// registered via RegisterPlugin() will have their constructors passed
// to the Matches() function of the ConstructorWrapper returned by
// wrapper, if they match, then instead of that constructor getting
// called when that plugin would normally be instantiated, first its
// constructor reflect.Value, then its dependencies (constructor
// arguments) are instead passed to the GetConstructorReturn()
// function, and the plugin returned by that function will be used
// instead of what would have been returned by the regular
// constructor.
func (control *PluginControl) RegisterConstructorWrapper(wrapper ConstructorWrapperFactory) {
	if err := control.Provide(wrapper); err != nil {
		panic(fmt.Sprintf("couldn't provide wrapper: %s", err))
	}
	if err := control.Invoke(func(w ConstructorWrapper) {
		control.wrapper = w
	}); err != nil {
		panic(fmt.Sprintf("couldn't instantiate wrapper: %s", err))
	}
}

// RegisterPlugin registers a plugin that will be created during the
// Startup() method and provided with its dependencies. constructor is
// a function that takes arbitrary types of arguments to be provided
// by the DI container and returns a plugin object. This function will
// not provide the plugin as a potential dependency for other plugins.
func (control *PluginControl) RegisterPlugin(constructor PluginConstructor, metadata ...any) {
	constructorType := reflect.TypeOf(constructor)
	constructorVal := reflect.ValueOf(constructor)
	inputs := []reflect.Type{}
	for i := 0; i < constructorType.NumIn(); i++ {
		inputs = append(inputs, constructorType.In(i))
	}

	if control.wrapper != nil && control.wrapper.Matches(constructor, metadata...) {
		constructorVal = makeWrapperConstructor(control.wrapper, constructor, metadata)
	}
	// create a func at runtime that we can invoke that calls the
	// constructor and appends the return value to the list of plugins.
	saverFunc := reflect.MakeFunc(
		reflect.FuncOf(inputs, []reflect.Type{}, false),
		func(vals []reflect.Value) []reflect.Value {
			output := constructorVal.Call(vals)
			plugin := output[0].Interface()
			pluginIntf := plugin.(Plugin)
			control.plugins = append(control.plugins, pluginIntf)
			return []reflect.Value{}
		})
	control.saverFuncs = append(control.saverFuncs, saverFunc)
}

// RegisterAndProvidePlugin registers a plugin that may be consumed by
// other plugins. This constructor function therefore needs a unique
// type. The constructor will be added to the DI container, and other
// plugins may require it. It is not instantiated until the Startup()
// method is called.
func (control *PluginControl) RegisterAndProvidePlugin(constructor PluginConstructor) {
	constructorType := reflect.TypeOf(constructor)
	outputType := constructorType.Out(0)

	// create a func at runtime that we can invoke that requires
	// the plugin to ensure it gets instantiated, and also appends
	// it to the list of registered plugins.
	saverFunc := reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{outputType}, []reflect.Type{}, false),
		func(vals []reflect.Value) []reflect.Value {
			plugin := vals[0].Interface()
			pluginIntf := plugin.(Plugin)
			control.plugins = append(control.plugins, pluginIntf)
			return []reflect.Value{}
		})
	control.saverFuncs = append(control.saverFuncs, saverFunc)
	if err := control.Provide(constructor); err != nil {
		panic(fmt.Sprintf(
			"couldn't register plugin constructor as a provider: %v, err: %s", constructor, err))
	}
}

// Unregister removes the a plugin from list of plugins
func (control *PluginControl) Unregister(name string) error {

	// Search the plugins slice and remove
	// TODO: When we migrate to go1.21, do this with the "slice" package instead of manually searching
	for indx, plugin := range control.plugins {
		if plugin.Name() == name {
			// remove the index from the slice
			control.plugins = append(control.plugins[:indx], control.plugins[indx+1:]...)
		} else {
		logger.Info("Could not unregister plugin. Plugin not found in list %s\n", plugin.Name())
            return errors.New("could not unregister plugin")
		}
	}
	return nil
}

// Startup constructs and then starts all registered plugins. It
// panics if any don't start. So it will call the constructor passed
// to RegisterPlugin with whatever arguments it requires (obtained via
// the DI container), and then call the Startup() method. Finally, if
// the plugin satisfies any of the interfaces ConnectionTrackerPlugin,
// NetlogHandler, or PacketProcessorPlugin, their handler methods are
// registered with the backend so they will receive these events.
func (control *PluginControl) Startup() {
	for _, saverFunc := range control.saverFuncs {
		if err := control.Invoke(saverFunc.Interface()); err != nil {
			panic(fmt.Sprintf("couldn't instantiate plugin: %s", err))
		}
	}

	for _, plugin := range control.plugins {
		logger.Info("Starting plugin: %s\n", plugin.Name())
		if err := plugin.Startup(); err != nil {

			if unregErr := control.Unregister(plugin.Name()); unregErr != nil {
				logger.Crit("couldn't unregister plugin %s after failed startup: %s",
					plugin.Name(),
					unregErr)
			}

			if control.enableStartupPanic {
				panic(fmt.Sprintf("couldn't startup plugin %s: %s",
					plugin.Name(),
					err))
			} else {
				logger.Crit("couldn't startup plugin %s: %s",
					plugin.Name(),
					err)
			}
		} else {
			control.findConsumers(plugin)
		}

	}

}

// LogStartupErrors sets the PluginControl to just log errors when
// plugins start rather than panicking.
func (control *PluginControl) LogStartupErrors() {
	control.enableStartupPanic = false
}

// PanicOnStartupErrors sets the PluginControl to panic() when a
// Startup() method returns an error.
func (control *PluginControl) PanicOnStartupErrors() {
	control.enableStartupPanic = true
}

// find and call the consumer functions that consume plugins.
func (control *PluginControl) findConsumers(plugin interface{}) {
	pluginType := reflect.TypeOf(plugin)
	pluginValue := reflect.ValueOf(plugin)
	for _, consumer := range control.consumers {
		if pluginType.Implements(consumer.consumedType) {
			theFunc := consumer.consumerFunc
			args := []reflect.Value{pluginValue}
			theFunc.Call(args)
		}
	}
}

// RegisterConsumer registers the plugin consumer function. The
// function should take a single argument which is some interface it
// wants the plugin to satisfy. Then after plugin startup, the
// consumer will be passed the started plugin. This allows you to
// define your own plugin interface and consume it.
func (control *PluginControl) RegisterConsumer(theConsumer PluginConsumer) {
	consumerFunc := reflect.TypeOf(theConsumer)
	expectedType := consumerFunc.In(0)
	control.consumers = append(control.consumers,
		consumer{
			consumedType: expectedType,
			consumerFunc: reflect.ValueOf(theConsumer),
		})

}

// StopAllPlugins stops all registered plugins.
func (control *PluginControl) Shutdown() {
	for _, plugin := range control.plugins {
		if err := plugin.Shutdown(); err != nil {
			logger.Warn("Plugin %s failed to stop: %s\n", plugin.Name(), err)
		} else {
			logger.Info("Shutdown: %s\n", plugin.Name())
		}
	}
}

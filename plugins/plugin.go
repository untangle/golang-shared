package plugins

import (
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

type pluginInfo struct {
	plugin      Plugin
	metadata    []any
	constructor PluginConstructor
	isProvider  bool

	// a function built using the reflect package that will set
	// the plugin value of this struct. Executed during Startup().
	saverFunc interface{}
}

type predicateInfo struct {
	invokerFunc     interface{}
	pluginPredicate PluginPredicate
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
	wrapperInvokerFunc interface{}
	predicates         []*predicateInfo
	pluginInfo         []*pluginInfo

	consumers          []consumer
	enableStartupPanic bool
}

// NewPluginControl creates an empty PluginControl
func NewPluginControl() *PluginControl {
	container := dig.New()
	ctrl := &PluginControl{
		Container:          *container,
		enableStartupPanic: false}
	err := ctrl.Provide(func() *PluginControl {
		return ctrl
	})
	if err != nil {
		logger.Warn("Failed to provide return value: %v\n", err.Error())
	}
	return ctrl
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

type PluginPredicate interface {
	// IsRelevant determines if a plugin constructor should be
	// included based on platform or other criteria.
	IsRelevant(constructor PluginConstructor, metadata ...any) bool
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

// PluginPredicateFactory is a function that takes any number of
// arguments (which will be supplied by DI in the typical case) and
// returns a PluginPredicate.
type PluginPredicateFactory any

func makeWrapperConstructor(
	wrapper ConstructorWrapper, ctor any, metadata []any) reflect.Value {
	ctorType := reflect.TypeOf(ctor)
	inputTypes := make([]reflect.Type, ctorType.NumIn())
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
	// Here we instantiate the constructor wrapper, of whatever
	// type it is.
	constructorType := reflect.TypeOf(wrapper)
	outputToGet := constructorType.Out(0)
	if err := control.Provide(wrapper); err != nil {
		panic(fmt.Sprintf("couldn't provide wrapper: %s", err))
	}
	invokerFunc := reflect.MakeFunc(reflect.FuncOf([]reflect.Type{outputToGet}, []reflect.Type{}, false),
		func(vals []reflect.Value) []reflect.Value {
			wrapper, ok := vals[0].Interface().(ConstructorWrapper)
			if !ok {
				panic("from RegisterConstructorWrapper: Unable to convert provided wrapper to PluginWrapper")
			}
			control.wrapper = wrapper
			return []reflect.Value{}
		})
	control.wrapperInvokerFunc = invokerFunc.Interface()
}

// RegisterPluginPredicate registers to the 'control' object a plugin
// predicate. predicateFactory is a function that returns a
// PluginPredicate, the arguments will be provided by the DI container
// so you can use dependency injection here. Immediately before
// plugins are instantiated (i.e. their constructors called), all
// predicates are instantiated by calling the factories. Then, the
// newly-instantiated predicate's "IsRelevant" method will then be
// used to conditionalize plugin instantiation/creation: if _any_
// plugin does not pass _all_ predicate checks, it will not be
// instantiated by the DI framework (the constructor will not be
// called), and will also of course not be started.
func (control *PluginControl) RegisterPluginPredicate(predicateFactory PluginPredicateFactory) {
	constructorType := reflect.TypeOf(predicateFactory)
	outputToGet := constructorType.Out(0)
	if err := control.Provide(predicateFactory); err != nil {
		panic(fmt.Sprintf("couldn't provide predicate: %s", err))
	}
	predicateInfo := &predicateInfo{
		invokerFunc:     nil,
		pluginPredicate: nil,
	}
	invokerFunc := reflect.MakeFunc(reflect.FuncOf([]reflect.Type{outputToGet}, []reflect.Type{}, false),
		func(vals []reflect.Value) []reflect.Value {
			pred, ok := vals[0].Interface().(PluginPredicate)
			if !ok {
				panic("from RegisterPluginPredicate: Unable to convert provided wrapper to PluginPredicate")
			}
			predicateInfo.pluginPredicate = pred
			return []reflect.Value{}
		})
	predicateInfo.invokerFunc = invokerFunc.Interface()
	control.predicates = append(control.predicates, predicateInfo)
}

// RegisterPlugin registers a plugin that will be created during the
// Startup() method and provided with its dependencies. constructor is
// a function that takes arbitrary types of arguments to be provided
// by the DI container and returns a plugin object. This function will
// not provide the plugin as a potential dependency for other plugins.
func (control *PluginControl) RegisterPlugin(constructor PluginConstructor, metadata ...any) {
	pluginInfo := &pluginInfo{
		constructor: constructor,
		metadata:    metadata,
		isProvider:  false,

		// to be set later by our constructed function, if
		// needed.
		plugin: nil,
	}
	control.pluginInfo = append(control.pluginInfo, pluginInfo)
}

// Function to return a count of the registered Plugins
// Intended for unit testing only
func (control *PluginControl) GetRegisteredPluginCount() int {
	// control.saverFuncs is appended to by the RegisterPlugin() above
	// and the RegistRegisterAndProvidePluginerPlugin() function below
	// The saverFuncs are invoked on Startup() to actually start plugins.
	return len(control.pluginInfo)
}

// RegisterAndProvidePlugin registers a plugin that may be consumed by
// other plugins. This constructor function therefore needs a unique
// type. The constructor will be added to the DI container, and other
// plugins may require it. It is not instantiated until the Startup()
// method is called.
func (control *PluginControl) RegisterAndProvidePlugin(constructor PluginConstructor, metadata ...any) {
	pluginInfo := &pluginInfo{
		plugin:      nil,
		metadata:    metadata,
		constructor: constructor,
		saverFunc:   nil,
		isProvider:  true,
	}
	control.pluginInfo = append(control.pluginInfo, pluginInfo)
}

// UnregisterPlugin removes a plugin from the list of plugins
func (control *PluginControl) UnregisterPluginByIndex(indx int) {
	control.pluginInfo = append(control.pluginInfo[:indx], control.pluginInfo[indx+1:]...)
}

func (control *PluginControl) preparePlugins() {
	for _, pluginInfo := range control.pluginInfo {
		constructor := pluginInfo.constructor
		constructorVal := reflect.ValueOf(constructor)

		// first, see if we need to wrap the constructor.
		if control.wrapper != nil && control.wrapper.Matches(constructor, pluginInfo.metadata...) {
			constructorVal = makeWrapperConstructor(control.wrapper, constructor, pluginInfo.metadata)
		}

		actualConstructor := constructorVal.Interface()

		if pluginInfo.isProvider {
			constructorType := reflect.TypeOf(actualConstructor)
			outputType := constructorType.Out(0)

			saverFunc := reflect.MakeFunc(
				reflect.FuncOf([]reflect.Type{outputType}, []reflect.Type{}, false),
				func(vals []reflect.Value) []reflect.Value {
					plugin := vals[0].Interface()
					pluginIntf := plugin.(Plugin)
					pluginInfo.plugin = pluginIntf
					return []reflect.Value{}
				})
			pluginInfo.saverFunc = saverFunc.Interface()
			if err := control.Provide(actualConstructor); err != nil {
				panic(fmt.Sprintf(
					"couldn't register plugin constructor as a provider: %v, err: %s", constructor, err))
			}
		} else {
			constructorType := reflect.TypeOf(actualConstructor)
			inputs := []reflect.Type{}
			for i := 0; i < constructorType.NumIn(); i++ {
				inputs = append(inputs, constructorType.In(i))
			}
			saverFunc := reflect.MakeFunc(
				reflect.FuncOf(inputs, []reflect.Type{}, false),
				func(vals []reflect.Value) []reflect.Value {
					output := constructorVal.Call(vals)
					plugin := output[0].Interface()
					pluginIntf := plugin.(Plugin)
					pluginInfo.plugin = pluginIntf
					return []reflect.Value{}
				})
			pluginInfo.saverFunc = saverFunc.Interface()
		}
	}
}

// filterPlugins filters the plugin list using the predicates.
func (control *PluginControl) filterPlugins() {
	// instantiate the predicates.
	for _, pred := range control.predicates {
		if err := control.Invoke(pred.invokerFunc); err != nil {
			panic(fmt.Sprintf("couldn't instantiate plugin predicate: %v",
				err))
		}

		if pred.pluginPredicate == nil {
			panic("plugin predicate is nil")
		}
	}

	plugins := []*pluginInfo{}
PluginLoop:
	for _, plugin := range control.pluginInfo {
		for _, pred := range control.predicates {
			// Filter out irrelevant plugins *before* creating the saver function
			if !pred.pluginPredicate.IsRelevant(
				plugin.constructor,
				plugin.metadata...) {
				continue PluginLoop // Skip this plugin
			}
		}
		plugins = append(plugins, plugin)
	}
	control.pluginInfo = plugins
}

// Startup constructs and then starts all registered plugins. It
// panics if any don't start. So it will call the constructor passed
// to RegisterPlugin with whatever arguments it requires (obtained via
// the DI container), and then call the Startup() method. Finally, if
// the plugin satisfies any of the interfaces ConnectionTrackerPlugin,
// NetlogHandler, or PacketProcessorPlugin, their handler methods are
// registered with the backend so they will receive these events.
func (control *PluginControl) Startup() {
	if control.wrapperInvokerFunc != nil && control.wrapper == nil {
		if err := control.Invoke(control.wrapperInvokerFunc); err != nil {
			panic(fmt.Sprintf("couldn't instantiate wrapper: %s", err))
		}
	}
	control.filterPlugins()
	control.preparePlugins()
	for _, pluginInf := range control.pluginInfo {
		if err := control.Invoke(pluginInf.saverFunc); err != nil {
			err = fmt.Errorf("couldn't instantiate plugin with constructor %T: %w", pluginInf.constructor, err)
			if control.enableStartupPanic {
				panic(err)
			}
			logger.Crit("%s\n", err.Error())
		}
	}

	successfulPlugins := []*pluginInfo{}
	for _, pluginInf := range control.pluginInfo {
		if pluginInf.plugin == nil {
			// Instantiation failed for this plugin, skip it.
			continue
		}
		plugin := pluginInf.plugin
		logger.Info("Starting plugin: %s\n", plugin.Name())
		if err := plugin.Startup(); err != nil {

			if control.enableStartupPanic {
				panic(fmt.Sprintf("couldn't startup plugin %s: %s",
					plugin.Name(),
					err))
			} else {
				logger.Crit("couldn't startup plugin %s: %s\n",
					plugin.Name(),
					err)
			}

		} else {
			successfulPlugins = append(successfulPlugins, pluginInf)
			control.findConsumers(pluginInf.plugin)
		}
	}

	control.pluginInfo = successfulPlugins

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
	for _, pluginInf := range control.pluginInfo {
		plugin := pluginInf.plugin
		if err := plugin.Shutdown(); err != nil {
			logger.Warn("Plugin %s failed to stop: %s\n", plugin.Name(), err)
		} else {
			logger.Info("Shutdown: %s\n", plugin.Name())
		}
	}
}

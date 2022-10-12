package plugins

import (
	"fmt"
	"reflect"
	"syscall"

	"github.com/untangle/golang-shared/services/logger"
	"go.uber.org/dig"
)

// Plugin is an interface for (right now only nfqueue) plugins.
type Plugin interface {
	Startup() error
	Name() string
	Signal(syscall.Signal) error
	Shutdown() error
}

// PluginConstructor is a function that generates a plugin.
type PluginConstructor interface{}

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
	plugins            []Plugin
	pluginConstructors []PluginConstructor
	saverFuncs         []reflect.Value
	consumers          []consumer
}

// NewPluginControl creates an empty PluginControl
func NewPluginControl() *PluginControl {
	container := dig.New()
	return &PluginControl{Container: *container}
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

// RegisterPlugin registers a plugin via its constructor so that when
// PluginControl.Startup(), Signal(), or Shutdown() are called, the
// plugin's methods are invoked. The plugin is not actually
// constructed until the Startup() method is called.  This method also
// adds the plugin to the container -- it registers that the
// constructor provides its return value and needs whatever arguments
// it takes. When PluginControl.Startup() is called, the DI container
// will will wire up dependencies.
func (control *PluginControl) RegisterPlugin(constructor PluginConstructor) {
	control.pluginConstructors = append(control.pluginConstructors, constructor)
	constructorType := reflect.TypeOf(constructor)
	outputType := constructorType.Out(0)

	// create a func at runtime that we can invoke that calls the
	// constructor and appends the return value to the list of plugins.
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
		panic(fmt.Sprintf("couldn't register plugin constructor: %v, err: %s", constructor, err))
	}
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
			panic(fmt.Sprintf("couldn't startup plugin %s: %s",
				plugin.Name(),
				err))
		} else {
			control.findConsumers(plugin)
		}

	}

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
func (control *PluginControl) RegisterConsumer(theConsumer interface{}) {
	consumerFunc := reflect.TypeOf(theConsumer)
	expectedType := consumerFunc.In(0)
	control.consumers = append(control.consumers,
		consumer{
			consumedType: expectedType,
			consumerFunc: reflect.ValueOf(theConsumer),
		})

}

// Signal calls the Signal() method of all registered plugins with sig.
func (control *PluginControl) Signal(sig syscall.Signal) {
	for _, plugin := range control.plugins {
		if err := plugin.Signal(sig); err != nil {
			logger.Warn("Plugin %s returned error handling signal %v: %s\n",
				plugin.Name(),
				sig,
				err)
		}
	}
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

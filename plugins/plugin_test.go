package plugins

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Fake configuration to check DI works.
type Config struct {
	Name string
}

// Mock plugin to check method calls.
type MockPlugin struct {
	mock.Mock
	config *Config
}

// global instance of the MockPlugin, to be returned by NewPlugin.
var baseMockPluginSave *MockPlugin = &MockPlugin{}

func NewPlugin(config *Config) *MockPlugin {
	baseMockPluginSave.config = config
	return baseMockPluginSave
}

func GetMockPluginConstructor() (*MockPlugin, *mock.Mock, func(config *Config) *MockPlugin) {
	baseMockPluginSave = &MockPlugin{}
	fmt.Printf("plugin save: %v (intf: %v)\n", baseMockPluginSave, (interface{}(baseMockPluginSave)).(Plugin))
	return baseMockPluginSave, &baseMockPluginSave.Mock, NewPlugin
}

func (plugin *MockPlugin) Startup() error {
	returns := plugin.Mock.Called()
	return returns.Error(0)
}

func (plugin *MockPlugin) Signal(sig syscall.Signal) error {
	returns := plugin.Called(sig)
	return returns.Error(0)
}

func (plugin *MockPlugin) Name() string {
	return "MockPlugin"
}

func (plugin *MockPlugin) Shutdown() error {
	returns := plugin.Called()
	return returns.Error(0)
}

// TestPlugin tests that basic DI works (we inject the proper config
// into the constructor). Also tests that we propagate the proper
// calls to plugin control to the actual plugin.
func TestPlugin(t *testing.T) {
	configName := "Hello Config"
	baseMockPluginSave.On("Startup").Return(nil)
	baseMockPluginSave.On("Shutdown").Return(nil)
	pluginController := NewPluginControl()
	pluginController.RegisterPlugin(
		NewPlugin)
	assert.Nil(t, pluginController.Provide(
		func() *Config {
			return &Config{Name: configName}
		}))
	pluginController.Startup()
	baseMockPluginSave.AssertNumberOfCalls(t, "Startup", 1)
	assert.Equal(t, baseMockPluginSave.config.Name, configName)
	pluginController.Shutdown()
	baseMockPluginSave.AssertNumberOfCalls(t, "Shutdown", 1)
}

/*
   Type declarations so we can test plugin consumers and test that we
   consume when something implements an interface and also test that
   when multiple interfaces are implemented, it is registered with all
   applicable consumers.
*/
type HelloType interface {
	SayHello()
}

type HelloWorldPlugin struct {
	MockPlugin
}

func (mock *HelloWorldPlugin) SayHello() {
	mock.Called()
	fmt.Println("hello!")
}

var consumedPluginSave *HelloWorldPlugin

func NewConsumedPlugin() *HelloWorldPlugin {
	return consumedPluginSave
}

func GetConsumedPluginConstructor() (*HelloWorldPlugin, *mock.Mock, func() *HelloWorldPlugin) {
	consumedPluginSave = &HelloWorldPlugin{}
	return consumedPluginSave, &consumedPluginSave.Mock, NewConsumedPlugin
}

type GoodbyeType interface {
	SayGoodbye()
}

type GoodbyeWorldPlugin struct {
	HelloWorldPlugin
}

func (mock *GoodbyeWorldPlugin) SayGoodbye() {
	mock.Called()
	fmt.Println("goodbye")
}

var multiInterfacePluginSave *GoodbyeWorldPlugin

func NewMultiInterfacePlugin() *GoodbyeWorldPlugin {
	return multiInterfacePluginSave
}

func GetMultiInterfaceConstructor() (*GoodbyeWorldPlugin, *mock.Mock, func() *GoodbyeWorldPlugin) {
	multiInterfacePluginSave = &GoodbyeWorldPlugin{}
	return multiInterfacePluginSave, &multiInterfacePluginSave.Mock, NewMultiInterfacePlugin
}

// Plugin to check that we can require other plugins.
type DependantPlugin struct {
	mock.Mock
	dependency *MockPlugin
	config     *Config
}

var dependantPluginSave *DependantPlugin

func (plugin *DependantPlugin) Startup() error {
	rvals := plugin.Called()
	return rvals.Error(0)
}

func (plugin *DependantPlugin) Shutdown() error {
	rvals := plugin.Called()
	return rvals.Error(0)

}

func (plugin *DependantPlugin) Name() string {
	return "Dependant plugin"
}

func NewDependantPlugin(config *Config, otherPlugin *MockPlugin) *DependantPlugin {
	dependantPluginSave.config = config
	dependantPluginSave.dependency = otherPlugin
	return dependantPluginSave
}

func GetDependantPluginConstructor() (*DependantPlugin, *mock.Mock, func(*Config, *MockPlugin) *DependantPlugin) {
	dependantPluginSave = &DependantPlugin{}
	return dependantPluginSave, &dependantPluginSave.Mock, NewDependantPlugin
}

// TestPluginDependenciesAndConsumption tests various use-cases for
// plugins.
func TestPluginDependenciesAndConsumption(t *testing.T) {
	helloConsumerPluginRegistry := []HelloType{}
	helloConsumer := func(thePlugin HelloType) {
		_, ok := thePlugin.(Plugin)
		assert.True(t, ok)
		helloConsumerPluginRegistry = append(helloConsumerPluginRegistry, thePlugin)
	}

	goodbyeConsumerPluginRegistry := []GoodbyeType{}
	goodbyeConsumer := func(thePlugin GoodbyeType) {
		goodbyeConsumerPluginRegistry = append(goodbyeConsumerPluginRegistry, thePlugin)
	}
	config := &Config{
		Name: "hello",
	}

	type constructorPluginPair struct {
		pluginMock  *mock.Mock  // the mock contained in the plugin.
		plugin      interface{} // the actual plugin instance.
		constructor interface{} // constructor function for the plugin (it will return the plugin instance).
		isProvider  bool        // is the plugin meant to be provided to others as a dependency?
	}

	type testConfig struct {
		plugins    []func() constructorPluginPair // list of functions that generate constructorPluginPairs for the test.
		assertions func()                         // assertions to make after Startup()
		consumers  []interface{}                  // list of plugin consumers to register as consumers.
	}
	makeConstructorPluginPair := func(
		mockPlugin interface{},
		mockValue *mock.Mock,
		constructor interface{}) constructorPluginPair {
		return constructorPluginPair{pluginMock: mockValue, plugin: mockPlugin, constructor: constructor}
	}
	makeConstructorProviderPluginPair := func(
		mockPlugin interface{},
		mockValue *mock.Mock,
		constructor interface{}) constructorPluginPair {
		return constructorPluginPair{pluginMock: mockValue, plugin: mockPlugin, constructor: constructor, isProvider: true}
	}

	tests := []testConfig{
		{
			// test that a plugin can be 'consumed' that
			// is, a function can be registered that is
			// passed all instances of a plugin that
			// implement the interface taken by that
			// function.
			plugins: []func() constructorPluginPair{
				func() constructorPluginPair { return makeConstructorPluginPair(GetConsumedPluginConstructor()) },
			},
			assertions: func() {
				require.Equal(t, 1, len(helloConsumerPluginRegistry))
				_, ok := helloConsumerPluginRegistry[0].(HelloType)
				assert.True(t, ok)
				helloConsumerPluginRegistry = []HelloType{}
			},
			consumers: []interface{}{
				helloConsumer,
			},
		},
		{
			// Test that if we do nothing, nothing happens.
			plugins: []func() constructorPluginPair{},
			assertions: func() {
				require.Equal(t, 0, len(helloConsumerPluginRegistry))
			},
			consumers: []interface{}{
				helloConsumer,
			},
		},
		{
			// Test that unrelated plugins are not consumed.
			plugins: []func() constructorPluginPair{
				func() constructorPluginPair { return makeConstructorPluginPair(GetMockPluginConstructor()) }},
			assertions: func() {
				require.Equal(t, 0, len(helloConsumerPluginRegistry))
				baseMockPluginSave.AssertCalled(t, "Startup")
			},
			consumers: []interface{}{
				helloConsumer,
			},
		},
		{
			// Test that if something implements multiple
			// interfaces, consumers listening to each are
			// notified.
			plugins: []func() constructorPluginPair{
				func() constructorPluginPair {
					return makeConstructorPluginPair(GetMultiInterfaceConstructor())
				}},
			assertions: func() {
				require.Equal(t, 1, len(helloConsumerPluginRegistry))
				require.Equal(t, 1, len(goodbyeConsumerPluginRegistry))
				_, ok := helloConsumerPluginRegistry[0].(HelloType)
				assert.True(t, ok)
				_, ok = goodbyeConsumerPluginRegistry[0].(GoodbyeType)
				assert.True(t, ok)

				// here check that the same plugin,
				// which implements multiple
				// interfaces, is registered to both
				// consumers.
				goodbye, ok := goodbyeConsumerPluginRegistry[0].(*GoodbyeWorldPlugin)
				assert.True(t, ok)
				hello, ok := helloConsumerPluginRegistry[0].(*GoodbyeWorldPlugin)
				assert.True(t, ok)
				assert.Same(t, hello, goodbye)
				goodbyeConsumerPluginRegistry = nil
				helloConsumerPluginRegistry = nil

			},
			consumers: []interface{}{
				helloConsumer,
				goodbyeConsumer,
			},
		},
		{
			// Test that inter-plugin dependencies work.
			plugins: []func() constructorPluginPair{
				func() constructorPluginPair { return makeConstructorProviderPluginPair(GetMockPluginConstructor()) },
				func() constructorPluginPair { return makeConstructorPluginPair(GetDependantPluginConstructor()) },
			},
			assertions: func() {
				// check that the dependant plugin got provided the right object during construction.
				require.Same(t, dependantPluginSave.dependency, baseMockPluginSave)
			},
			consumers: []interface{}{},
		},
	}
	for _, test := range tests {
		controller := NewPluginControl()
		for _, consumer := range test.consumers {
			controller.RegisterConsumer(consumer)
		}

		for _, pluginFunc := range test.plugins {
			plugin := pluginFunc()
			plugin.pluginMock.On("Startup").Return(nil)
			plugin.pluginMock.On("Shutdown").Return(nil)
			if plugin.isProvider {
				controller.RegisterAndProvidePlugin(plugin.constructor)
			} else {
				controller.RegisterPlugin(plugin.constructor)
			}
		}
		assert.Nil(t, controller.Provide(func() *Config { return config }))
		controller.Startup()
		test.assertions()
		controller.Shutdown()
	}
}
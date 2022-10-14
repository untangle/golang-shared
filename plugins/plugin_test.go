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

var pluginSave *MockPlugin = &MockPlugin{}

func NewPlugin(config *Config) *MockPlugin {
	pluginSave.config = config
	return pluginSave
}

func GetMockPluginConstructor() (*MockPlugin, *mock.Mock, func(config *Config) *MockPlugin) {
	pluginSave = &MockPlugin{}
	fmt.Printf("plugin save: %v (intf: %v)\n", pluginSave, (interface{}(pluginSave)).(Plugin))
	return pluginSave, &pluginSave.Mock, NewPlugin
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
	pluginSave.On("Startup").Return(nil)
	pluginSave.On("Shutdown").Return(nil)
	pluginController := NewPluginControl()
	pluginController.RegisterPlugin(
		NewPlugin)
	assert.Nil(t, pluginController.Provide(
		func() *Config {
			return &Config{Name: configName}
		}))
	pluginController.Startup()
	pluginSave.AssertNumberOfCalls(t, "Startup", 1)
	assert.Equal(t, pluginSave.config.Name, configName)
	pluginController.Shutdown()
	pluginSave.AssertNumberOfCalls(t, "Shutdown", 1)
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

type ConsumerPlugin struct {
	mock.Mock
	dependency *MockPlugin
	config     *Config
}

var consumerPluginSave *ConsumerPlugin

func (plugin *ConsumerPlugin) Startup() error {
	rvals := plugin.Called()
	return rvals.Error(0)
}

func (plugin *ConsumerPlugin) Shutdown() error {
	rvals := plugin.Called()
	return rvals.Error(0)

}

func (plugin *ConsumerPlugin) Name() string {
	return "Dependant plugin"
}

func NewConsumerPlugin(config *Config, otherPlugin *MockPlugin) *ConsumerPlugin {
	consumerPluginSave.config = config
	consumerPluginSave.dependency = otherPlugin
	return consumerPluginSave
}

func GetConsumerPluginConstructor() (*ConsumerPlugin, *mock.Mock, func(*Config, *MockPlugin) *ConsumerPlugin) {
	consumerPluginSave = &ConsumerPlugin{}
	return consumerPluginSave, &consumerPluginSave.Mock, NewConsumerPlugin
}

// TestPluginConsumer tests that we can 'consume' plugins. That is, we
// can call RegisterPluginConsumer() to show we are interested in a
// particular consumer, and have the function provided called.
func TestPluginConsumer(t *testing.T) {
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
		pluginMock  *mock.Mock
		plugin      interface{}
		constructor interface{}
		isProvider  bool
	}
	type testConfig struct {
		plugins    []func() constructorPluginPair
		assertions func()
		consumers  []interface{}
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
			plugins: []func() constructorPluginPair{},
			assertions: func() {
				require.Equal(t, 0, len(helloConsumerPluginRegistry))
			},
			consumers: []interface{}{
				helloConsumer,
			},
		},
		{
			plugins: []func() constructorPluginPair{
				func() constructorPluginPair { return makeConstructorPluginPair(GetMockPluginConstructor()) }},
			assertions: func() {
				require.Equal(t, 0, len(helloConsumerPluginRegistry))
			},
			consumers: []interface{}{
				helloConsumer,
			},
		},
		{
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
			plugins: []func() constructorPluginPair{
				func() constructorPluginPair { return makeConstructorProviderPluginPair(GetMockPluginConstructor()) },
				func() constructorPluginPair { return makeConstructorPluginPair(GetConsumerPluginConstructor()) },
			},
			assertions: func() {
				require.Same(t, consumerPluginSave.dependency, pluginSave)
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

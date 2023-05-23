package plugins

import (
	"fmt"
	"syscall"

	"github.com/stretchr/testify/mock"
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

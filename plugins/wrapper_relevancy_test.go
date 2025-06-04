package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPlatform string

// relevancyWrapper is a wrapper which looks at specific metadata to
// decide if a plugin is relevant or not.
type relevancyWrapper struct {
	WrapperHelper
}

// a mock plugin with 'type A' so we can provide multiple plugins to
// the DI system, to allow one but not the other based on metadata.
type MockPluginA struct {
	*MockPlugin
}

// a mock plugin -- see MockPluginA.
type MockPluginB struct {
	*MockPlugin
}

// Name returns the name
func (m *MockPluginA) Name() string {
	return "MockPluginA"
}

// Name returns the name
func (m *MockPluginB) Name() string {
	return "MockPluginB"
}

func newRelevancyWrapper() *relevancyWrapper {
	w := &relevancyWrapper{}
	w.SetConstructorReturn(
		ConstructorWrapperPluginFactory(func(gen PluginGeneratorCallback, _ ...any) Plugin {
			return gen()
		}))
	return w
}

var (
	os1Constructed = false
	os2Constructed = false
)

func NewMockPluginOS1(config *Config) *MockPluginA {
	m := &MockPluginA{MockPlugin: &MockPlugin{config: config}}
	m.MockPlugin.config = config
	m.On("Startup").Maybe().Return(nil)
	m.On("Shutdown").Maybe().Return(nil)
	os1Constructed = true
	return m
}

// NewMockPluginOS2 returns an OS2 mock plugin. It is just a different
// type than MockPlugin with the same implementations.
func NewMockPluginOS2(config *Config) *MockPluginB {
	m := &MockPluginB{MockPlugin: &MockPlugin{config: config}}
	m.On("Startup").Maybe().Return(nil)
	m.On("Shutdown").Maybe().Return(nil)
	os2Constructed = true
	return m
}

// IsRelevant only returns true if the metadata slice contains a
// testPlatform string of "OS2".
func (w *relevancyWrapper) IsRelevant(val PluginConstructor, metadata ...any) bool {
	for _, m := range metadata {

		switch v := m.(type) {
		case testPlatform:
			if v == "OS2" {
				return true
			}
		}
	}
	return false
}

var _ PluginPredicate = &relevancyWrapper{}

// Test the IsRelevant method by providing a fake wrapper that decides
// something is relevant if the platform metadata is equal to "OS2".
func TestIsRelevant(t *testing.T) {
	controller := NewPluginControl()

	controller.RegisterPluginPredicate(func() *relevancyWrapper {
		return newRelevancyWrapper()
	})

	err := controller.Provide(func() *Config {
		return &Config{}
	})
	if err != nil {
		logger.Warn("Failed to provide return value: %v\n", err.Error())
	}

	controller.RegisterPlugin(NewMockPluginOS1, testPlatform("OS1"))
	controller.RegisterPlugin(NewMockPluginOS2, testPlatform("OS2"))
	controller.Startup()
	assert.True(t, os2Constructed)
	assert.False(t, os1Constructed)
}

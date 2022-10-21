package plugins

import (
	"syscall"
	"testing"
)

type MockSignalHandlingPlugin struct {
	MockPlugin
}

func TestSignalHandler(t *testing.T) {
	handler := NewSignalHandler()
	config := &Config{
		Name: "helloworld",
	}
	plugin := &MockSignalHandlingPlugin{MockPlugin{config: config}}
	plugin.On("Signal", syscall.SIGHUP).Return(nil)
	handler.RegisterPlugin(plugin)
	handler.Signal(syscall.SIGHUP)
	plugin.AssertCalled(t, "Signal", syscall.SIGHUP)
	plugin.On("Signal", syscall.SIGQUIT).Return(nil)
	handler.Signal(syscall.SIGQUIT)
	plugin.AssertCalled(t, "Signal", syscall.SIGQUIT)
}

package plugins

import (
	"syscall"
)

type SignalHandlingPlugin interface {
	Signal(syscall.Signal) error
}

type SignalHandler struct {
	plugins []SignalHandlingPlugin
}

func NewSignalHandler() *SignalHandler {
	return &SignalHandler{}
}

func (handler *SignalHandler) RegisterPlugin(plugin SignalHandlingPlugin) {
	handler.plugins = append(handler.plugins, plugin)
}

// Signal calls the Signal() method of all registered plugins with sig.
func (handler *SignalHandler) Signal(sig syscall.Signal) {
	for _, sigHandler := range handler.plugins {
		if err := sigHandler.Signal(sig); err != nil {
			plugin := sigHandler.(Plugin)
			logger.Warn("Plugin %s returned error handling signal %v: %s\n",
				plugin.Name(),
				sig,
				err)
		}
	}
}

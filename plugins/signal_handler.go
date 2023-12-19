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
	logger.Warn("--------In new SignalHandler----")
	return &SignalHandler{}
}

func (handler *SignalHandler) RegisterPlugin(plugin SignalHandlingPlugin) {
	handler.plugins = append(handler.plugins, plugin)
}

// Signal calls the Signal() method of all registered plugins with sig.
func (handler *SignalHandler) Signal(sig syscall.Signal) {
	logger.Warn("--------In sending Signal----")
	for _, sigHandler := range handler.plugins {
		logger.Warn("--------In FORLOOP for sif : %v ----", sig)
		if err := sigHandler.Signal(sig); err != nil {
			logger.Warn("--------In ERROR 1 for sig : %v ----", sig)
			plugin := sigHandler.(Plugin)
			logger.Warn("--------In ERROR 2 for sig : %v ----", sig)
			logger.Warn("Plugin %s returned error handling signal %v: %s\n",
				plugin.Name(),
				sig,
				err)
			logger.Warn("--------In ERROR 3 for sig : %v ----", sig)
		}
	}
	logger.Warn("--------Out sending Signal----")
}

package example

import (
	"github.com/untangle/golang-shared/services/logger"
)

// Startup starts the gin server
func Startup() {
	logger.Info("This example has been started\n")
}

// Shutdown function here to stop gind service
func Shutdown() {
}

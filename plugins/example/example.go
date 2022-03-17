package example

import (
	"github.com/untangle/golang-shared/services/logger"
)

// Start starts QoS
func Start() {
	logger.Info("Starting Example plugin\n")
}

// Stop stops QoS
func Stop() {
}


package arp

import (
	"github.com/untangle/discoverd/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
)

// Start starts the ARP collector
func Start() {
	logger.Info("Starting ARP collector plugin\n")
	discovery.RegisterCollector(NetlinkNeighbourCallbackController)
	// Lets do a first run to get the initial data
	NetlinkNeighbourCallbackController(nil)
}

// Stop stops QoS
func Stop() {
}

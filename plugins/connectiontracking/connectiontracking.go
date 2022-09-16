package connectiontracking

import (
	"github.com/untangle/discoverd/plugins/connectiontracking/connectiondetailer"
	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
)

type ConnnectionTracking struct {
	connectionDetails connectiondetailer.ConnectionDetailer
}

func NewConnectionTracking(connectionDetailer connectiondetailer.ConnectionDetailer) *ConnnectionTracking {
	return &ConnnectionTracking{connectionDetails: connectionDetailer}
}

// Starts the Conntrack collector
func (connectionTracking *ConnnectionTracking) Start() {
	logger.Info("Starting Conntrack collector plugin\n")
	discovery.RegisterCollector(connectionTracking.ConnectionTrackingBackHandler)

	// initial run
	connectionTracking.ConnectionTrackingBackHandler(nil)
}

// Stops Conntrack collector
func (connectionTracking *ConnnectionTracking) Stop() {
}

func (connectionTracking *ConnnectionTracking) ConnectionTrackingBackHandler(commands []discovery.Command) {
	logger.Debug("ConnectionTracking callback handler: Received %d commands\n", len(commands))

	if fetchErr := connectionTracking.connectionDetails.FetchSystemConnections(); fetchErr == nil {
		if deviceToConnections, getErr := connectionTracking.connectionDetails.GetDeviceToConnections(); getErr == nil {
			for device, connections := range deviceToConnections {
				entry := disc.DeviceEntry{}
				entry.Init()
				entry.IPv4Address = device
				entry.Connections = connections

				discovery.UpdateDiscoveryEntry("", &entry)

				logger.Debug("Created connection details for device with IPv4 address: %d\n", device)
			}
		} else {
			logger.Err("Couldn't get the connection list: %s", getErr.Error())
		}

	} else {
		logger.Err("Couldn't fetch the system's connections: %s", fetchErr.Error())
	}

}

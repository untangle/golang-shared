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

// Returns a new ConnectionTracking struct
func NewConnectionTracking(connectionDetailer connectiondetailer.ConnectionDetailer) *ConnnectionTracking {
	return &ConnnectionTracking{connectionDetails: connectionDetailer}
}

// Starts the Conntrack collector
func Start() {
	logger.Info("Starting Conntrack collector plugin\n")
	discovery.RegisterCollector(ConnectionTrackingBackHandler)

	// initial run
	ConnectionTrackingBackHandler(nil)
}

var connectionTracking *ConnnectionTracking

// Used to swap which detailer will be used to gather connection info depending on which
// OS is being used. There is only one detailer created at the moment. Update when more
// are added to differentiate between EOS and openWRT
func init() {
	connectionTracking = NewConnectionTracking(new(connectiondetailer.ConnTrackerDetails))
}

// ConnectionTrackingBackHandler is the callback handler for the connection tracker collector.
func ConnectionTrackingBackHandler(commands []discovery.Command) {
	logger.Debug("ConnectionTracking callback handler: Received %d commands\n", len(commands))

	if fetchErr := connectionTracking.connectionDetails.FetchSystemConnections(); fetchErr == nil {
		if deviceToConnections, getErr := connectionTracking.connectionDetails.GetDeviceToConnections(); getErr == nil {
			for device, connections := range deviceToConnections {
				entry := disc.DeviceEntry{}
				entry.Init()
				entry.IPv4Address = device
				entry.Connections = connections

				discovery.UpdateDiscoveryEntry("", &entry)

				logger.Debug("Created connection details for device with IPv4 address: %s\n", device)
			}
		} else {
			logger.Err("Couldn't get the connection list: %s", getErr.Error())
		}

	} else {
		logger.Err("Couldn't fetch the system's connections: %s", fetchErr.Error())
	}

}

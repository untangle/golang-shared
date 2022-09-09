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
		if connections, getErr := connectionTracking.connectionDetails.GetConnectionList(); getErr == nil {
			for _, connection := range connections {

				// Discovery entries can only be linked up by mac/ipv4
				if connection.Original.LayerThree.Protoname == "ipv4" {

					entry := disc.DeviceEntry{}
					entry.Init()
					entry.Data.ConnectionTracking = connection

					logger.Debug("Created connection for %d\n", entry.Data.ConnectionTracking.Independent.Id)

					// No mac address can be provided, so hope UpdateDiscoveryEntry can update an entry with just
					// the ipv4 address
					discovery.UpdateDiscoveryEntry("", entry)

				}

			}
		} else {
			logger.Err("Couldn't get the connection list")
		}

	} else {
		logger.Err("Couldn't fetch the system's connections")
	}

}

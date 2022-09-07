package connectiontracking

import (
	"github.com/untangle/discoverd/plugins/connectiontracking/connectiondetailer"
	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
)

// Starts the Conntrack collector
func Start() {
	logger.Info("Starting Conntrack collector plugin\n")
	discovery.RegisterCollector(ConnectionTrackingBackHandler)

	// initial run
	ConnectionTrackingBackHandler(nil)
}

// Stops Conntrack collector
func Stop() {
}


func ConnectionTrackingBackHandler(commands []discovery.Command) {
	logger.Debug("ConnectionTracking callback handler: Received %d commands\n", len(commands))

	// If EOS doesn't have ConnectionTracking swap this to a newly implemented struct fulfilling the connectinodetailer
	// interface
	var connectionDetails *connectiondetailer.ConnectionDetailer = new(connectiondetailer.ConnTrackerDetails)

	if err := connectionDetails.InitializeConnectionDetails(); err == nil {
		connDetails := connectionDetails.GetConnectionDetails()
		if len(connDetails.Flows) > 0 {
			for i := 0; i < len(conntrack.Flows); i++ {
				// initialize the discovery entry
				entry := disc.DeviceEntry{}
				entry.
			}
		} else {
			logger.Debug("No connections in conntrack nf table!\n")
		}
	} else {
		logger.Err("Unable to unmarshal xml output of the conntrack command: %s\n", err)
	}

}
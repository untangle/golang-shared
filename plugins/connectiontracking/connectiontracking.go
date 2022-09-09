package connectiontracking

import (
	"github.com/untangle/discoverd/plugins/connectiontracking/connectiondetailer"
	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
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
				if connection.Original.Layer3.Protoname == "ipv4" {

					entry := disc.DeviceEntry{}
					entry.Init()
					entry.Data.ConnectionTracking = &Discoverd.ConnectionTracking{}

					// Set IPv4 in discovery entry
					entry.Data.IPv4Address = connection.Original.Layer3.Src

					entry.Data.ConnectionTracking.Independent = &Discoverd.Independent{}
					entry.Data.ConnectionTracking.Independent.Id = connection.Independent.Id
					entry.Data.ConnectionTracking.Independent.Mark = connection.Independent.Mark
					entry.Data.ConnectionTracking.Independent.Timeout = connection.Independent.Timeout
					entry.Data.ConnectionTracking.Independent.Use = connection.Independent.Use

					entry.Data.ConnectionTracking.Original = &Discoverd.Original{}
					entry.Data.ConnectionTracking.Original.LayerThree = &Discoverd.LayerThree{}
					entry.Data.ConnectionTracking.Original.LayerFour = &Discoverd.LayerFour{}
					entry.Data.ConnectionTracking.Original.LayerThree.Protoname = connection.Original.Layer3.Protoname
					entry.Data.ConnectionTracking.Original.LayerThree.Protonum = connection.Original.Layer3.Protonum
					entry.Data.ConnectionTracking.Original.LayerThree.Src = connection.Original.Layer3.Src
					entry.Data.ConnectionTracking.Original.LayerThree.Dest = connection.Original.Layer3.Dst

					entry.Data.ConnectionTracking.Original.LayerFour.Protoname = connection.Original.Layer4.Protoname
					entry.Data.ConnectionTracking.Original.LayerFour.Protonum = connection.Original.Layer4.Protonum
					entry.Data.ConnectionTracking.Original.LayerFour.SPort = connection.Original.Layer4.SPort
					entry.Data.ConnectionTracking.Original.LayerFour.DPort = connection.Original.Layer4.DPort

					entry.Data.ConnectionTracking.Reply = &Discoverd.Reply{}
					entry.Data.ConnectionTracking.Reply.LayerThree = &Discoverd.LayerThree{}
					entry.Data.ConnectionTracking.Reply.LayerFour = &Discoverd.LayerFour{}
					entry.Data.ConnectionTracking.Reply.LayerThree.Protoname = connection.Reply.Layer3.Protoname
					entry.Data.ConnectionTracking.Reply.LayerThree.Protonum = connection.Reply.Layer3.Protonum
					entry.Data.ConnectionTracking.Reply.LayerThree.Src = connection.Reply.Layer3.Src
					entry.Data.ConnectionTracking.Reply.LayerThree.Dest = connection.Reply.Layer3.Dst

					entry.Data.ConnectionTracking.Reply.LayerFour.Protoname = connection.Reply.Layer4.Protoname
					entry.Data.ConnectionTracking.Reply.LayerFour.Protonum = connection.Reply.Layer4.Protonum
					entry.Data.ConnectionTracking.Reply.LayerFour.SPort = connection.Reply.Layer4.SPort
					entry.Data.ConnectionTracking.Reply.LayerFour.DPort = connection.Reply.Layer4.DPort

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

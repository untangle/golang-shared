package arp

import (
	"net"
	"os/exec"
	"strings"

	"github.com/untangle/discoverd/services/discovery"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// Start starts the ARP collector
func Start() {
	logger.Info("Starting ARP collector plugin\n")
	discovery.RegisterCollector(ArpcallBackHandler)
	// Lets do a first run to get the initial data
	ArpcallBackHandler(nil)
}

// Stop stops QoS
func Stop() {
}

// ArpcallBackHandler is the callback handler for the ARP collector
func ArpcallBackHandler(commands []discovery.Command) {
	logger.Debug("Arp Callback handler: Received %d commands\n", len(commands))
	cmd := exec.Command("cat", "/proc/net/arp")
	output, _ := cmd.CombinedOutput()

	// Parse each line
	for _, line := range strings.Split(string(output), "\n") {
		// Parse each field
		fields := strings.Fields(line)

		// If empty or mac address is not valid, skip
		if len(fields) == 0 || fields[3] == "00:00:00:00:00:00" {
			continue
		}

		// Initialize the entry
		entry := disc.DeviceEntry{}
		entry.Init()
		entry.Data.Arp = &Discoverd.ARP{}

		// Populate the entry
		entry.Data.Arp.Ip = fields[0]
		entry.Data.Arp.Mac = fields[3]
		entry.Data.MacAddress = entry.Data.Arp.Mac

		// Make sure the IP is valid before updating the entry, this also excludes headings
		if net.ParseIP(entry.Data.Arp.Ip) != nil {
			entry.Data.IPv4Address = entry.Data.Arp.Ip
			entry.Data.MacAddress = entry.Data.Arp.Mac
			discovery.UpdateDiscoveryEntry(entry.Data.Arp.Mac, entry)
		}
	}
}

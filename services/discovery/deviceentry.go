package discovery

import (
	"time"

	"github.com/untangle/discoverd/utils"
	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"google.golang.org/protobuf/proto"
)

// UpdateDiscoveryEntry updates the discovery list with the new entry and publishes the entry.
// If existing entry is present we update only fields that are set in the new entry
func UpdateDiscoveryEntry(mac string, entry *discovery.DeviceEntry) {
	// Check if an invalid mac was provided, not just an empty one
	// Fail if the MAC in invalid since IPv6 may be used
	if mac != "" {
		if utils.IsMacAddress(mac) {
			entry.MacAddress = mac
		} else {
			logger.Warn("UpdateDiscoveryEntry called with invalid mac: %s\n", mac)
			return
		}
	}

	// Check if either the IPv4 or MAC address was valid
	// If the IPv4 is invalid and MAC valid, publish message anyway
	// since the layer 4 protocol used might be IPv6
	if utils.IsMacAddress(mac) || utils.IsIpv4Address(entry.IPv4Address) {
		if entry.Connections != nil {
			logger.Err("This should have been published\n")
		}

		// ZMQ publish the entry
		entry.LastUpdate = time.Now().Unix()
		logger.Debug("Publishing discovery entry for %s, %s\n", mac, entry.IPv4Address)

		zmqpublishEntry(entry)

	} else {
		logger.Warn("UpdateDiscoveryEntry called with invalid IPv4 and MAC addresses\n")
	}
}

func zmqpublishEntry(entry *discovery.DeviceEntry) {
	message, err := proto.Marshal(entry)
	if err != nil {
		logger.Err("Unable to marshal discovery entry: %s\n", err)
		return
	}
	messagePublisherChannel <- &zmqMessage{"arista:discovery:device", message}
}

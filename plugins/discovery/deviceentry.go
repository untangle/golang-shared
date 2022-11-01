package discovery

import (
	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"google.golang.org/protobuf/proto"
)

// UpdateDiscoveryEntry updates the discovery list with the new entry and publishes the entry.
// If existing entry is present we update only fields that are set in the new entry
/*func UpdateDiscoveryEntry(mac string, entry *discovery.DeviceEntry) {
// Check if an invalid mac was provided, not just an empty one
// Fail if the MAC in invalid since IPv6 may be used
/*if utils.IsMacAddress(mac) {
	entry.SetMac(mac)
} else if mac != "" {
	logger.Warn("UpdateDiscoveryEntry called with invalid mac: %s\n", mac)
	return
}*/

// Check if either the IPv4 or MAC address was valid
// If the IPv4 is invalid and MAC valid, publish message anyway
// since the layer 4 protocol used might be IPv6
//if utils.IsMacAddress(mac) {
// ZMQ publish the entry
/*entry.LastUpdate = time.Now().Unix()
logger.Debug("Attempting to send discovery entry for %s to the ZMQ publisher\n", mac)

zmqpublishEntry(entry)

/*} else {
	logger.Warn("UpdateDiscoveryEntry called with invalid IPv4 and MAC addresses\n")
}*/
//}*/

func ZmqpublishEntry(entry *discovery.DeviceEntry, topic string) {
	message, err := proto.Marshal(entry)
	if err != nil {
		logger.Err("Unable to marshal discovery entry: %s\n", err)
		return
	}

	// Do not block if message can't be sent. Just log that it was dropped
	select {
	case NewDiscovery().messagePublisherChannel <- &zmqMessage{topic, message}:
		logger.Debug("Sent discovery entry to ZMQ publisher %s\n", entry.MacAddress)
	default:
		logger.Debug("Dropped discovery entry meant for the ZMQ publisher %s\n", entry.MacAddress)
	}

}

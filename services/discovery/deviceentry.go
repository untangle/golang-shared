package discovery

import (
	"time"

	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"google.golang.org/protobuf/proto"
)

// Indexed list of discovered devices. Index is the MAC address
// var deviceList = make(map[string]discovery.DeviceEntry)
// var deviceListLock sync.RWMutex = sync.RWMutex{}

// UpdateDiscoveryEntry updates the discovery list with the new entry and publishes the entry.
// If existing entry is present we update only fields that are set in the new entry
func UpdateDiscoveryEntry(mac string, entry *discovery.DeviceEntry) {

	if entry.IPv4Address == "" && mac == "" {
		logger.Warn("UpdateDiscoveryEntry called with empty IP and MAC address\n")
		return
	}

	// ZMQ publish the entry
	entry.LastUpdate = time.Now().Unix()
	logger.Debug("Publishing discovery entry for %s, %s\n", mac, entry.IPv4Address)
	zmqpublishEntry(entry)
}

func zmqpublishEntry(entry *discovery.DeviceEntry) {
	message, err := proto.Marshal(entry)
	if err != nil {
		logger.Err("Unable to marshal discovery entry: %s\n", err)
		return
	}
	messagePublisherChannel <- &zmqMessage{"arista:discovery:device", message}
}

package discovery

import (
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
	"google.golang.org/protobuf/proto"
)

// Indexed list of discovered devices. Index is the MAC address
var deviceList = make(map[string]discovery.DeviceEntry)
var deviceListLock sync.RWMutex = sync.RWMutex{}

// UpdateDiscoveryEntry updates the discovery list with the new entry and publishes the entry.
// If existing entry is present we update only fields that are set in the new entry
func UpdateDiscoveryEntry(mac string, entry discovery.DeviceEntry) {

	if entry.IPv4Address == "" && mac == "" {
		logger.Warn("UpdateDiscoveryEntry called with empty IP and MAC address\n")
		return
	}

	mac = strings.ToLower(mac)
	entry.MacAddress = mac
	logger.Debug("Received %+v\n", entry)
	// If there is no Mac address, lets see if there is an existing entry with the IP address
	if mac == "" {
		if entry.IPv4Address != "" {
			existingEntry, ok := getDeviceEntryFromIP(entry.IPv4Address)
			if ok {
				entry.MacAddress = existingEntry.MacAddress
				mac = existingEntry.MacAddress
			} else {
				logger.Warn("No entry found for IP address %s, which is missing Mac Address. Can't add\n", entry.IPv4Address)
				return
			}
		}
	}
	// Do a check to see if mac is really a Mac Address
	if !isMacAddress(mac) {
		logger.Warn("UpdateDiscoveryEntry: Invalid MAC address: %s\n", mac)
		return
	}
	deviceListLock.Lock()
	if oldEntry, ok := deviceList[mac]; ok {
		// Merge the old entry with the new one
		entry.Merge(&oldEntry)
	}
	entry.LastUpdate = time.Now().Unix()
	deviceList[mac] = entry
	deviceListLock.Unlock()

	// ZMQ publish the entry
	logger.Debug("Publishing discovery entry for %s, %s\n", mac, entry.IPv4Address)

	if entry.Connections != nil {
		logger.Err("yoyoyo: %s", entry.Connections[0].Original.LayerThree.Src)
	}

	zmqpublishEntry(entry)
}

func zmqpublishEntry(entry discovery.DeviceEntry) {
	message, err := proto.Marshal(&entry)
	if err != nil {
		logger.Err("Unable to marshal discovery entry: %s\n", err)
		return
	}
	messagePublisherChannel <- &zmqMessage{"arista:discovery:device", message}
}

func publishAll() {
	deviceListLock.RLock()
	for _, entry := range deviceList {
		zmqpublishEntry(entry)
	}
	deviceListLock.RUnlock()
}

func isMacAddress(mac string) bool {
	var validChars = "0123456789abcdefABCDEF:"
	if len(mac) != 17 {
		return false
	}
	for _, c := range mac {
		if !strings.Contains(validChars, string(c)) {
			return false
		}
	}
	return true
}

func getDeviceEntryFromIP(ip string) (discovery.DeviceEntry, bool) {
	deviceListLock.RLock()
	for _, entry := range deviceList {
		if entry.IPv4Address == ip {
			deviceListLock.RUnlock()
			return entry, true
		}
	}
	return discovery.DeviceEntry{}, false
}

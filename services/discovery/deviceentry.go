package discovery

import (
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/logger"
	disc "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

// DeviceEntry is the data structure for a device entry
type DeviceEntry struct {
	Data disc.DiscoveryEntry
}

// Indexed list of discovered devices. Index is the MAC address
var deviceList = make(map[string]DeviceEntry)
var deviceListLock sync.RWMutex = sync.RWMutex{}

// We need to merge the new entry with the existing entry. Not all data providers will
// have all data at all times, hence gathering the widest set of data possible.
func (n DeviceEntry) merge(o DeviceEntry) {
	if n.Data.IPv4Address != "" {
		o.Data.IPv4Address = n.Data.IPv4Address
	}
	if n.Data.Lldp != nil {
		o.Data.Lldp = n.Data.Lldp
	}
	if n.Data.Arp != nil {
		o.Data.Arp = n.Data.Arp
	}
	if n.Data.Nmap != nil {
		o.Data.Nmap = n.Data.Nmap
	}
}

// Init initialize a new DeviceEntry
func (n DeviceEntry) Init() {
	n.Data.IPv4Address = ""
	n.Data.MacAddress = ""
	n.Data.Lldp = nil
	n.Data.Arp = nil
	n.Data.Nmap = nil
}

// UpdateDiscoveryEntry updates the discovery list with the new entry and publishes the entry.
// If existing entry is present we update only fields that are set in the new entry
func UpdateDiscoveryEntry(mac string, entry DeviceEntry) {

	deviceListLock.Lock()
	if oldEntry, ok := deviceList[mac]; ok {
		// Merge the old entry with the new one
		entry.merge(oldEntry)
	}
	entry.Data.LastUpdate = time.Now().Unix()
	deviceList[mac] = entry
	deviceListLock.Unlock()

	// ZMQ publish the entry
	logger.Debug("Publishing discovery entry for %s\n", mac)
	zmqpublishEntry(entry)
}

func zmqpublishEntry(entry DeviceEntry) {
	message, err := proto.Marshal(&entry.Data)
	if err != nil {
		logger.Err("Unable to marshal discovery entry: %s\n", err)
		return
	}
	messageChannel <- &zmqMessage{"arista:discovery:device", message}
}

package discovery

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/ActiveSessions"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

// DeviceEntry represents a device found via discovery and methods on
// it. Mostly used as a key of DevicesList.
type DeviceEntry struct {
	disco.DiscoveryEntry
	sessions []*ActiveSessions.Session
}

// DevicesList is an in-memory 'list' of all known devices (stored as
// a map from mac address string to device).
type DevicesList struct {
	Devices map[string]*DeviceEntry

	// an index of the devices by IP.
	devicesByIP map[string]*DeviceEntry
	Lock        sync.RWMutex
}

// NewDevicesList returns a new DevicesList which is ready for use.
func NewDevicesList() *DevicesList {
	return &DevicesList{
		Devices:     map[string]*DeviceEntry{},
		devicesByIP: map[string]*DeviceEntry{},
	}
}

// ListPredicate is a function that when applied to an entry returns
// true if it 'accepts' the entry.
type ListPredicate func(entry *DeviceEntry) bool

// WithUpdatesWithinDuration returns a predicate ensuring that the
// LastUpdate member was within the period. So if period is an hour
// for example, the returned predicate will return true when the
// device has a LastUpdate time within the last hour.
func WithUpdatesWithinDuration(period time.Duration) ListPredicate {
	return func(entry *DeviceEntry) bool {
		lastUpdated := time.Unix(entry.LastUpdate, 0)
		now := time.Now()
		return (now.Sub(lastUpdated) <= period)
	}
}

// putDeviceUnsafe puts the device in the list without locking it.
func (list *DevicesList) putDeviceUnsafe(entry *DeviceEntry) {
	list.Devices[entry.MacAddress] = entry
	if entry.IPv4Address != "" {
		list.devicesByIP[entry.IPv4Address] = entry
	}
}

func (list *DevicesList) PutDevice(entry *DeviceEntry) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	list.putDeviceUnsafe(entry)

}

// Get 24hours older device discovery entry from device list and delete the entry from device list
func (list *DevicesList) CleanOldDeviceEntry(preds ...ListPredicate) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	listOfDevs := list.listDevices(preds...)
	list.CleanDevices(listOfDevs)
}

// Clean device discovery entry from devices list if the entry lastUpdate is 24 hours older
func (list *DevicesList) CleanDevices(devices []*DeviceEntry) {

	for _, device := range devices {
		delete(list.Devices, device.MacAddress)
		if device.IPv4Address != "" {
			delete(list.devicesByIP, device.IPv4Address)
		}
		logger.Debug("Deleted entry %s:%s\n", device.MacAddress, device.IPv4Address)
	}
}

// Get the device discovery entries lastUpdate time is older than the duration(24 hours)
func LastUpdateOlderThanDuration(period time.Duration) ListPredicate {
	return func(entry *DeviceEntry) bool {
		lastUpdated := time.Unix(entry.LastUpdate, 0)
		now := time.Now()
		return (now.Sub(lastUpdated) >= period)
	}

}

// listDevices returns a list of devices matching all predicates. It
// doesn't do anything with locks so without the outer function
// locking appropriately is unsafe.
func (list *DevicesList) listDevices(preds ...ListPredicate) (returns []*DeviceEntry) {

	returns = []*DeviceEntry{}
search:
	for _, device := range list.Devices {
		for _, pred := range preds {
			if !pred(device) {
				continue search
			}
		}
		returns = append(returns, device)
	}
	return
}

// ApplyToDeviceList applys doToList to the list of device entries in
// all the DeviceList device entries that match all the predicates. It
// does this with the lock taken.
func (list *DevicesList) ApplyToDeviceList(
	doToList func([]*DeviceEntry) (interface{}, error),
	preds ...ListPredicate) (interface{}, error) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	listOfDevs := list.listDevices(preds...)
	return doToList(listOfDevs)
}

// getDeviceFromIPUnsafe gets a device in the table by IP address.
func (list *DevicesList) getDeviceFromIPUnsafe(ip string) *DeviceEntry {
	return list.devicesByIP[ip]
}

// GetDeviceEntryFromIP returns a copy of a device entry in the list with IP
// address ip. *Currently only works with ipv4*.
func (list *DevicesList) GetDeviceEntryFromIP(ip string) *disco.DiscoveryEntry {
	list.Lock.RLock()
	defer list.Lock.RUnlock()
	if entry := list.getDeviceFromIPUnsafe(ip); entry != nil {
		cloned := proto.Clone(&entry.DiscoveryEntry)
		return cloned.(*disco.DiscoveryEntry)
	}
	return nil
}

// MergeOrAddDeviceEntry merges the new entry if an entry can be found
// that corresponds to the same MAC or IP. If none can be found, we
// put the new entry in the table. The provided callback function is
// called after everything is merged but before the lock is
// released. This can allow you to clone/copy the merged device.
func (list *DevicesList) MergeOrAddDeviceEntry(entry *DeviceEntry, callback func()) {
	if entry.MacAddress == "00:00:00:00:00:00" {
		return
	}
	list.Lock.Lock()
	defer list.Lock.Unlock()
	if entry.MacAddress == "" && entry.IPv4Address != "" {
		if found := list.getDeviceFromIPUnsafe(entry.IPv4Address); found != nil {
			entry.Merge(found)
		} else {
			return
		}
	} else if entry.MacAddress == "" {
		return
	} else if oldEntry, ok := list.Devices[entry.MacAddress]; ok {
		entry.Merge(oldEntry)
	}
	list.putDeviceUnsafe(entry)
	callback()
}

// MergeSessions merges the sessions into the devices.
func (list *DevicesList) MergeSessions(sessions []*ActiveSessions.Session) {
	for _, device := range list.Devices {
		device.sessions = nil
	}
	for _, session := range sessions {
		if entry, ok := list.devicesByIP[session.ClientAddress]; ok {
			entry.sessions = append(entry.sessions, session)
		}
	}
}

// Init initialize a new DeviceEntry
func (n *DeviceEntry) Init() {
	n.IPv4Address = ""
	n.MacAddress = ""
	n.Lldp = nil
	n.Arp = nil
	n.Nmap = nil
}

// Merge fills the relevant fields of n that are not present with ones
// of newEntry that are.
func (n *DeviceEntry) Merge(newEntry *DeviceEntry) {
	if n.IPv4Address == "" {
		n.IPv4Address = newEntry.IPv4Address
	}
	if n.MacAddress == "" {
		n.MacAddress = newEntry.MacAddress
	}
	if n.Lldp == nil {
		n.Lldp = newEntry.Lldp
	}
	if n.Arp == nil {
		n.Arp = newEntry.Arp
	}
	if n.Nmap == nil {
		n.Nmap = newEntry.Nmap
	}
	if n.LastUpdate < newEntry.LastUpdate {
		n.LastUpdate = newEntry.LastUpdate
	}

}

// SessionDetail is a summary of active session details for a device.
type SessionDetail struct {
	// Total byte transfer rate of all sessions.
	ByteTransferRate int64 `json:"byteTransferRate"`
	// Total number of active sessions on this device.
	NumSessions int64 `json:"numSessions"`
}

func (n *DeviceEntry) calcSessionDetails() (output SessionDetail) {
	for _, session := range n.sessions {
		output.ByteTransferRate += int64(session.ByteRate)
		output.NumSessions++
	}
	return
}

func (n *DeviceEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		*disco.DiscoveryEntry
		SessionDetail SessionDetail `json:"sessionDetail"`
	}{
		DiscoveryEntry: &n.DiscoveryEntry,
		SessionDetail:  n.calcSessionDetails(),
	})
}

// SetMac sets the mac address of the device entry. It 'normalizes' it.
func (n *DeviceEntry) SetMac(mac string) {
	n.MacAddress = strings.ToLower(mac)

}

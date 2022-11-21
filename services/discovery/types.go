package discovery

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/untangle/golang-shared/structs/protocolbuffers/ActiveSessions"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

// DeviceEntry represents a device found via discovery and methods on
// it. Mostly used as a key of DevicesList.
type DeviceEntry struct {
	disco.DiscoveryEntry
	sessions []*ActiveSessions.Session
	rxTotal  uint
	txTotal  uint
}

// DevicesList is an in-memory 'list' of all known devices (stored as
// a map from mac address string to device).
type DevicesList struct {
	Devices map[string]*DeviceEntry

	// TODO: Prior to 1.18, net.HardwareAddr and net.IpAddr cannot be used as a map key.
	// Currently we have to be careful about making our Mac and IPV6 string uppercase
	// before using them as a key in devicesByIp and Devices maps to avoid duplicates. Once upgraded
	// to 1.18+, swap these string type keys to net.HardwareAddr for Devices and net.IpAddr for devicesByIp
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

	for _, ip := range entry.getDeviceIpsUnsafe() {
		list.devicesByIP[ip] = entry
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
		for _, ip := range device.getDeviceIpsUnsafe() {
			delete(list.devicesByIP, ip)
		}
		logger.Debug("Deleted entry %s:%s\n", device.MacAddress, device.MacAddress)
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

// ApplyToDeviceWithMac will apply doToDev to the device with the
// given mac, if present, and return the result, otherwise it will
// return an error.
func (list *DevicesList) ApplyToDeviceWithMac(
	doToDev func(*DeviceEntry) (interface{}, error),
	mac string) (interface{}, error) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	if dev, wasFound := list.Devices[mac]; wasFound {
		return doToDev(dev)
	}
	return nil, fmt.Errorf("device with mac: %s not found", mac)
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
// Make sure to merge new into old.
func (list *DevicesList) MergeOrAddDeviceEntry(entry *DeviceEntry, callback func()) {
	// Lock the entry down before reading from it.
	// Otherwise the read in Merge causes a data race
	if entry.MacAddress == "00:00:00:00:00:00" {
		return
	}

	list.Lock.Lock()
	defer list.Lock.Unlock()

	deviceIps := entry.getDeviceIpsUnsafe()
	if entry.MacAddress == "" && len(deviceIps) <= 0 {
		return
	} else if oldEntry, ok := list.Devices[entry.MacAddress]; ok {
		oldEntry.Merge(entry)
		list.putDeviceUnsafe(oldEntry)
	} else if len(deviceIps) > 0 {
		// See if the IPs of entry correspond to any others
		found := false
		for _, ip := range deviceIps {
			// Once an old entry is oldEntry and the new entry is merged with it,
			// break out of the loop since any device oldEntry is a pointer that
			// every IP for a device points to
			if oldEntry := list.getDeviceFromIPUnsafe(ip); oldEntry != nil {
				oldEntry.Merge(entry)
				list.putDeviceUnsafe(oldEntry)
				found = true
				break
			}
		}
		if !found {
			list.putDeviceUnsafe(entry)
		}
	} else {
		list.putDeviceUnsafe(entry)
	}

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
	n.MacAddress = ""
	n.Lldp = nil
	n.Neigh = nil
	n.Nmap = nil
}

// Returns the list of IPs being used by a device. Does not acquire any locks
// before accessing device list elements. The IPs are fetched by going through
// each collector entry and adding any IPs found to a set
func (n *DeviceEntry) getDeviceIpsUnsafe() []string {
	// Use a set to easily get the list of unique IPs assigned to a device
	ipSet := make(map[string]struct{})

	for ip := range n.Neigh {
		ipSet[ip] = struct{}{}
	}

	for ip := range n.Lldp {
		ipSet[ip] = struct{}{}
	}

	for ip := range n.Nmap {
		ipSet[ip] = struct{}{}
	}

	var ipList []string
	for ip := range ipSet {
		ipList = append(ipList, ip)
	}

	return ipList
}

// Merge fills the relevant fields of n that are not present with ones
// of newEntry that are.
func (n *DeviceEntry) Merge(newEntry *DeviceEntry) {
	// The protobuf library has a merge function that merges exactly as needed,
	// except for the case where the LastUpdated time coming in is less than
	// The current LastUpdated time. Take a snapshot of the original before merging
	lastUpdate := n.LastUpdate

	proto.Merge(n, newEntry)

	if lastUpdate > n.LastUpdate {
		n.LastUpdate = lastUpdate
	}
}

// SessionDetail is a summary of active session details for a device.
type SessionDetail struct {
	// Total byte transfer rate of all sessions.
	ByteTransferRate int64 `json:"byteTransferRate"`

	// Total number of active sessions on this device.
	NumSessions int64 `json:"numSessions"`

	// Total number of bytes used by a device
	DataUsage int64 `json:"dataUsage"`

	RxTotal int64 `json:"rxTotal"`

	TxTotal int64 `json:"txTotal"`
}

func (n *DeviceEntry) calcSessionDetails() (output SessionDetail) {
	output.DataUsage = int64(n.rxTotal + n.txTotal)
	output.RxTotal = int64(n.rxTotal)
	output.TxTotal = int64(n.txTotal)
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

// IncrTx increments total tx bytes by tx, returns updated total.
func (n *DeviceEntry) IncrTx(tx uint) uint {
	n.txTotal += tx
	return n.txTotal
}

// IncrRx increments total rx bytes by rx, returns updated total.
func (n *DeviceEntry) IncrRx(rx uint) uint {
	n.rxTotal += rx
	return n.rxTotal
}

// RxTotal returns total rx bytes for this device.
func (n *DeviceEntry) RxTotal() uint {
	return n.rxTotal
}

// TxTotal returns total tx bytes for this device
func (n *DeviceEntry) TxTotal() uint {
	return n.txTotal
}

// SetMac sets the mac address of the device entry. It 'normalizes' it.
func (n *DeviceEntry) SetMac(mac string) {
	n.MacAddress = strings.ToLower(mac)

}

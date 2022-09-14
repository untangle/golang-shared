package discovery

import (
	"fmt"
	"sync"
	"time"

	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

type DeviceEntry struct {
	disco.DiscoveryEntry
}

// DevicesList is an in-memory 'list' of all known devices (stored as
// a map from mac address string to device).
type DevicesList struct {
	Devices map[string]*DeviceEntry
	Lock    sync.RWMutex
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

func (list *DevicesList) PutDevice(entry *DeviceEntry) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	list.Devices[entry.MacAddress] = entry
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
	for _, entry := range list.Devices {
		if entry.IPv4Address == ip {
			return entry
		}
	}
	return nil
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
	list.Lock.Lock()
	defer list.Lock.Unlock()
	if entry.MacAddress == "" && entry.IPv4Address != "" {
		if found := list.getDeviceFromIPUnsafe(entry.IPv4Address); found != nil {
			entry.Merge(found)
		} else {
			fmt.Printf("RRR: %+v\n", entry)
			fmt.Printf("L: %+#v\n", list)
			for k, v := range list.Devices {
				fmt.Printf("%s: %+v\n", k, v)
			}
			return
		}
	} else if entry.MacAddress == "" {
		fmt.Printf("REarly: %+v\n", entry)
		return
	} else if oldEntry, ok := list.Devices[entry.MacAddress]; ok {
		entry.Merge(oldEntry)
	}
	list.Devices[entry.MacAddress] = entry
	callback()
}

// Init initialize a new DeviceEntry
func (n *DeviceEntry) Init() {
	n.IPv4Address = ""
	n.MacAddress = ""
	n.Lldp = nil
	n.Arp = nil
	n.Nmap = nil
}

func (n *DeviceEntry) Merge(o *DeviceEntry) {
	if n.IPv4Address == "" {
		n.IPv4Address = o.IPv4Address
	}
	if n.MacAddress == "" {
		n.MacAddress = o.MacAddress
	}
	if n.Lldp == nil {
		n.Lldp = o.Lldp
	}
	if n.Arp == nil {
		n.Arp = o.Arp
	}
	if n.Nmap == nil {
		n.Nmap = o.Nmap
	}
	if n.LastUpdate < o.LastUpdate {
		n.LastUpdate = o.LastUpdate
	}
}

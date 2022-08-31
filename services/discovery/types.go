package discovery

import (
	"net"
	"sync"
	"time"

	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type DeviceEntry struct {
	Data disco.DiscoveryEntry
}

type DevicesList struct {
	Devices map[string]DeviceEntry
	Lock    sync.RWMutex
}

type ListPredicate func(entry *DeviceEntry) bool

func WithUpdatesWithinDuration(period time.Duration) ListPredicate {
	return func(entry *DeviceEntry) bool {
		lastUpdated := time.Unix(entry.Data.LastUpdate, 0)
		now := time.Now()
		return (now.Sub(lastUpdated) <= period)
	}
}

func IsNotFromLocalInterface(osListedInterfaces []net.Interface) ListPredicate {
	mapOfMacs := map[string]*net.Interface{}
	for _, intf := range osListedInterfaces {
		mapOfMacs[intf.HardwareAddr.String()] = &intf
	}
	return func(entry *DeviceEntry) bool {
		_, isMACFromThisMachine := mapOfMacs[entry.Data.MacAddress]
		return !isMACFromThisMachine
	}
}

func (list *DevicesList) ListDevices(preds ...ListPredicate) (returns []DeviceEntry) {
	list.Lock.RLock()
	defer list.Lock.RUnlock()
	returns = []DeviceEntry{}
search:
	for _, device := range list.Devices {
		for _, pred := range preds {
			if !pred(&device) {
				continue search
			}
		}
		returns = append(returns, device)
	}
	return
}

func (list *DevicesList) GetDeviceEntryFromIP(ip string) *disco.DiscoveryEntry {
	list.Lock.RLock()
	defer list.Lock.RUnlock()

	for _, entry := range list.Devices {
		if entry.Data.IPv4Address == ip {
			return &entry.Data
		}
	}

	return nil
}

// Init initialize a new DeviceEntry
func (n *DeviceEntry) Init() {
	n.Data.IPv4Address = ""
	n.Data.MacAddress = ""
	n.Data.Lldp = nil
	n.Data.Arp = nil
	n.Data.Nmap = nil
}

func (n *DeviceEntry) Merge(o DeviceEntry) {
	if n.Data.IPv4Address == "" {
		n.Data.IPv4Address = o.Data.IPv4Address
	}
	if n.Data.Lldp == nil {
		n.Data.Lldp = o.Data.Lldp
	}
	if n.Data.Arp == nil {
		n.Data.Arp = o.Data.Arp
	}
	if n.Data.Nmap == nil {
		n.Data.Nmap = o.Data.Nmap
	}
}

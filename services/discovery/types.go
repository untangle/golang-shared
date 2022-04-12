package discovery

import (
	"sync"

	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type DeviceEntry struct {
	Data disco.DiscoveryEntry
}

type DevicesList struct {
	Devices map[string]DeviceEntry
	Lock    sync.RWMutex
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

func (n *DeviceEntry) Merge(o DeviceEntry) {
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
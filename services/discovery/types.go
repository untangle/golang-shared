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

	if entry, ok := list.Devices[ip]; ok {
		return &entry.Data
	}
	return nil
}

package discovery

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	mfw_ifaces "github.com/untangle/golang-shared/util/net"
	"google.golang.org/protobuf/proto"
)

type DeviceEntry struct {
	disco.DiscoveryEntry
}

type DevicesList struct {
	Devices map[string]*DeviceEntry
	Lock    sync.RWMutex
}

type SystemNetInterface interface {
	Addrs() ([]net.Addr, error)
	GetName() string
}

type SystemInterfaceProxy net.Interface

func (iface *SystemInterfaceProxy) GetName() string {
	return iface.Name
}

type ListPredicate func(entry *DeviceEntry) bool

func WithUpdatesWithinDuration(period time.Duration) ListPredicate {
	return func(entry *DeviceEntry) bool {
		lastUpdated := time.Unix(entry.LastUpdate, 0)
		now := time.Now()
		return (now.Sub(lastUpdated) <= period)
	}
}

func IsNotFromLocalInterface(osListedInterfaces []net.Interface) ListPredicate {
	mapOfMacs := map[string]*net.Interface{}
	for _, intf := range osListedInterfaces {
		mapOfMacs[strings.ToUpper(intf.HardwareAddr.String())] = &intf
	}
	fmt.Printf("map of macs: %#v\n", mapOfMacs)
	return func(entry *DeviceEntry) bool {
		fmt.Printf("getting entry for: %s", strings.ToUpper(entry.MacAddress))
		_, isMACFromThisMachine := mapOfMacs[strings.ToUpper(entry.MacAddress)]
		return !isMACFromThisMachine
	}
}

func IsNotFromWANDevice(allInterfaces []*mfw_ifaces.Interface, osListedInetfaces []SystemNetInterface) ListPredicate {
	//bannedNetworks := []net.IPNet{}
	return func(entry *DeviceEntry) bool {
		return false
	}
}

func (list *DevicesList) PutDevice(entry *DeviceEntry) {
	list.Lock.Lock()
	defer list.Lock.Unlock()
	list.Devices[entry.MacAddress] = entry
}

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
// does this with the read lock taken.
func (list *DevicesList) ApplyToDeviceList(
	doToList func([]*DeviceEntry) (interface{}, error),
	preds ...ListPredicate) (interface{}, error) {
	list.Lock.RLock()
	defer list.Lock.RUnlock()
	listOfDevs := list.listDevices(preds...)
	return doToList(listOfDevs)
}

func (list *DevicesList) GetDeviceEntryFromIP(ip string) *disco.DiscoveryEntry {
	list.Lock.RLock()
	defer list.Lock.RUnlock()

	for _, entry := range list.Devices {
		if entry.IPv4Address == ip {
			cloned := proto.Clone(&entry.DiscoveryEntry)
			return cloned.(*disco.DiscoveryEntry)
		}
	}
	return nil
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
	if n.Lldp == nil {
		n.Lldp = o.Lldp
	}
	if n.Arp == nil {
		n.Arp = o.Arp
	}
	if n.Nmap == nil {
		n.Nmap = o.Nmap
	}
}

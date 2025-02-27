package discovery

import (
	"strings"

	logService "github.com/untangle/golang-shared/services/logger"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

var logger = logService.GetLoggerInstance()

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

// MergeZmqMessageIntoDeviceList merges a zmq message into the device
// list, by unmarshalling it, putting it into a DeviceEntry, and
// merging it. callback is called with a clone of the entry merged
// into the dictionary (it needs to be a clone otherwise a race might
// occur).
func MergeZmqMessageIntoDeviceList(devlist *DevicesList, device *DeviceEntry, callback func(*DeviceEntry)) error {
	clonedEntry := &DeviceEntry{}
	devlist.MergeOrAddDeviceEntry(
		device,
		func() {
			proto.Merge(&clonedEntry.DiscoveryEntry, device)
			callback(clonedEntry)
		},
	)
	return nil
}

// checkStaleNeigh would check if the encoming neighbour is stale and already marked stale in the deviceList
func checkStaleNeigh(devices map[string]*DeviceEntry, neighDevice map[string]*disco.NEIGH, macAddress string) bool {
	oldEntry, ok := devices[macAddress]
	if !ok {
		logger.Debug("Error reading device list\n")
		return false
	}
	var knownNeighState, newNeighState *disco.NEIGH
	for ip := range neighDevice {
		knownNeighState = oldEntry.Neigh[ip]
		newNeighState = neighDevice[ip]
	}
	if knownNeighState != nil && newNeighState.State == "STALE" && knownNeighState.State == "STALE" {
		return true
	}
	return false
}

// FillDeviceListWithZMQDeviceMessages will run an infinite loop
// receiving messages from channel and putting them into the device
// list. Call callback on each new device.
func FillDeviceListWithZMQDeviceMessages(
	devlist *DevicesList,
	channel chan *ZmqMessage,
	shutdownChannel chan bool,
	callback func(*DeviceEntry)) {
	for {
		select {
		case msg := <-channel:
			// The nicest way to merge collector messages into device entries
			// is to make a new device entry and add the collector message to it.
			// Then, use the protobuf library's merge function
			switch msg.Topic {
			case LLDPDeviceZMQTopic:
				lldp := &disco.LLDP{}
				logger.Info("Attempting to unmarshall LLDP ZMQ message\n")
				if err := proto.Unmarshal(msg.Message, lldp); err != nil {
					logger.Warn("Could not unmarshal LLDP ZMQ Message: %s\n", err.Error())
					break
				}
				if strings.Contains(lldp.Interface, "ma") {
					// skipping Management interface
					break
				}

				macAddress := strings.ToUpper(lldp.Mac)
				ipAddr := strings.ToUpper(lldp.Ip)

				entryMap := make(map[string]*disco.LLDP)
				entryMap[ipAddr] = lldp

				lldpDeviceEntry := &DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{Lldp: entryMap,
					MacAddress: macAddress,
					LastUpdate: lldp.LastUpdate}}
				if err := MergeZmqMessageIntoDeviceList(devlist, lldpDeviceEntry, callback); err != nil {
					logger.Warn("Could not process LLDP ZMQ message: %\n", err.Error())
				}
			case NEIGHDeviceZMQTopic:
				neigh := &disco.NEIGH{}
				logger.Info("Attempting to unmarshall NEIGH ZMQ message\n")
				if err := proto.Unmarshal(msg.Message, neigh); err != nil {
					logger.Warn("Could not unmarshal NEIGH ZMQ Message: %s\n", err.Error())
					break
				}
				if strings.Contains(neigh.Interface, "ma") {
					// skipping Management interfaces
					break
				}

				macAddress := strings.ToUpper(neigh.Mac)
				ipAddr := strings.ToUpper(neigh.Ip)

				entryMap := make(map[string]*disco.NEIGH)
				entryMap[ipAddr] = neigh

				neighDeviceEntry := &DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{Neigh: entryMap,
					MacAddress: macAddress,
					LastUpdate: neigh.LastUpdate}}
				if !checkStaleNeigh(devlist.Devices, entryMap, macAddress) {
					if err := MergeZmqMessageIntoDeviceList(devlist, neighDeviceEntry, callback); err != nil {
						logger.Warn("Could not process NEIGH ZMQ message: %\n", err.Error())
					}
				}
			case NMAPDeviceZMQTopic:
				nmap := &disco.NMAP{}
				logger.Info("Attempting to unmarshall NMAP ZMQ message\n")
				if err := proto.Unmarshal(msg.Message, nmap); err != nil {
					logger.Warn("Could not unmarshal NMAP ZMQ Message: %s\n", err.Error())
					break
				}
				if strings.Contains(nmap.Interface, "ma") {
					// skipping Management interfaces
					break
				}

				macAddress := strings.ToUpper(nmap.Mac)
				ipAddr := strings.ToUpper(nmap.Ip)

				entryMap := make(map[string]*disco.NMAP)
				entryMap[ipAddr] = nmap

				nmapDeviceEntry := &DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{Nmap: entryMap,
					MacAddress: macAddress,
					LastUpdate: nmap.LastUpdate}}
				if err := MergeZmqMessageIntoDeviceList(devlist, nmapDeviceEntry, callback); err != nil {
					logger.Warn("Could not process NMAP ZMQ message: %\n", err.Error())
				}
			}
		case <-shutdownChannel:
			shutdownChannel <- true // for testing purposes
			return
		}
	}
}

package discovery

import (
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
	devlist.MergeOrAddDeviceEntry(device,
		func() {
			proto.Merge(&clonedEntry.DiscoveryEntry, device)
			callback(clonedEntry)
		})
	return nil
}

// TODO: When mergiing Csaba's changes in, make sure to add a shutdown in reportd for this function

// FillDeviceListWithZMQDeviceMessages will run an infinite loop
// receiving messages from channel and putting them into the device
// list. Call callback on each new device.
func FillDeviceListWithZMQDeviceMessages(
	devlist *DevicesList,
	channel chan *ZmqMessage,
	shutdownChannel chan bool,
	callback func(*DeviceEntry)) {
Out:
	for {
		select {
		case msg := <-channel:
			// The nicest way to merge collector messages into device entries
			// is to make a new device entry and add the collector message to it.
			// Then, use the protobuf library's merge function
			switch msg.Topic {
			case LLDPDeviceZMQTopic:
				lldp := &disco.LLDP{}
				if err := proto.Unmarshal(msg.Message, lldp); err != nil {
					logger.Warn("Could not unmarshal LLDP ZMQ Message: %s", err.Error())
					break
				}

				lldpDeviceEntry := &DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{Lldp: []*disco.LLDP{lldp}, MacAddress: lldp.Mac}}
				if err := MergeZmqMessageIntoDeviceList(devlist, lldpDeviceEntry, callback); err != nil {
					logger.Warn("Could not process LLDP ZMQ message: %\n", err.Error())
				}
			case NEIGHDeviceZMQTopic:
				neigh := &disco.NEIGH{}
				if err := proto.Unmarshal(msg.Message, neigh); err != nil {
					logger.Warn("Could not unmarshal NEIGH ZMQ Message: %s", err.Error())
					break
				}

				neighDeviceEntry := &DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{Neigh: []*disco.NEIGH{neigh}, MacAddress: neigh.Mac}}
				if err := MergeZmqMessageIntoDeviceList(devlist, neighDeviceEntry, callback); err != nil {
					logger.Warn("Could not process NEIGH ZMQ message: %\n", err.Error())
				}
			case NMAPDeviceZMQTopic:
				nmap := &disco.NMAP{}
				if err := proto.Unmarshal(msg.Message, nmap); err != nil {
					logger.Warn("Could not unmarshal NMAP ZMQ Message: %s", err.Error())
					break
				}

				nmapDeviceEntry := &DeviceEntry{DiscoveryEntry: disco.DiscoveryEntry{Nmap: []*disco.NMAP{nmap}, MacAddress: nmap.Mac}}
				if err := MergeZmqMessageIntoDeviceList(devlist, nmapDeviceEntry, callback); err != nil {
					logger.Warn("Could not process NMAP ZMQ message: %\n", err.Error())
				}
			}
		case <-shutdownChannel:
			shutdownChannel <- true
			break Out
		}
	}
}

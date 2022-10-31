package discovery

import (
	"github.com/untangle/golang-shared/services/logger"
	"google.golang.org/protobuf/proto"
)

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
func MergeZmqMessageIntoDeviceList(devlist *DevicesList, msg *ZmqMessage, callback func(*DeviceEntry)) error {
	device := &DeviceEntry{}
	if err := proto.Unmarshal(msg.Message, &device.DiscoveryEntry); err != nil {
		return err
	}
	clonedEntry := &DeviceEntry{}
	devlist.MergeOrAddDeviceEntry(device,
		func() {
			proto.Merge(&clonedEntry.DiscoveryEntry, device)
			callback(clonedEntry)
		})
	return nil
}

// FillDeviceListWithZMQDeviceMessages will run an infinite loop
// receiving messages from channel and putting them into the device
// list. Call callback on each new device.
func FillDeviceListWithZMQDeviceMessages(
	devlist *DevicesList,
	channel chan *ZmqMessage,
	callback func(*DeviceEntry)) {
	for {
		msg := <-channel
		logger.Debug("Received ZMQ message: %s\n", msg.Topic)
		if err := MergeZmqMessageIntoDeviceList(devlist, msg, callback); err != nil {
			logger.Warn("Couldn't process device ZMQ message: %s\n", err)
		}
	}
}

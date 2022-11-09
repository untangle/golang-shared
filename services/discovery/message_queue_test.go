package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"google.golang.org/protobuf/proto"
)

// Test if messages read by FillDeviceListWithZMQDeviceMessages() get placed on a device table.
func TestFillDeviceListWithZMQDeviceMessages(t *testing.T) {
	deviceList := NewDevicesList()
	zmqChan := make(chan *ZmqMessage)
	shutdownChannel := make(chan bool)

	// Start processing ZMQ messages
	go FillDeviceListWithZMQDeviceMessages(deviceList, zmqChan, shutdownChannel, func(de *DeviceEntry) {})

	totalSentMessages := 3
	lldpMessage, _ := proto.Marshal(&disco.LLDP{Mac: "11:11:11:11:11:11"})
	neighMessage, _ := proto.Marshal(&disco.NEIGH{Mac: "22:22:22:22:22:22"})
	nmapMessage, _ := proto.Marshal(&disco.NMAP{Mac: "33:33:33:33:33:33"})

	zmqChan <- &ZmqMessage{Topic: LLDPDeviceZMQTopic, Message: lldpMessage}
	zmqChan <- &ZmqMessage{Topic: NEIGHDeviceZMQTopic, Message: neighMessage}
	zmqChan <- &ZmqMessage{Topic: NMAPDeviceZMQTopic, Message: nmapMessage}

	// Sleep to give the ZMQ processor a change to process the sent messages
	time.Sleep(1 * time.Second)

	// Get all device entries. Acquire and release the lock to prevent data races.
	deviceList.Lock.Lock()
	allDevices := deviceList.listDevices(func(entry *DeviceEntry) bool { return true })
	deviceList.Lock.Unlock()
	assert.Equal(t, totalSentMessages, len(allDevices))

	shutdownChannel <- true
	shutdownSuccess := false
	select {
	case <-shutdownChannel:
		shutdownSuccess = true
	case <-time.After(5 * time.Second):
	}
	assert.True(t, shutdownSuccess, "The goroutine processing ZMQ messages never shut down\n")
}

package events

import (
	"fmt"
	"testing"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"github.com/untangle/golang-shared/testing/mocks"

	"google.golang.org/protobuf/proto"
)

const testSocketAddress = "inproc://testInProcDescriptor"
const testSocketAddress1 = "inproc://testInProcDescriptor1"

// TestEventPublisher tests the functionality of the ZmqEventPublisher.
// It sets up a test event handler, creates a test subscriber socket, and tests
// sending both ZmqMessage and Event messages.
// It verifies that the messages are received correctly by the subscriber socket.
func TestEventPublisher(t *testing.T) {
	// Set up event Publisher
	handler := newTestEventHandler()
	startupErr := handler.Startup()
	assert.Nil(t, startupErr, "Failed to startup Zmq publisher in background")

	subscriberSocket, err := createTestSubscriberSocket(testSocketAddress)
	assert.Nil(t, err)
	defer subscriberSocket.Close() // Close socket in case of error

	// test ZmqMessage send functionality
	t.Run("send message", func(t *testing.T) {
		testZmqMessageParam := ZmqMessage{
			Topic:   AlertZMQTopic,
			Message: []byte(`{"message": "test message"}`),
		}

		// Send ZmqMessage
		handler.SendZmqMessage(&testZmqMessageParam)
		resultMessage, resultTopic, err := receiveSentZmqMessage(subscriberSocket)

		assert.Nil(t, err)
		assert.Equal(t, AlertZMQTopic, resultTopic)
		assert.Equal(t, testZmqMessageParam.Message, resultMessage)
	})

	// test Event send functionality
	t.Run("send message", func(t *testing.T) {
		testAlertParam := Alerts.Alert{
			Type:     Alerts.AlertType_USER,
			Severity: Alerts.AlertSeverity_INFO,
			Message:  "test message",
		}

		// Send Event
		handler.Send(&testAlertParam)
		resultMessage, resultTopic, err := receiveSentEventMessage(subscriberSocket)

		assert.Nil(t, err)
		assert.Equal(t, AlertZMQTopic, resultTopic)
		assert.Equal(t, testAlertParam.GetType(), resultMessage.GetType())
		assert.Equal(t, testAlertParam.GetSeverity(), resultMessage.GetSeverity())
		assert.Equal(t, testAlertParam.GetMessage(), resultMessage.GetMessage())
	})

	// Tear down
	if shutdownErr := handler.Shutdown(); shutdownErr != nil {
		fmt.Printf("Failed to stop the goroutine running the ZMQ subscriber %v\n", shutdownErr.Error())
	}
}

// TestSubscriberNotReady tests the functionality of the ZmqEventPublisher.
// It sets up a test event handler, creates a test subscriber socket, and tests
// sending both ZmqMessage and Event messages when Subscriber is not running.
func TestSubscriberNotReady(t *testing.T) {
	// Set up event Publisher
	handler := newTestEventHandler()
	startupErr := handler.Startup()
	assert.Nil(t, startupErr, "Failed to startup Zmq publisher in background")

	subscriberSocket, err := createTestSubscriberSocket(testSocketAddress1)
	assert.Nil(t, err)
	defer subscriberSocket.Close() // Close socket in case of error

	// test ZmqMessage send functionality
	t.Run("send message", func(t *testing.T) {
		testZmqMessageParam := ZmqMessage{
			Topic:   AlertZMQTopic,
			Message: []byte(`{"message": "test message"}`),
		}

		// Send ZmqMessage
		handler.SendZmqMessage(&testZmqMessageParam)
		time.Sleep(time.Second * 1)
	})

	// test Event send functionality
	t.Run("send message", func(t *testing.T) {
		testAlertParam := Alerts.Alert{
			Type:     Alerts.AlertType_USER,
			Severity: Alerts.AlertSeverity_INFO,
			Message:  "test message",
		}

		// Send Event
		handler.Send(&testAlertParam)
		time.Sleep(time.Second * 1)
	})

	// Tear down
	if shutdownErr := handler.Shutdown(); shutdownErr != nil {
		fmt.Printf("Failed to stop the goroutine running the ZMQ subscriber %v\n", shutdownErr.Error())
	}
}

// receiveSentZmqMessage receives a sent ZMQ message from the subscriber socket.
// It returns the message, topic, and any error that occurred during reception.
func receiveSentZmqMessage(subscriberSocket *zmq.Socket) ([]byte, string, error) {
	msg, err := subscriberSocket.RecvMessageBytes(0)
	if err != nil {
		return []byte(""), "", err
	}

	resultTopic := string(msg[0])
	resultMessage := msg[1]

	return resultMessage, resultTopic, nil
}

// receiveSentEventMessage receives a sent Event message from the subscriber socket.
// It returns the Event message, topic, and any error that occurred during reception.
// The Event message is unmarshaled from the received bytes using protocol buffers.
func receiveSentEventMessage(subscriberSocket *zmq.Socket) (*Alerts.Alert, string, error) {
	msg, err := subscriberSocket.RecvMessageBytes(0)
	if err != nil {
		return nil, "", err
	}

	resultTopic := string(msg[0])
	resultMessage := &Alerts.Alert{}

	err = proto.Unmarshal(msg[1], resultMessage)
	if err != nil {
		return nil, "", err
	}

	return resultMessage, resultTopic, nil
}

// createTestSubscriberSocket creates a new subscriber socket for testing.
// It binds the socket to the specified address and sets the subscription to the AlertZMQTopic.
func createTestSubscriberSocket(socket string) (*zmq.Socket, error) {
	subSocket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}

	if err = subSocket.Bind(socket); err != nil {
		return nil, err
	}

	if err = subSocket.SetSubscribe(AlertZMQTopic); err != nil {
		return nil, err
	}

	return subSocket, nil
}

// newTestEventHandler creates a new test event handler for testing.
// It returns a new ZmqEventPublisher instance with a mock logger and test socket address.
func newTestEventHandler() *ZmqEventPublisher {
	return &ZmqEventPublisher{
		logger:                  mocks.NewMockLogger(),
		messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
		zmqPublisherShutdown:    make(chan bool),
		zmqPublisherStarted:     make(chan int32, 1),
		socketAddress:           testSocketAddress,
	}
}

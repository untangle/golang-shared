package events

import (
	"fmt"
	"testing"

	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/testing/mocks"
)

const testSocketAddress = "inproc://testInProcDescriptor"

func TestEventPublisher(t *testing.T) {
	//Set up
	handler := newTestEventHandler()
	startupErr := handler.Startup()
	assert.Nil(t, startupErr, "Failed to startup Zmq publisher in background")

	subscriberSocket, err := createTestSubscriberSocket(testSocketAddress)
	assert.Nil(t, err)
	defer subscriberSocket.Close() // Close socket in case of error

	t.Run("send message", func(t *testing.T) {
		testParam := ZmqMessage{
			Topic:   AlertZMQTopic,
			Message: []byte(`{"message": "test message"}`),
		}

		// Send Event
		handler.Send(&testParam)
		resultMessage, resultTopic, err := receiveSentMessage(subscriberSocket)

		assert.Nil(t, err)
		assert.Equal(t, AlertZMQTopic, resultTopic)
		assert.Equal(t, testParam.Message, resultMessage)
	})

	// Tear down
	if shutdownErr := handler.Shutdown(); shutdownErr != nil {
		fmt.Printf("Failed to stop the goroutine running the ZMQ subscriber %v\n", shutdownErr.Error())
	}
}

func receiveSentMessage(subscriberSocket *zmq.Socket) ([]byte, string, error) {
	msg, err := subscriberSocket.RecvMessageBytes(0)
	if err != nil {
		return []byte(""), "", err
	}

	resultTopic := string(msg[0])
	resultMessage := msg[1]

	return resultMessage, resultTopic, nil
}

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

func newTestEventHandler() *ZmqEventPublisher {
	return &ZmqEventPublisher{
		logger:                  mocks.NewMockLogger(),
		messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
		zmqPublisherShutdown:    make(chan bool),
		zmqPublisherStarted:     make(chan int32, 1),
		socketAddress:           testSocketAddress,
	}
}

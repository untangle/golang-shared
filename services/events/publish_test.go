package events

import (
	"encoding/json"
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

	t.Run("send message", func(t *testing.T) {
		testParam := ZmqMessage{
			Topic:   ZMQTopicRestdEvents,
			Message: []byte(`{"type": "user", "severity": "info", "message": "test message"}`),
		}

		// Send Event
		handler.Send(&testParam)
		resultMessage, resultTopic, err := receiveSentMessage(subscriberSocket)

		assert.Nil(t, err)
		assert.Equal(t, ZMQTopicRestdEvents, resultTopic)
		assert.Equal(t, testParam.Topic, resultMessage.Topic)
		assert.Equal(t, testParam.Message, resultMessage.Message)
	})

	// Tear down
	if shutdownErr := handler.Shutdown(); shutdownErr != nil {
		fmt.Printf("Failed to stop the goroutine running the ZMQ subscriber %v\n", shutdownErr.Error())
	}
	_ = subscriberSocket.Close()
}

func receiveSentMessage(subscriberSocket *zmq.Socket) (*ZmqMessage, string, error) {
	// fmt.Printf("Inside receiveSentMessage\n")
	msg, err := subscriberSocket.RecvMessageBytes(0)
	if err != nil {
		return nil, "", err
	}
	// fmt.Printf("message received: %v\n", msg)

	resultTopic := string(msg[0])
	resultMessage := &ZmqMessage{}

	err = json.Unmarshal(msg[1], &resultMessage.Message)
	if err != nil {
		return nil, "", err
	}

	resultMessage.Topic = resultTopic

	return resultMessage, resultTopic, nil
}

func createTestSubscriberSocket(socket string) (*zmq.Socket, error) {
	subSocket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}

	if err = subSocket.Connect(socket); err != nil {
		return nil, err
	}

	if err = subSocket.SetSubscribe(ZMQTopicRestdEvents); err != nil {
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

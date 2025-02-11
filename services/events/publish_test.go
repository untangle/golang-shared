package events

import (
	"fmt"
	"testing"

	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"github.com/untangle/golang-shared/testing/mocks"
	"google.golang.org/protobuf/proto"
)

const testSocketAddress = "inproc://testInProcDescriptor"

func TestEventPublisher(t *testing.T) {
	//Set up
	handler := newTestAlertHandler()
	startupErr := handler.Startup()
	assert.Nil(t, startupErr, "Failed to startup Zmq publisher in background")

	subscriberSocket, err := createTestSubscriberSocket(testSocketAddress)
	assert.Nil(t, err)

	t.Run("send message", func(t *testing.T) {
		testParam := Alerts.Alert{
			Type:     Alerts.AlertType_USER,
			Severity: Alerts.AlertSeverity_INFO,
			Message:  "test message",
		}

		// Send Event
		handler.Send(&testParam)
		resultMessage, resultTopic, err := receiveSentMessage(subscriberSocket)

		assert.Nil(t, err)
		assert.Equal(t, EventZMQTopic, resultTopic)
		assert.Equal(t, testParam.GetType(), resultMessage.GetType())
		assert.Equal(t, testParam.GetSeverity(), resultMessage.GetSeverity())
		assert.Equal(t, testParam.GetMessage(), resultMessage.GetMessage())
	})

	// Tear down
	if shutdownErr := handler.Shutdown(); shutdownErr != nil {
		fmt.Printf("Failed to stop the goroutine running the ZMQ subscriber %v\n", shutdownErr.Error())
	}
	_ = subscriberSocket.Close()
}

func receiveSentMessage(subscriberSocket *zmq.Socket) (*Alerts.Alert, string, error) {
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

func createTestSubscriberSocket(socket string) (*zmq.Socket, error) {
	pubSocket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}

	if err = pubSocket.Bind(socket); err != nil {
		return nil, err
	}

	if err = pubSocket.SetSubscribe(EventZMQTopic); err != nil {
		return nil, err
	}

	return pubSocket, nil
}

func newTestAlertHandler() *ZmqEventPublisher {
	return &ZmqEventPublisher{
		logger:                  mocks.NewMockLogger(),
		messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
		zmqPublisherShutdown:    make(chan bool),
		zmqPublisherStarted:     make(chan int32, 1),
		socketAddress:           testSocketAddress,
	}
}

package alerts

import (
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"testing"
)

const testSocketAddress = "inproc://testInProcDescriptor"

func TestAlertPublisher(t *testing.T) {
	//Set up
	handler := newTestAlertHandler()
	handler.startup()

	subscriberSocket, err := createTestSubscriberSocket(testSocketAddress)
	assert.Nil(t, err)

	t.Run("send message", func(t *testing.T) {
		testParam := Alerts.Alert{
			Type:     Alerts.AlertType_USER,
			Severity: Alerts.AlertSeverity_INFO,
			Message:  "test message",
		}

		// Send alert
		handler.Send(&testParam)
		resultMessage, resultTopic, err := receiveSentMessage(subscriberSocket)

		assert.Nil(t, err)
		assert.Equal(t, AlertZMQTopic, resultTopic)
		assert.Equal(t, testParam.GetType(), resultMessage.GetType())
		assert.Equal(t, testParam.GetSeverity(), resultMessage.GetSeverity())
		assert.Equal(t, testParam.GetMessage(), resultMessage.GetMessage())
	})

	// Tear down
	handler.Shutdown()
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

	if err = pubSocket.SetSubscribe(AlertZMQTopic); err != nil {
		return nil, err
	}

	return pubSocket, nil
}

func newTestAlertHandler() *AlertPublisher {
	return &AlertPublisher{
		logger:                  logger.NewLogger(),
		messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
		zmqPublisherShutdown:    make(chan bool),
		socketAddress:           testSocketAddress,
	}
}

package alerts

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"google.golang.org/protobuf/proto"
	"sync"
)

var alertPublisherSingleton *AlertPublisher
var once sync.Once

// AlertPublisher runs a ZMQ publisher socket in the background.
// When the Send method is called the alert is passed down to the
// ZMQ socket using a chanel and the message is published to ZMQ
// using the alert specific topic.
type AlertPublisher struct {
	logger                  *logger.Logger
	messagePublisherChannel chan ZmqMessage
	zmqPublisherShutdown    chan bool
	zmqPublisherStarted     chan bool
	socketAddress           string
}

// newAlertPublisher Gets the singleton instance of AlertPublisher.
func newAlertPublisher(logger *logger.Logger) *AlertPublisher {
	once.Do(func() {
		alertPublisherSingleton = &AlertPublisher{
			logger:                  logger,
			messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
			zmqPublisherShutdown:    make(chan bool),
			zmqPublisherStarted:     make(chan bool, 1),
			socketAddress:           PublisherSocketAddressConnect,
		}
	})

	return alertPublisherSingleton
}

// startup starts the ZMQ publisher in the background.
func (publisher *AlertPublisher) startup() {
	go publisher.zmqPublisher()

	// Blocks until the publisher starts.
	<-publisher.zmqPublisherStarted
}

// Shutdown stops the goroutine running the ZMQ subscriber and closes the channels used in the service.
func (publisher *AlertPublisher) Shutdown() {
	publisher.zmqPublisherShutdown <- true
	close(publisher.zmqPublisherShutdown)
	close(publisher.zmqPublisherStarted)
	close(publisher.messagePublisherChannel)
}

// Send publishes the alert to on the ZMQ publishing socket.
func (publisher *AlertPublisher) Send(alert *Alerts.Alert) {
	logger.Debug("Publish alert %v\n", alert)
	alertMessage, err := proto.Marshal(alert)
	if err != nil {
		logger.Err("Unable to marshal alert entry: %s\n", err)
		return
	}

	publisher.messagePublisherChannel <- ZmqMessage{Topic: AlertZMQTopic, Message: alertMessage}
}

// zmqPublisher initializes a ZMQ publishing socket and listens on the
// messagePublisherChannel. The received messages are published to the
// ZMQ publisher socket.
//
// This method should be run as a goroutine. The goroutine can be stopped
// by sending a message on the zmqPublisherShutdown channel.
func (publisher *AlertPublisher) zmqPublisher() {
	socket, err := publisher.setupZmqPubSocket()
	if err != nil {
		logger.Warn("Unable to setup ZMQ publisher socket: %s\n", err)
		return
	}
	defer socket.Close()

	publisher.zmqPublisherStarted <- true

	for {
		select {
		case msg := <-publisher.messagePublisherChannel:
			sentBytes, err := socket.SendMessage(msg.Topic, msg.Message)
			if err != nil {
				logger.Err("Publisher Send error: %s\n", err)
				continue
			}
			logger.Debug("Message sent: %v bytes\n", sentBytes)
		case <-publisher.zmqPublisherShutdown:
			logger.Info("ZMQ Publisher shutting down\n")
			return
		}
	}
}

// setupZmqPubSocket sets up the ZMQ socket for publishing alerts
func (publisher *AlertPublisher) setupZmqPubSocket() (soc *zmq.Socket, err error) {
	publisher.logger.Info("Setting up Alerts ZMQ publisher socket...\n")

	socket, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		publisher.logger.Err("Unable to open ZMQ publisher socket: %s\n", err)
		return nil, err
	}

	if err = socket.SetLinger(0); err != nil {
		publisher.logger.Err("Unable to SetLinger on ZMQ publisher socket: %s\n", err)
		return nil, err
	}

	if err = socket.Connect(publisher.socketAddress); err != nil {
		publisher.logger.Err("Unable to bind to ZMQ socket: %s\n", err)
		return nil, err
	}

	publisher.logger.Info("Alerts ZMQ Publisher started!\n")

	return socket, nil
}

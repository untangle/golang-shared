package alerts

import (
	"errors"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"google.golang.org/protobuf/proto"
)

var alertPublisherSingleton *ZmqAlertPublisher
var once sync.Once
var ErrPublisherStarted = errors.New("publisher already running")

type AlertPublisher interface {
	Send(alert *Alerts.Alert)
}

// ZmqAlertPublisher runs a ZMQ publisher socket in the background.
// When the Send method is called the alert is passed down to the
// ZMQ socket using a chanel and the message is published to ZMQ
// using the alert specific topic.
type ZmqAlertPublisher struct {
	logger                  logger.LoggerLevels
	messagePublisherChannel chan ZmqMessage
	zmqPublisherShutdown    chan bool
	zmqPublisherStarted     chan bool
	socketAddress           string
	started                 bool
}

func (publisher *ZmqAlertPublisher) Name() string {
	return "Alert publisher"
}

// NewZmqAlertPublisher Gets the singleton instance of ZmqAlertPublisher.
func NewZmqAlertPublisher(logger logger.LoggerLevels) *ZmqAlertPublisher {
	once.Do(func() {
		alertPublisherSingleton = &ZmqAlertPublisher{
			logger:                  logger,
			messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
			zmqPublisherShutdown:    make(chan bool),
			zmqPublisherStarted:     make(chan bool, 1),
			socketAddress:           PublisherSocketAddress,
		}
	})

	return alertPublisherSingleton
}

func NewDefaultAlertPublisher(logger logger.LoggerLevels) AlertPublisher {
	return NewZmqAlertPublisher(logger)
}

// Startup starts the ZMQ publisher in the background.
func (publisher *ZmqAlertPublisher) Startup() error {
	// Make sure it is not started twice.
	if publisher.started {
		return ErrPublisherStarted
	}

	go publisher.zmqPublisher()

	// Blocks until the publisher starts.
	publisher.started = <-publisher.zmqPublisherStarted

	return nil
}

// Shutdown stops the goroutine running the ZMQ subscriber and closes the channels used in the service.
func (publisher *ZmqAlertPublisher) Shutdown() error {
	publisher.zmqPublisherShutdown <- true
	close(publisher.zmqPublisherShutdown)
	close(publisher.zmqPublisherStarted)
	close(publisher.messagePublisherChannel)
	publisher.started = false

	return nil
}

// Send publishes the alert to on the ZMQ publishing socket.
func (publisher *ZmqAlertPublisher) Send(alert *Alerts.Alert) {
	// 2 reasons to set the timestamp here:
	// - the caller isn't responsible for setting the timestamp so we just need to set it in one place (here)
	// - we set it before putting it in queue, which means we have the timestamp of the alert creation, not the timestamp when it was processed
	alert.Timestamp = time.Now().Unix()

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
func (publisher *ZmqAlertPublisher) zmqPublisher() {
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
func (publisher *ZmqAlertPublisher) setupZmqPubSocket() (soc *zmq.Socket, err error) {
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

package events

import (
	"sync"
	"sync/atomic"
	"time"

	loggerModel "github.com/untangle/golang-shared/logger"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"

	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
)

const messageBuffer = 1000

var EventPublisherSingleton *ZmqEventPublisher
var once sync.Once

type EventPublisher interface {
	Send(alert *Alerts.Alert)
	SendZmqMessage(msg *ZmqMessage)
}

// ZmqEventPublisher runs a ZMQ publisher socket in the background.
// When the Send method is called the Events is passed down to the
// ZMQ socket using a chanel and the message is published to ZMQ
// using the Events specific topic.
type ZmqEventPublisher struct {
	logger                  loggerModel.LoggerLevels
	messagePublisherChannel chan ZmqMessage
	zmqPublisherShutdown    chan bool
	zmqPublisherStarted     chan int32
	socketAddress           string
	started                 int32
}

// Name returns the name of the event publisher.
func (publisher *ZmqEventPublisher) Name() string {
	return "Event publisher"
}

// NewZmqEventPublisher Gets the singleton instance of ZmqEventPublisher.
func NewZmqEventPublisher(logger loggerModel.LoggerLevels) *ZmqEventPublisher {
	once.Do(func() {
		EventPublisherSingleton = &ZmqEventPublisher{
			logger:                  logger,
			messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
			zmqPublisherShutdown:    make(chan bool),
			zmqPublisherStarted:     make(chan int32, 1),
			socketAddress:           PublisherSocketAddress,
		}
	})

	return EventPublisherSingleton
}

// NewDefaultEventPublisher returns a new instance of ZmqEventPublisher.
func NewDefaultEventPublisher(logger loggerModel.LoggerLevels) EventPublisher {
	return NewZmqEventPublisher(logger)
}

// Startup starts the ZMQ publisher in the background.
func (publisher *ZmqEventPublisher) Startup() error {
	publisher.logger.Info("Starting up the Events service\n")

	// Use CompareAndSwap to ensure only one thread can start the publisher.
	if !atomic.CompareAndSwapInt32(&publisher.started, 0, 1) {
		publisher.logger.Debug("Events service is already running.\n")
		return nil
	}

	go publisher.zmqPublisher()

	// Blocks until the publisher starts.
	<-publisher.zmqPublisherStarted

	return nil
}

// Shutdown stops the goroutine running the ZMQ subscriber and closes the channels used in the service.
func (publisher *ZmqEventPublisher) Shutdown() error {
	publisher.logger.Info("Shutting down the Events service\n")

	// Use CompareAndSwap to ensure only one thread can shut down the publisher.
	if !atomic.CompareAndSwapInt32(&publisher.started, 1, 0) {
		publisher.logger.Debug("Events service is not running.\n")
		return nil
	}

	publisher.zmqPublisherShutdown <- true
	close(publisher.zmqPublisherShutdown)
	close(publisher.zmqPublisherStarted)
	close(publisher.messagePublisherChannel)

	publisher.logger.Info("Events service has been shut down.\n")
	return nil
}

// Send publishes the alert to on the ZMQ publishing socket.
func (publisher *ZmqEventPublisher) Send(alert *Alerts.Alert) {
	// Make sure it is not shutdown.
	if atomic.LoadInt32(&publisher.started) == 0 {
		publisher.logger.Debug("Alerts service has been shutdown.\n")
		return
	}

	// 2 reasons to set the timestamp here:
	// - the caller isn't responsible for setting the timestamp so we just need to set it in one place (here)
	// - we set it before putting it in queue, which means we have the timestamp of the alert creation, not the timestamp when it was processed
	alert.Datetime = time.Now().Unix()

	publisher.logger.Debug("Publish alert %v\n", alert)
	alertMessage, err := proto.Marshal(alert)
	if err != nil {
		publisher.logger.Err("Unable to marshal alert entry: %s\n", err)
		return
	}

	publisher.messagePublisherChannel <- ZmqMessage{Topic: AlertZMQTopic, Message: alertMessage}
}

// SendZmqMessage publishes the event to on the ZMQ publishing socket.
func (publisher *ZmqEventPublisher) SendZmqMessage(event *ZmqMessage) {
	// Make sure it is not shutdown.
	if atomic.LoadInt32(&publisher.started) == 0 {
		publisher.logger.Debug("Events service has been shutdown.\n")
		return
	}
	publisher.logger.Debug("Publish event %v\n", event)
	// send event directly on messagePublisherChannel
	publisher.messagePublisherChannel <- *event
}

// zmqPublisher initializes a ZMQ publishing socket and listens on the
// messagePublisherChannel. The received messages are published to the
// ZMQ publisher socket.
//
// This method should be run as a goroutine. The goroutine can be stopped
// by sending a message on the zmqPublisherShutdown channel.
func (publisher *ZmqEventPublisher) zmqPublisher() {
	socket, err := publisher.setupZmqPubSocket()
	if err != nil {
		publisher.logger.Warn("Unable to setup ZMQ publisher socket: %s\n", err)
		return
	}
	defer socket.Close()

	publisher.zmqPublisherStarted <- 1

	for {
		select {
		case msg := <-publisher.messagePublisherChannel:
			sentBytes, err := socket.SendMessage(msg.Topic, msg.Message)
			if err != nil {
				publisher.logger.Err("Publisher Send error: %s\n", err)
				continue
			}
			publisher.logger.Debug("Message sent: %v bytes\n", sentBytes)
		case <-publisher.zmqPublisherShutdown:
			publisher.logger.Info("ZMQ Publisher shutting down\n")
			return
		}
	}
}

// setupZmqPubSocket sets up the ZMQ socket for publishing events
func (publisher *ZmqEventPublisher) setupZmqPubSocket() (soc *zmq.Socket, err error) {
	publisher.logger.Info("Setting up Events ZMQ publisher socket...\n")

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

	publisher.logger.Info("Events ZMQ Publisher started!\n")

	return socket, nil
}

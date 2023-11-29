package alerts

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	loggerModel "github.com/untangle/golang-shared/logger"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"google.golang.org/protobuf/proto"
)

const messageBuffer = 1000

var alertPublisherSingleton *ZmqAlertPublisher
var once sync.Once

type AlertPublisher interface {
	Send(alert *Alerts.Alert)
}

// ZmqAlertPublisher runs a ZMQ publisher socket in the background.
// When the Send method is called the alert is passed down to the
// ZMQ socket using a chanel and the message is published to ZMQ
// using the alert specific topic.
type ZmqAlertPublisher struct {
	logger                  loggerModel.LoggerLevels
	messagePublisherChannel chan ZmqMessage
	zmqPublisherShutdown    chan bool
	zmqPublisherStarted     chan int32
	socketAddress           string
	started                 int32
}

func (publisher *ZmqAlertPublisher) Name() string {
	return "Alert publisher"
}

// NewZmqAlertPublisher Gets the singleton instance of ZmqAlertPublisher.
func NewZmqAlertPublisher(logger loggerModel.LoggerLevels) *ZmqAlertPublisher {
	once.Do(func() {
		alertPublisherSingleton = &ZmqAlertPublisher{
			logger:                  logger,
			messagePublisherChannel: make(chan ZmqMessage, messageBuffer),
			zmqPublisherShutdown:    make(chan bool),
			zmqPublisherStarted:     make(chan int32, 1),
			socketAddress:           PublisherSocketAddress,
		}
	})

	return alertPublisherSingleton
}

func NewDefaultAlertPublisher(logger loggerModel.LoggerLevels) AlertPublisher {
	return NewZmqAlertPublisher(logger)
}

// Startup starts the ZMQ publisher in the background.
func (publisher *ZmqAlertPublisher) Startup() error {
	publisher.logger.Info("Starting up the Alerts service\n")

	// Make sure it is not started twice.
	if atomic.LoadInt32(&publisher.started) > 0 {
		publisher.logger.Debug("Alerts service is already running.\n")
		return nil
	}

	go publisher.zmqPublisher()

	// Blocks until the publisher starts.
	atomic.AddInt32(&publisher.started, <-publisher.zmqPublisherStarted)

	return nil
}

// Shutdown stops the goroutine running the ZMQ subscriber and closes the channels used in the service.
func (publisher *ZmqAlertPublisher) Shutdown() error {
	publisher.logger.Info("Shutting down the Alerts service\n")

	// Make sure it is not shutdown twice.
	if atomic.LoadInt32(&publisher.started) == 0 {
		publisher.logger.Debug("Alerts service is already shutdown.\n")
		return nil
	}

	publisher.zmqPublisherShutdown <- true
	close(publisher.zmqPublisherShutdown)
	close(publisher.zmqPublisherStarted)
	close(publisher.messagePublisherChannel)
	atomic.StoreInt32(&publisher.started, 0)

	return nil
}

// Send publishes the alert to on the ZMQ publishing socket.
func (publisher *ZmqAlertPublisher) Send(alert *Alerts.Alert) {
	// 2 reasons to set the timestamp here:
	// - the caller isn't responsible for setting the timestamp so we just need to set it in one place (here)
	// - we set it before putting it in queue, which means we have the timestamp of the alert creation, not the timestamp when it was processed
	fmt.Println("Inside Send Fn()1 ---------")
	alert.Timestamp = time.Now().Unix()
	fmt.Println("Inside Send Fn() 2 ALERT is %v:---------", alert)
	publisher.logger.Debug("Publish alert %v\n", alert)
	fmt.Println("Inside Send Fn() 3---------")
	alertMessage, err := proto.Marshal(alert)
	fmt.Println("Inside Send Fn() 4---------")
	if err != nil {
		fmt.Println("Inside Send Fn() 5---------")
		publisher.logger.Err("Unable to marshal alert entry: %s\n", err)
		fmt.Println("Inside Send Fn() 6---------")
		return
	}
	fmt.Println("Inside Send Fn() 7---------")
	publisher.messagePublisherChannel <- ZmqMessage{Topic: AlertZMQTopic, Message: alertMessage}
	fmt.Println("Inside Send Fn() 8---------")
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

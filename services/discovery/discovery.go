package discovery

import (
	"net"
	"net/http"
	"net/rpc"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
)

type zmqMessage struct {
	Topic   string
	Message []byte
}

// Messages to be published to the ZMQ socket
var messagePublisherChannel = make(chan *zmqMessage, 1000)

// List of registered collectors
var collectors []CollectorHandlerFunction = nil

// Channel to signal shutdown of periodic collector timer
var collectorTimerQuit = make(chan bool)

const (

	// CmdScanHost is a command to scan a host, argument is the hostnames
	CmdScanHost int = 1
	// CmdScanNet is a command to scan a network, argument is the networks (CIDR notation)
	CmdScanNet int = 2
)

// Command is commands that can be send back to the collector
type Command struct {
	Command   int
	Arguments []string
}

// CollectorHandlerFunction is the prototype for the registed call back handler.
// A callback handler should be able to handle an empty command set.
type CollectorHandlerFunction func([]Command)

// Startup the discovery service.
func Startup() {
	logger.Info("Starting discovery service\n")

	// Start the ZMQ publisher
	go zmqPublisher()

	// Start the collector timer
	collectorTimer := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-collectorTimer.C:
				callCollectors([]Command{})
			case <-collectorTimerQuit:
				collectorTimer.Stop()
				return
			}
		}
	}()

	rpcServer := new(DiscoveryRPCService)
	rpc.Register(rpcServer)
	rpc.HandleHTTP()

	lis, err := net.Listen("tcp", "127.0.0.1:5563")
	if err != nil {
		logger.Err("Failed to listen: %v\n", err)
		return
	}

	go http.Serve(lis, nil)

}

// Shutdown the discovery service.
func Shutdown() {
	logger.Info("Shutting down discovery service\n")
	// sgrpc.GracefulStop()
	collectorTimerQuit <- true
}

// zmqPublisher reads from the messageChannel and sends the events to the associated topic
func zmqPublisher() {
	socket, err := setupZmqPubSocket()
	if err != nil {
		logger.Warn("Unable to setup ZMQ Publishing socket: %s\n", err)
		return
	}
	defer socket.Close()

	for {
		select {
		case msg := <-messagePublisherChannel:
			sentBytes, err := socket.SendMessage(msg.Topic, msg.Message)
			if err != nil {
				logger.Err("Test publisher error: %s\n", err)
				break //  Interrupted
			}
			logger.Debug("Message sent, %v bytes sent.\n", sentBytes)

		}
	}
}

// setupZmqPubSocket sets up the ZMQ socket for publishing
func setupZmqPubSocket() (soc *zmq.Socket, err error) {
	logger.Info("Setting up ZMQ Publishing socket...\n")

	publisher, err := zmq.NewSocket(zmq.PUB)

	if err != nil {
		logger.Err("Unable to open ZMQ publisher socket: %s\n", err)
		return nil, err
	}

	err = publisher.SetLinger(0)
	if err != nil {
		logger.Err("Unable to SetLinger on ZMQ socket: %s\n", err)
		return nil, err
	}

	err = publisher.Bind("tcp://*:5562")

	if err != nil {
		logger.Err("Unable to bind to ZMQ socket: %s\n", err)
		return nil, err
	}

	logger.Info("ZMQ Publisher started!\n")

	return publisher, nil
}

func callCollectors(cmds []Command) {
	logger.Info("Calling collectors\n")
	for _, handler := range collectors {
		go handler(cmds)
	}
}

// RegisterCollector registers a collector callback function
func RegisterCollector(handler CollectorHandlerFunction) {
	collectors = append(collectors, handler)
}

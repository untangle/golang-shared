package discovery

import (
	"net"
	"net/http"
	"net/rpc"
	"sync"

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
//var collectors []CollectorHandlerFunction = nil
var collectors map[string]CollectorHandlerFunction
var collectorsLock sync.RWMutex

const (
	// CmdScanHost is a command to scan a host, argument is the hostnames
	CmdScanHost int = 1
	// CmdScanNet is a command to scan a network, argument is the networks (CIDR notation)
	CmdScanNet int = 2
)

func init() {
	collectors = make(map[string]CollectorHandlerFunction)
}

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
	for label, handler := range collectors {
		logger.Debug("Calling collector with label %s", label)
		go handler(cmds)
	}
}

// RegisterCollector registers a collector callback function.
// The collectorLabel is used for quick lookups of the collector function being registered
func RegisterCollector(collectorLabel string, handler CollectorHandlerFunction) {
	collectorsLock.Lock()
	defer collectorsLock.Unlock()
	collectors[collectorLabel] = handler
}

// Unregisters a collector function
func UnregisterCollector(collectorLabel string) {
	collectorsLock.Lock()
	defer collectorsLock.Unlock()
	delete(collectors, collectorLabel)
}

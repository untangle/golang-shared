package discovery

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	pluginName string = "discovery"

	// CmdScanHost is a command to scan a host, argument is the hostnames
	CmdScanHost int = 1
	// CmdScanNet is a command to scan a network, argument is the networks (CIDR notation)
	CmdScanNet int = 2
)

type zmqMessage struct {
	Topic   string
	Message []byte
}

var (
	discoverySingleton *Discovery
	once               sync.Once

	settingsPath []string = []string{"discovery"}
)

type discoveryPluginType struct {
	Enabled bool `json:"enabled"`
}

type Discovery struct {
	discoverySettings discoveryPluginType

	collectors              map[string]CollectorHandlerFunction
	collectorsLock          sync.RWMutex
	zmqPublisherShutdown    chan bool
	messagePublisherChannel chan *zmqMessage
}

func NewDiscovery() *Discovery {
	once.Do(func() {
		discoverySingleton = &Discovery{collectors: make(map[string]CollectorHandlerFunction),
			zmqPublisherShutdown: make(chan bool), messagePublisherChannel: make(chan *zmqMessage, 1000)}
	})

	return discoverySingleton
}

func (discovery *Discovery) InSync(settings interface{}) bool {
	newSettings, ok := settings.(discoveryPluginType)
	if !ok {
		logger.Warn("Discovery: Could not compare the settings file provided to the current plugin settings. The settings cannot be updated.")
		return false
	}

	if newSettings == discovery.discoverySettings {
		logger.Debug("Settings remain unchanged for the NMAP plugin\n")
		return true
	}

	logger.Info("Updating Discovery plugin settings\n")
	return false
}

func (discovery *Discovery) GetSettingsStruct() (interface{}, error) {
	var newSettings discoveryPluginType
	if err := settings.UnmarshalSettingsAtPath(&newSettings, settingsPath...); err != nil {
		return nil, fmt.Errorf("Discovery: %s", err.Error())
	}

	return newSettings, nil
}

func (discovery *Discovery) Name() string {
	return pluginName
}

func (discovery *Discovery) SyncSettings(settings interface{}) error {
	originalSettings := discovery.discoverySettings
	newSettings, ok := settings.(discoveryPluginType)
	if !ok {
		return fmt.Errorf("Discovery: Settings provided were %s but expected %s",
			reflect.TypeOf(settings).String(), reflect.TypeOf(discovery.discoverySettings).String())
	}

	discovery.discoverySettings = newSettings

	// If settings changed but the plugin was previously enabled, restart the plugin
	// for changes to take effect
	var shutdownError error
	if originalSettings.Enabled && discovery.discoverySettings.Enabled {
		shutdownError = discovery.Shutdown()
	}

	if discovery.discoverySettings.Enabled {
		discovery.startDiscovery()
	} else {
		shutdownError = discovery.Shutdown()
	}

	return shutdownError
}

// Command is commands that can be send back to the collector
type Command struct {
	Command   int
	Arguments []string
}

// CollectorHandlerFunction is the prototype for the registed call back handler.
// A callback handler should be able to handle an empty command set.
type CollectorHandlerFunction func([]Command)

// Startup the discovery plugin
func (discovery *Discovery) Startup() error {
	logger.Info("Starting Discovery Plugin\n")

	// Setup permanent resources
	rpcServer := new(DiscoveryRPCService)
	rpc.Register(rpcServer)
	rpc.HandleHTTP()

	lis, err := net.Listen("tcp", "127.0.0.1:5563")
	if err != nil {
		return fmt.Errorf("failed to listen %v", err)
	}

	go http.Serve(lis, nil)

	// Setup resources that can be turned on/off via settings
	// Grab the initial settings on startup
	settings, err := discovery.GetSettingsStruct()
	if err != nil {
		return err
	}

	// SyncSettings will start the plugin if it's enabled
	err = discovery.SyncSettings(settings)
	if err != nil {
		return err
	}

	return nil
}

// Startup func to be used on Restarts of the discovery plugin
func (discovery *Discovery) startDiscovery() {
	// Start the ZMQ publisher
	go discovery.zmqPublisher()
}

// Shutdown the discovery Plugin
func (discovery *Discovery) Shutdown() error {
	logger.Info("Shutting down discovery plugin\n")

	discovery.shutdownZmqPublisher()

	return nil
}

// zmqPublisher reads from the messageChannel and sends the events to the associated topic
func (discovery *Discovery) zmqPublisher() {
	socket, err := setupZmqPubSocket()
	if err != nil {
		logger.Warn("Unable to setup ZMQ Publishing socket: %s\n", err)
		return
	}
	defer socket.Close()

out:
	for {
		select {
		case msg := <-discovery.messagePublisherChannel:
			sentBytes, err := socket.SendMessage(msg.Topic, msg.Message)
			if err != nil {
				logger.Err("Test publisher error: %s\n", err)
				break //  Interrupted
			}
			logger.Debug("Message sent, %v bytes sent\n", sentBytes)
		case <-discovery.zmqPublisherShutdown:
			logger.Info("ZMQ Publisher shutting down\n")
			discovery.zmqPublisherShutdown <- true
			break out
		}
	}
}

func (discovery *Discovery) shutdownZmqPublisher() {
	// The send to kill the zmqPublisher goroutine must be non-blocking for
	// the case where the goroutine wasn't started in the first place.
	// The goroutine never starting occurs when the plugin is disabled
	select {
	case discovery.zmqPublisherShutdown <- true:
		// Send message
	default:
		// Do nothing if the message couldn't be sent
	}

	select {
	case <-discovery.zmqPublisherShutdown:
		logger.Info("Successful shutdown of the Discovery ZMQ Publisher\n")
	case <-time.After(1 * time.Second):
		logger.Warn("Failed to shutdown the Discovery ZMQ Publisher. It may never have been started\n")
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

func (discovery *Discovery) callCollectors(cmds []Command) {
	discovery.collectorsLock.Lock()
	defer discovery.collectorsLock.Unlock()
	logger.Info("Calling collectors\n")
	for label, handler := range discovery.collectors {
		logger.Debug("Calling collector with label %s", label)
		go handler(cmds)
	}
}

// RegisterCollector registers a collector callback function.
// The collectorLabel is used for quick lookups of the collector function being registered
// Will only run if the discovery plugin is enabled
func (discovery *Discovery) RegisterCollector(collectorLabel string, handler CollectorHandlerFunction) {
	if discovery.discoverySettings.Enabled {
		discovery.collectorsLock.Lock()
		defer discovery.collectorsLock.Unlock()
		discovery.collectors[collectorLabel] = handler
	}
}

// Unregisters a collector function
// Will only run if the discovery plugin is enabled
func (discovery *Discovery) UnregisterCollector(collectorLabel string) {
	if discovery.discoverySettings.Enabled {
		discovery.collectorsLock.Lock()
		defer discovery.collectorsLock.Unlock()
		delete(discovery.collectors, collectorLabel)
	}
}

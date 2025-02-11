// Package events provides a service for alert publishing on a ZMQ socket.
// Usage:
// 		events.Startup(logger) // this is just so it initializes the publisher on service Startup, not on the first call
// 		events.Publisher().Send(alert1)
// 		events.Publisher().Send(alert2)
//		...
//		events.Shutdown()

package events

import "github.com/untangle/golang-shared/logger"

const (
	// AlertZMQTopic Topic name to be used when sending alerts.
	ZMQTopicRestdEvents     = "arista:restd:alert"
	ZMQTopicPakcetdEvents   = "arista:packetd:alert"
	sessionsZMQTopic        = "untangle:packetd:sessions"
	interfaceStatsZMQTopic  = "untangle:packetd:interface-stats"
	sessionStatsZMQTopic    = "untangle:packetd:session-stats"
	activeSessionsZMQTopic  = "untangle:packetd:active-sessions"
	dosEventsZMQTopic       = "untangle:packetd:dos-events"
	PublisherSocketAddress  = "ipc:///var/zmq_event_publisher"
	SubscriberSocketAddress = "ipc:///var/zmq_event_subscriber"
)

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

var publisher EventPublisher

// Publisher returns the Publisher singleton.
func Publisher(logger logger.LoggerLevels) EventPublisher {
	if publisher == nil {
		zmqPublisher := NewZmqEventPublisher(logger)
		_ = zmqPublisher.Startup()

		publisher = zmqPublisher
	}

	return publisher
}

// Startup is called when the service that uses Events starts
func Startup(logger logger.LoggerLevels) {
	logger.Info("Starting up the Events service\n")
	Publisher(logger)
}

// Shutdown is called when the service that uses Events stops
func Shutdown() {
	if publisher == nil {
		return
	}

	var zmqPublisher interface{} = Publisher(nil)
	_ = zmqPublisher.(*ZmqEventPublisher).Shutdown()
	publisher = nil
}

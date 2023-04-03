// Package alerts provides a service for alert publishing on a ZMQ socket.
// Usage:
// 		alerts.Startup() // this is just so it initializes the publisher on service Startup, not on the first call
// 		alerts.Publisher().Send(alert1)
// 		alerts.Publisher().Send(alert2)
//		...
//		alerts.Shutdown()

package alerts

// AlertZMQTopic Topic name to be used when sending alerts.
const AlertZMQTopic string = "arista:alertd:alert"

const PublisherSocketAddress = "ipc:///var/zmq_alert_publisher"
const SubscriberSocketAddress = "ipc:///var/zmq_alert_subscriber"

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

var publisher AlertPublisher

// Publisher returns the Publisher singleton.
func Publisher(logger AlertsLogger) AlertPublisher {
	if publisher == nil {
		zmqPublisher := NewZmqAlertPublisher(logger)
		_ = zmqPublisher.Startup()

		publisher = zmqPublisher
	}

	return publisher
}

// Startup is called when the service that uses alerts starts
func Startup(logger AlertsLogger) {
	logger.Info("Starting up the Alerts service\n")
	Publisher(logger)
}

// Shutdown is called when the service that uses alerts stops
func Shutdown() {
	if publisher == nil {
		return
	}

	var zmqPublisher interface{} = Publisher(nil)
	_ = zmqPublisher.(*ZmqAlertPublisher).Shutdown()
	publisher = nil
}

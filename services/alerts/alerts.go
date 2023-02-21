// Package alerts provides a service for alert publishing on a ZMQ socket.
// Usage:
// 		alerts.Startup() // this is just so it initializes the publisher on service Startup, not on the first call
// 		alerts.Publisher().Send(alert1)
// 		alerts.Publisher().Send(alert2)
//		...
//		alerts.Shutdown()

package alerts

import logService "github.com/untangle/golang-shared/services/logger"

// AlertZMQTopic Topic name to be used when sending alerts.
const AlertZMQTopic string = "arista:alertd:alert"

const PublisherSocketAddress = "ipc:///var/zmq_alert_publisher"
const SubscriberSocketAddress = "ipc:///var/zmq_alert_subscriber"

const messageBuffer = 1000

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

var loggerInstance = logService.GetLoggerInstance()
var publisher AlertPublisher

// Publisher returns the Publisher singleton.
func Publisher() AlertPublisher {
	if publisher == nil {
		zmqPublisher := NewZmqAlertPublisher(loggerInstance)
		_ = zmqPublisher.Startup()

		publisher = zmqPublisher
	}

	return publisher
}

// Startup is called when the service that uses alerts starts
func Startup() {
	loggerInstance.Info("Starting up the Alerts service\n")
	Publisher()
}

// Shutdown is called when the service that uses alerts stops
func Shutdown() {
	loggerInstance.Info("Shutting down the Alerts service\n")
	if publisher == nil {
		return
	}

	var zmqPublisher interface{} = Publisher()
	_ = zmqPublisher.(*ZmqAlertPublisher).Shutdown()
	publisher = nil
}

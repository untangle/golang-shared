// Package alerts provides a service for alert publishing on a ZMQ socket.
// Usage:
// 		alerts.Publisher().Send(alert1)
// 		alerts.Publisher().Send(alert2)
//		...
//		alerts.Publisher().Shutdown()

package alerts

import logService "github.com/untangle/golang-shared/services/logger"

// AlertZMQTopic Topic name to be used when sending alerts.
const AlertZMQTopic string = "arista:alertd:alert"

const PublisherSocketAddressConnect = "tcp://127.0.0.1:5565"
const PublisherSocketAddressBind = "tcp://*:5565"

const SubscriberSocketAddressConnect = "tcp://127.0.0.1:5564"
const SubscriberSocketAddressBind = "tcp://*:5564"

const messageBuffer = 1000

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

var loggerInstance = logService.GetLoggerInstance()
var publisher *AlertPublisher

// Publisher returns the Publisher singleton.
func Publisher() *AlertPublisher {
	if publisher == nil {
		publisher = newAlertPublisher(loggerInstance)
		publisher.startup()
	}

	return publisher
}

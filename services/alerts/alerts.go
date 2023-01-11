// Package alerts provides a service for alert publishing on a ZMQ socket.
// Usage:
// 		alerts.Startup()
// 		alerts.Publisher.Send(alert)
//		...
//		alerts.Shutdown()

package alerts

import logService "github.com/untangle/golang-shared/services/logger"

// AlertZMQTopic Topic name to be used when sending alerts.
const AlertZMQTopic string = "arista:reportd:alertd:alert"

const socketAddress = "tcp://*:5562"
const messageBuffer = 1000

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

var loggerInstance = logService.GetLoggerInstance()
var Publisher *AlertPublisher

// Startup initializes the Publisher
func Startup() {
	if Publisher != nil {
		return
	}

	Publisher = NewAlertPublisher(loggerInstance)
	Publisher.startup()
}

// Shutdown stops and resets the publisher.
func Shutdown() {
	if Publisher == nil {
		return
	}

	Publisher.shutdown()
	Publisher = nil
}

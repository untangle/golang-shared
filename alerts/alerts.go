package alerts

// Topic name to be used when sending alerts.
const AlertZMQTopic string = "arista:reportd:alertd:alert"

// ZmqMessage is a message sent over a zmq bus for us to consume.
type ZmqMessage struct {
	Topic   string
	Message []byte
}

package events

import (
	"sync"

	protoAlerts "github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"google.golang.org/protobuf/proto"
)

// MockEventPublisher is a test stub that keeps all the Alerts,
// Events, last alert, last Event that were sent to it.
type MockEventPublisher struct {
	sync.Mutex
	Alerts    []*protoAlerts.Alert
	Events    []*ZmqMessage
	LastAlert *protoAlerts.Alert
	LastEvent *ZmqMessage
}

func (m *MockEventPublisher) Send(alert *protoAlerts.Alert) {
	m.Lock()
	defer m.Unlock()
	m.LastAlert = alert
	m.Alerts = append(m.Alerts, alert)
}

func (m *MockEventPublisher) SendZmqMessage(event *ZmqMessage) {
	m.Lock()
	defer m.Unlock()
	m.LastEvent = event
	m.Events = append(m.Events, event)
}

// GetLastAlert gets the last alert.
func (m *MockEventPublisher) GetLastAlert() *protoAlerts.Alert {
	m.Lock()
	defer m.Unlock()
	copy := proto.Clone(m.LastAlert)
	return copy.(*protoAlerts.Alert)
}

// GetLastEvent gets the last event.
func (m *MockEventPublisher) GetLastEvent() *ZmqMessage {
	m.Lock()
	defer m.Unlock()
	copy := *m.LastEvent
	return &copy
}

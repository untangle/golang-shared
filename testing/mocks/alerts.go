package mocks

import (
	"sync"

	protoAlerts "github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	"google.golang.org/protobuf/proto"
)

// MockAlertPublisher is not a real mock, it is a test stub that keeps
// the last alert. TODO: make this a real mock.
type MockAlertPublisher struct {
	sync.Mutex
	LastAlert *protoAlerts.Alert
}

func (m *MockAlertPublisher) Send(alert *protoAlerts.Alert) {
	m.Lock()
	defer m.Unlock()
	m.LastAlert = alert
}

// GetLastAlert gets the last alert.
func (m *MockAlertPublisher) GetLastAlert() *protoAlerts.Alert {
	m.Lock()
	defer m.Unlock()
	copy := proto.Clone(m.LastAlert)
	return copy.(*protoAlerts.Alert)
}

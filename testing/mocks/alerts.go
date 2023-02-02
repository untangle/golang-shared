package mocks

import (
	protoAlerts "github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
)

type MockAlertPublisher struct {
	LastAlert *protoAlerts.Alert
}

func (m *MockAlertPublisher) Send(alert *protoAlerts.Alert) {
	m.LastAlert = alert
}

package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	protoAlerts "github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	protoDiscoverd "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

func TestNewDeviceDiscoveredAlert(t *testing.T) {
	devicesList := NewDevicesList()
	emptyCb := func() {}

	alerts := []*protoAlerts.Alert{}
	sendAlert := func(a *protoAlerts.Alert) {
		alerts = append(alerts, a)
	}

	devicesToAdd := []*DeviceEntry{
		{
			// empty device entry that shouldn't generate any alert
		},
		{
			// all zeros mac shouldn't generate any alert
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				MacAddress: "00:00:00:00:00:00",
			},
		},
		{
			// new mac address generates alert
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				MacAddress: "11:11:11:11:11:11",
			},
		},
		{
			// same mac as the previous one, even with new IPs, shouldn't generate
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				MacAddress: "11:11:11:11:11:11",
				Neigh: map[string]*protoDiscoverd.NEIGH{
					"192.168.56.1": {},
				},
			},
		},
		{
			// new mac with some IPs
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				MacAddress: "22:22:22:22:22:22",
				Neigh: map[string]*protoDiscoverd.NEIGH{
					"192.168.56.2": {},
				},
				Lldp: map[string]*protoDiscoverd.LLDP{
					"192.168.56.3": {},
				},
			},
		},
		{
			// unknown mac but existing IP shouldn't generate alert
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				Neigh: map[string]*protoDiscoverd.NEIGH{
					"192.168.56.2": {},
				},
			},
		},
		{
			// unknown mac and new IP generates new alert
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				Neigh: map[string]*protoDiscoverd.NEIGH{
					"192.168.56.56": {},
				},
			},
		},
		// TODO: this seems to fall into the `if !found` case
		//  but should it ???
		// TODO: A S K   A B O U T   T H I S
		// ./types.go:309  || oldEntry.MacAddress == ""
		{
			// new mac and existing IP shouldn't generate alert
			DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
				MacAddress: "33:33:33:33:33:33",
				Neigh: map[string]*protoDiscoverd.NEIGH{
					"192.168.56.56": {},
				},
			},
		},
	}
	expectedAlerts := []*protoAlerts.Alert{
		{
			Params: map[string]string{
				"ips":        "",
				"macAddress": "11:11:11:11:11:11",
			},
		},
		{
			Params: map[string]string{
				"ips":        "192.168.56.2,192.168.56.3",
				"macAddress": "22:22:22:22:22:22",
			},
		},
		{
			Params: map[string]string{
				"ips":        "192.168.56.56",
				"macAddress": "",
			},
		},
	}

	for _, d := range devicesToAdd {
		devicesList.mergeOrAddEntry(sendAlert, d, emptyCb, false)
	}

	assert.Equal(t, len(expectedAlerts), len(alerts))
	for i, a := range alerts {
		e := expectedAlerts[i]

		assert.Equal(t, protoAlerts.AlertType_DISCOVERY, a.Type)
		assert.Equal(t, "ALERT_NEW_DEVICE_DISCOVERED", a.Message)
		assert.Equal(t, e.Params["ips"], a.Params["ips"])
		assert.Equal(t, e.Params["macAddress"], a.Params["macAddress"])
	}

	// check that no alerts are created when reading from DB at startup
	warmUpDeviceList := NewDevicesList()
	alerts = []*protoAlerts.Alert{}
	for _, d := range devicesToAdd {
		warmUpDeviceList.mergeOrAddEntry(sendAlert, d, emptyCb, true)
	}
	assert.Equal(t, 0, len(alerts))
}

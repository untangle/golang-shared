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

	var alert *protoAlerts.Alert
	sendAlert := func(a *protoAlerts.Alert) {
		alert = a
	}

	type testParams struct {
		ips        string
		macAddress string
	}
	type testCase struct {
		name   string
		entry  DeviceEntry
		params *testParams // nil if a new alert should not be created
	}

	testCases := []testCase{
		{
			name: "empty device shouldn't generate alert",
			// empty device entry that shouldn't generate any alert
			params: nil,
		},
		{
			name: "all zeros mac shouldn't generate alert",
			entry: DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					MacAddress: "00:00:00:00:00:00",
				},
			},
			params: nil,
		},
		{
			name: "new mac address generates alert",
			entry: DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					MacAddress: "11:11:11:11:11:11",
				},
			},
			params: &testParams{
				ips:        "",
				macAddress: "11:11:11:11:11:11",
			},
		},
		{
			name: "same mac as the previous one, even with new IPs, shouldn't generate alert",
			entry: DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					MacAddress: "11:11:11:11:11:11",
					Neigh: map[string]*protoDiscoverd.NEIGH{
						"192.168.56.1": {},
					},
				},
			},
			params: nil,
		},
		{
			name: "new mac with some IPs generates alert",
			entry: DeviceEntry{
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
			params: &testParams{
				ips:        "192.168.56.2,192.168.56.3",
				macAddress: "22:22:22:22:22:22",
			},
		},
		{
			name: "unknown mac but existing IP shouldn't generate alert",
			entry: DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					Neigh: map[string]*protoDiscoverd.NEIGH{
						"192.168.56.2": {},
					},
				},
			},
			params: nil,
		},
		{
			name: "unknown mac and new IP generates new alert",
			entry: DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					Neigh: map[string]*protoDiscoverd.NEIGH{
						"192.168.56.56": {},
					},
				},
			},
			params: &testParams{
				ips:        "192.168.56.56",
				macAddress: "",
			},
		},
		// TODO: this seems to fall into the `if !found` case
		//  but should it ???
		// TODO: A S K   A B O U T   T H I S
		// ./types.go:309  || oldEntry.MacAddress == ""
		{
			name: "new mac and existing IP shouldn't generate alert",
			entry: DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					MacAddress: "33:33:33:33:33:33",
					Neigh: map[string]*protoDiscoverd.NEIGH{
						"192.168.56.56": {},
					},
				},
			},
			params: nil,
		},
	}
	for _, test := range testCases {
		// reset alert so we can check if a new one was created during this step
		alert = nil
		devicesList.mergeOrAddWithAlert(&test.entry, emptyCb, sendAlert)

		if test.params == nil {
			assert.Nil(t, alert)
			return
		}

		assert.NotNil(t, alert)

		assert.Equal(t, protoAlerts.AlertType_DISCOVERY, alert.Type)
		assert.Equal(t, "ALERT_NEW_DEVICE_DISCOVERED", alert.Message)
		assert.Equal(t, test.params.ips, alert.Params["ips"])
		assert.Equal(t, test.params.macAddress, alert.Params["macAddress"])
	}
}

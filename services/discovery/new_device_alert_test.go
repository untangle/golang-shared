package discovery

import (
	"github.com/untangle/golang-shared/testing/mocks"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	protoAlerts "github.com/untangle/golang-shared/structs/protocolbuffers/Alerts"
	protoDiscoverd "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

func TestNewDeviceAlert(t *testing.T) {
	devicesList := NewDevicesList()
	emptyCb := func() {}

	alertPublisher := &mocks.MockAlertPublisher{}

	type testParams struct {
		ips        string
		macAddress string
	}
	type testCase struct {
		name   string
		entry  *DeviceEntry
		params *testParams // nil if a new alert should not be created
	}

	testCases := []testCase{
		{
			name:  "empty device shouldn't generate alert",
			entry: &DeviceEntry{},
			// empty device entry that shouldn't generate any alert
			params: nil,
		},
		{
			name: "all zeros mac shouldn't generate alert",
			entry: &DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					MacAddress: "00:00:00:00:00:00",
				},
			},
			params: nil,
		},
		{
			name: "new mac address generates alert",
			entry: &DeviceEntry{
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
			entry: &DeviceEntry{
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
			entry: &DeviceEntry{
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
			entry: &DeviceEntry{
				DiscoveryEntry: protoDiscoverd.DiscoveryEntry{
					Neigh: map[string]*protoDiscoverd.NEIGH{
						"192.168.56.2": {},
					},
				},
			},
			params: nil,
		},
		{
			name: "unknown mac and new IP does generates new alert",
			entry: &DeviceEntry{
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
		{
			name: "new mac and existing IP shouldn't generate alert",
			entry: &DeviceEntry{
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
		alertPublisher.LastAlert = nil
		devicesList.mergeOrAddWithAlert(test.entry, emptyCb, alertPublisher)

		if test.params == nil {
			assert.Nil(t, alertPublisher.LastAlert, test.name)
			continue
		}

		assert.NotNil(t, alertPublisher.LastAlert, test.name)

		assert.Equal(t, protoAlerts.AlertType_DISCOVERY, alertPublisher.LastAlert.Type, test.name)
		assert.Equal(t, "ALERT_NEW_DEVICE_DISCOVERED", alertPublisher.LastAlert.Message, test.name)
		assert.Equal(t, len(test.params.ips), len(alertPublisher.LastAlert.Params["ips"]), test.name)

		// in alert.Params["ips"] we do not know the order of the ips
		// we check that it contains all the expected ips
		for _, ip := range strings.Split(test.params.ips, ",") {
			assert.Contains(t, alertPublisher.LastAlert.Params["ips"], ip, test.name)
		}

		assert.Equal(t, test.params.macAddress, alertPublisher.LastAlert.Params["macAddress"], test.name)
	}
}

// Mocks

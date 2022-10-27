package arp

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	disc "github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/settings"
	disc_pb "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	"github.com/vishvananda/netlink"
	"google.golang.org/protobuf/proto"
)

func TestIpNeighbourEntries(t *testing.T) {
	for _, testcase := range ipNeighDevicesTestcases {
		t.Run(testcase.name, func(t *testing.T) {
			handler := newMockNetlinkHandler(testcase.devices)
			scanner := &netlinkScanner{
				settings:      settings.NewSettingsFile(testcase.settingsFile),
				handler:       handler,
				timestamp:     timestampCallback,
				skippedStates: testcase.skippStates,
			}

			entries, err := scanner.getIpNeighbourEntries()
			assert.Nil(t, err)

			// Check if the functions where called with the right parameters.
			// LinkByName should be called with the name of the interfaces that are not WAM.
			// NeighList should be called with the indexes of the non WAM interfaces (see mapping in mock).
			assert.Equal(t, testcase.linkByNameParams, handler.linkByNameParams)
			assert.Equal(t, testcase.neighListParams, handler.neighListParams)
			assert.Equal(t, len(testcase.entries), len(entries))

			for i := range testcase.entries {
				assert.True(t, proto.Equal(entries[i], testcase.entries[i]))
			}
		})
	}
}

type ipNeighProviderParam struct {
	name             string
	settingsFile     string
	devices          []netlink.Neigh
	entries          []*disc.DeviceEntry
	linkByNameParams string
	neighListParams  string
	skippStates      map[int]bool
}

var ipNeighDevicesTestcases = []ipNeighProviderParam{
	{
		name:             "Test when there are no devices",
		settingsFile:     "./testdata/settings.json",
		devices:          []netlink.Neigh{},
		entries:          []*disc.DeviceEntry{},
		linkByNameParams: ",eth1",
		neighListParams:  ",2-0",
		skippStates:      make(map[int]bool),
	},
	{
		name:             "Test mulltiple non WAN interfaces",
		settingsFile:     "./testdata/settings-2-non-wans.json",
		linkByNameParams: ",eth1,eth2",
		neighListParams:  ",2-0,3-0", skippStates: make(map[int]bool),
	},
	{
		name:         "Test serialization",
		settingsFile: "./testdata/settings.json",
		devices:      devices,
		entries: []*disc.DeviceEntry{
			{DiscoveryEntry: disc_pb.DiscoveryEntry{
				MacAddress:  "00:11:22:33:44:55",
				IPv4Address: "192.168.68.1",
				Arp: &disc_pb.ARP{
					Ip:  "192.168.68.1",
					Mac: "00:11:22:33:44:55", LastUpdate: timestampCallback(),
					State: "REACHABLE",
				}}},
			{DiscoveryEntry: disc_pb.DiscoveryEntry{
				MacAddress: "00:11:22:33:44:55",
				// We do not have an IPV4 ip.
				IPv4Address: "",
				Arp: &disc_pb.ARP{
					Ip:         "f0e:d0c:b0a:908:706:504:302:101",
					Mac:        "00:11:22:33:44:55",
					LastUpdate: timestampCallback(),
					State:      "STALE",
				}}},
			{DiscoveryEntry: disc_pb.DiscoveryEntry{
				MacAddress: "11:11:11:33:44:55",
				// We do not have an IPV4 ip.
				IPv4Address: "192.168.56.3",
				Arp: &disc_pb.ARP{
					Ip:         "192.168.56.3",
					Mac:        "11:11:11:33:44:55",
					LastUpdate: timestampCallback(),
					State:      "NOARP",
				}}},
		},
		linkByNameParams: ",eth1",
		neighListParams:  ",2-0",
		skippStates:      make(map[int]bool),
	},
	{
		name:         "Test serialization with skipped statuses",
		settingsFile: "./testdata/settings.json",
		devices:      devices,
		entries: []*disc.DeviceEntry{
			{DiscoveryEntry: disc_pb.DiscoveryEntry{
				MacAddress:  "00:11:22:33:44:55",
				IPv4Address: "192.168.68.1",
				Arp: &disc_pb.ARP{
					Ip:  "192.168.68.1",
					Mac: "00:11:22:33:44:55", LastUpdate: timestampCallback(),
					State: "REACHABLE",
				}}},
		},
		linkByNameParams: ",eth1",
		neighListParams:  ",2-0",
		skippStates:      map[int]bool{netlink.NUD_STALE: true, netlink.NUD_NOARP: true},
	},
}

// Mocks.

// The netlink handler can be mocked without abstractization as well (see netlink tests), but test require root priviledges to run.
// mockNetlinkHandler helps substituting the netlink componenet in order to get consistent results for tests.
type mockNetlinkHandler struct {
	devices        []netlink.Neigh
	interfaceNames map[string]int

	linkByNameParams string
	neighListParams  string
}

func newMockNetlinkHandler(devices []netlink.Neigh) *mockNetlinkHandler {
	interfaceName := make(map[string]int)
	interfaceName["eth0"] = 1
	interfaceName["eth1"] = 2
	interfaceName["eth2"] = 3

	return &mockNetlinkHandler{
		interfaceNames: interfaceName,
		devices:        devices,
	}
}

// LinkByName returns mock details for a network device name
// It stores the arguments of the function in order to check if it was called in the right context.
func (m *mockNetlinkHandler) LinkByName(name string) (netlink.Link, error) {
	m.linkByNameParams = fmt.Sprintf("%s,%s", m.linkByNameParams, name)

	return &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: m.interfaceNames[name], Name: name}}, nil
}

// NeighList returns mock neighbours (hosts) for a device.
// It stores the arguments of the function in order to check if it was called in the right context.
func (m *mockNetlinkHandler) NeighList(linkIndex, family int) ([]netlink.Neigh, error) {
	m.neighListParams = fmt.Sprintf("%s,%d-%d", m.neighListParams, linkIndex, family)

	return m.devices, nil
}

// Delete releases the used netlink sockets
func (m *mockNetlinkHandler) Delete() {
}

var timestampCallback = func() int64 {
	return 123456
}

// Test data.

var hardware1, _ = net.ParseMAC("00:11:22:33:44:55")
var hardware2, _ = net.ParseMAC("11:11:11:33:44:55")

var devices = []netlink.Neigh{
	{
		LinkIndex:    2,
		Family:       2,
		HardwareAddr: hardware1,
		IP:           net.IPv4(192, 168, 68, 1),
		State:        netlink.NUD_REACHABLE,
	},
	{
		LinkIndex:    2,
		Family:       10,
		HardwareAddr: hardware1,
		IP:           net.IP{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 1},
		State:        netlink.NUD_STALE,
	},
	{
		LinkIndex:    2,
		Family:       2,
		HardwareAddr: hardware2,
		IP:           net.IPv4(192, 168, 56, 3),
		State:        netlink.NUD_NOARP,
	},
}

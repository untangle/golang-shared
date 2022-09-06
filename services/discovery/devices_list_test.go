package discovery

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
	mfw_ifaces "github.com/untangle/golang-shared/util/net"
	"google.golang.org/protobuf/proto"
)

type DeviceListTestSuite struct {
	suite.Suite

	mac1         string
	mac2         string
	mac3         string
	devicesTable *DevicesList
	now          time.Time
	oneHourAgo   time.Time
	halfHourAgo  time.Time
}

func (suite *DeviceListTestSuite) SetupTest() {
	suite.now = time.Now()
	suite.oneHourAgo = suite.now.Add(-1 * time.Hour)
	suite.halfHourAgo = suite.now.Add(-30 * time.Minute)
	suite.mac1 = "00:11:22:33:44:55"
	suite.mac2 = "00:11:22:33:44:66"
	suite.mac3 = "00:11:33:44:55:66"
	suite.devicesTable = &DevicesList{
		Devices: map[string]*DeviceEntry{
			suite.mac1: {disco.DiscoveryEntry{
				MacAddress:  suite.mac1,
				LastUpdate:  suite.oneHourAgo.Unix(),
				IPv4Address: "192.168.56.1",
			}},
			suite.mac2: {disco.DiscoveryEntry{
				MacAddress:  suite.mac2,
				LastUpdate:  suite.halfHourAgo.Unix(),
				IPv4Address: "192.168.56.2",
			}},
			suite.mac3: {disco.DiscoveryEntry{
				MacAddress:  suite.mac3,
				LastUpdate:  suite.halfHourAgo.Unix(),
				IPv4Address: "192.168.56.3",
			}},
		},
	}
}

func (suite *DeviceListTestSuite) TestListing() {

	getMacs := func(devs []*DeviceEntry) (output []string) {
		output = []string{}
		for _, dev := range devs {
			output = append(output, dev.MacAddress)
		}
		return
	}
	mac1hwaddr, _ := net.ParseMAC(suite.mac1)
	mac3hwaddr, _ := net.ParseMAC(suite.mac3)
	localInterfaces := []net.Interface{
		{Index: 1, MTU: 1500, Name: "bar", HardwareAddr: mac1hwaddr},
		{Index: 2, MTU: 1500, Name: "bar", HardwareAddr: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x88}},
		{Index: 3, MTU: 1500, Name: "bar", HardwareAddr: mac3hwaddr},
	}

	type testSpec struct {
		predicates       []ListPredicate
		expectedListMacs []string
		description      string
	}
	tests := []testSpec{
		{
			predicates:       []ListPredicate{},
			expectedListMacs: []string{suite.mac1, suite.mac2, suite.mac3},
			description:      "Empty predicate list should return all MACs.",
		},
		{
			predicates:       []ListPredicate{WithUpdatesWithinDuration(time.Hour - time.Second)},
			expectedListMacs: []string{suite.mac2, suite.mac3},
			description:      "List of macs with an update time within the hour should be mac2, mac3.",
		},
		{
			predicates: []ListPredicate{
				IsNotFromLocalInterface(localInterfaces),
			},
			expectedListMacs: []string{suite.mac2},
			description:      "MAcs not on the local interface list should be one, mac2.",
		},
		{
			predicates: []ListPredicate{
				IsNotFromLocalInterface(localInterfaces),
				WithUpdatesWithinDuration(time.Minute * 20),
			},
			expectedListMacs: []string{},
			description:      "There are no devices in the list that have an update time within 20 minutes and are not local.",
		},
	}

	for _, test := range tests {
		theList := suite.devicesTable.listDevices(test.predicates...)
		assert.ElementsMatch(suite.T(), test.expectedListMacs, getMacs(theList), test.description)
	}
}

func (suite *DeviceListTestSuite) TestMarshallingList() {
	mac1hwaddr, _ := net.ParseMAC(suite.mac1)
	mac3hwaddr, _ := net.ParseMAC(suite.mac3)
	localInterfaces := []net.Interface{
		{Index: 1, MTU: 1500, Name: "bar", HardwareAddr: mac1hwaddr},
		{Index: 2, MTU: 1500, Name: "bar", HardwareAddr: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x88}},
		{Index: 3, MTU: 1500, Name: "bar", HardwareAddr: mac3hwaddr},
	}
	output, err := suite.devicesTable.ApplyToDeviceList(
		func(list []*DeviceEntry) (interface{}, error) {
			return json.Marshal(list)
		},
		IsNotFromLocalInterface(localInterfaces))
	suite.Nil(err)
	bytes := output.([]byte)
	fmt.Printf("JSON string output: %s\n", string(bytes))
}

func (suite *DeviceListTestSuite) TestGetDevFromIP() {
	dev := suite.devicesTable.Devices[suite.mac1]
	foundDev := suite.devicesTable.GetDeviceEntryFromIP(dev.IPv4Address)
	suite.True(proto.Equal(dev, foundDev))
}

type MockNetInterface struct {
	net.Interface
	mock.Mock
}

func (mock *MockNetInterface) GetName() string {
	return mock.Name
}

func (mock *MockNetInterface) Addrs() ([]net.Addr, error) {
	outputs := mock.Called()
	return outputs.Get(0).([]net.Addr), outputs.Error(1)
}

func (suite *DeviceListTestSuite) TestWANDeviceFilter() {
	mfwInterfaces := []*mfw_ifaces.Interface{
		{
			Name:   "wan0",
			Device: "eth0",
			IsWAN:  true,
		},
		{
			Name:  "lan0",
			IsWAN: false,
		},
	}
	localInterfaces := []*MockNetInterface{
		{
			Interface: net.Interface{Name: "eth0"},
		},
		{
			Interface: net.Interface{Name: "eth1"},
		},
	}
	localInterfacesSysNet := []SystemNetInterface{}
	for _, iface := range localInterfaces {
		localInterfacesSysNet = append(localInterfacesSysNet, iface)
	}
	_, wanNet, err := net.ParseCIDR("192.168.56.2/24")
	suite.Nil(err)
	// eth0 is the WAN and has this network, 192.168.56.1/24.  So
	// when something shows up from there, we disregard it.
	localInterfaces[0].On("Addrs").Return(
		[]net.Addr{
			wanNet,
		})
	predicate := IsNotFromWANDevice(mfwInterfaces, localInterfacesSysNet)
	suite.False(
		predicate(&DeviceEntry{}))
}

func TestDeviceList(t *testing.T) {
	testSuite := &DeviceListTestSuite{}
	suite.Run(t, testSuite)
}

package interfaces

import (
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
)

var (
	settingsPath string
	settingsFile *settings.SettingsFile
	addr0        string
	mask0        uint8
	device0      string
	addr1        string
	mask1        uint8
)

func netTestSetup(t *testing.T) func(t *testing.T) {
	settingsPath = "./testdata/settings.json"
	settingsFile = settings.NewSettingsFile(settingsPath)
	addr0 = "192.168.56.1"
	mask0 = 24
	device0 = "eth0"
	addr1 = "192.168.56.10"
	mask1 = 24

	return func(t *testing.T) {
		// shutdown
	}
}

func TestGetLocalInterfaceFromIp(t *testing.T) {
	tearDownSuite := netTestSetup(t)
	defer tearDownSuite(t)

	ip, _, parseErr := net.ParseCIDR(addr1 + "/" + strconv.Itoa(int(mask1)))
	if parseErr == nil {
		intf, getErr := GetLocalInterfaceFromIpAndPath(ip, settingsPath)
		if getErr == nil {
			assert.NotNil(t, intf)
			assert.NotNil(t, intf.Device)
			assert.Equal(t, device0, intf.Device)
		}
	}
}

func TestSettings(t *testing.T) {
	tearDownSuite := netTestSetup(t)
	defer tearDownSuite(t)

	var interfaces []Interface
	err := settingsFile.UnmarshalSettingsAtPath(&interfaces, "network", "interfaces")
	assert.Nil(t, err)
	assert.Equal(
		t,
		interfaces,
		[]Interface{
			{
				IsWAN:           false,
				Enabled:         true,
				V4StaticAddress: addr0,
				V4StaticPrefix:  mask0,
				Device:          device0,
			},
			{
				IsWAN:           false,
				Enabled:         true,
				V4StaticAddress: "asdf",
				V4StaticPrefix:  0,
				Device:          "qwer",
			},
		},
	)
}

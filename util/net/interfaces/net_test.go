package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
)

var (
	testSettingsPath   string
	testSettingsFile   *settings.SettingsFile
	expectedInterfaces []Interface
)

func netTestSetup(t *testing.T) func(t *testing.T) {
	testSettingsPath = "./testdata/settings.json"
	testSettingsFile = settings.NewSettingsFile(testSettingsPath)
	expectedInterfaces = []Interface{
		{
			IsWAN:           false,
			Enabled:         true,
			V4StaticAddress: "192.168.56.1",
			V4StaticPrefix:  24,
			Device:          "eth0",
		},
		{
			IsWAN:           false,
			Enabled:         true,
			V4StaticAddress: "asdf",
			V4StaticPrefix:  0,
			Device:          "qwer",
		},
	}

	return func(t *testing.T) {
		// shutdown
	}
}

func TestGetLocalInterfacesFromPath(t *testing.T) {
	tearDownSuite := netTestSetup(t)
	defer tearDownSuite(t)

	intfSettings := InterfaceSettings{
		file:     settings.NewSettingsFile(testSettingsPath),
		jsonPath: []string{},
	}
	intfSettings.SetJsonPath(defaultJsonParent, defaultJsonChild)
	interfaces := GetLocalInterfacesFromSettings(intfSettings)
	assert.Equal(t, expectedInterfaces, interfaces)
}

func TestSettings(t *testing.T) {
	tearDownSuite := netTestSetup(t)
	defer tearDownSuite(t)

	var interfaces []Interface
	err := testSettingsFile.UnmarshalSettingsAtPath(&interfaces, "network", "interfaces")
	assert.Nil(t, err)
	assert.Equal(t, interfaces, expectedInterfaces)
}

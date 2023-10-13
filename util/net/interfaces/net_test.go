package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	testSettingsPath = "./testdata/settings.json"
)

type interfacesTestFixture struct {
	settingsFile      *settings.SettingsFile
	interfaceSettings *InterfaceSettings

	lanOneExpected  Interface
	lanTwoExpected  Interface
	vlanOneExpected Interface
	vlanTwoExpected Interface
	wanOneExpected  Interface
}

// Setup objects used by all tests
func setupNewTestFixture() *interfacesTestFixture {
	f := &interfacesTestFixture{}
	f.lanOneExpected = Interface{
		IsWAN:           false,
		Enabled:         true,
		V4StaticAddress: "192.168.56.1",
		V4StaticPrefix:  24,
		Device:          "eth0",
	}
	f.lanTwoExpected = Interface{
		IsWAN:           false,
		Enabled:         true,
		V4StaticAddress: "asdf",
		V4StaticPrefix:  0,
		Device:          "qwer",
	}
	f.vlanOneExpected = Interface{
		BoundInterfaceID: 3,
		BridgedTo:        2,
		Enabled:          true,
		V4StaticAddress:  "192.168.59.2",
		V4StaticPrefix:   24,
		IsVirtual:        true,
		VlanID:           "2",
		IsWAN:            false,
	}
	f.vlanTwoExpected = Interface{
		BoundInterfaceID: 4,
		BridgedTo:        2,
		Enabled:          false,
		V4StaticAddress:  "192.168.88.2",
		V4StaticPrefix:   24,
		IsVirtual:        true,
		VlanID:           "3",
		IsWAN:            false,
	}
	f.wanOneExpected = Interface{
		IsWAN:           true,
		Enabled:         true,
		V4StaticAddress: "50.50.50.50",
		V4StaticPrefix:  8,
		Device:          "wan0",
	}

	f.settingsFile = settings.NewSettingsFile(testSettingsPath)
	f.interfaceSettings = NewInterfaceSettings(f.settingsFile)

	return f
}

// Test GetLocalInterfaces
func TestGetLocalInterfaces(t *testing.T) {
	f := setupNewTestFixture()

	actual := f.interfaceSettings.GetLocalInterfaces()
	expected := []Interface{f.lanOneExpected, f.lanTwoExpected, f.vlanOneExpected}

	assert.ElementsMatch(t, expected, actual)
}

// Test GetVLANInterfaces
func TestGetVLANInterfaces(t *testing.T) {
	f := setupNewTestFixture()

	actual := f.interfaceSettings.GetVLANInterfaces()
	expected := []Interface{f.vlanOneExpected}

	assert.ElementsMatch(t, expected, actual)
}

// Test GetVLANInterfaces
func TestGetAllInterfaces(t *testing.T) {
	f := setupNewTestFixture()

	actual := f.interfaceSettings.GetAllInterfaces()
	expected := []Interface{f.lanOneExpected, f.lanTwoExpected, f.vlanOneExpected, f.vlanTwoExpected, f.wanOneExpected}

	assert.ElementsMatch(t, expected, actual)
}

// Tests TestGetInterfacesWithFilter
func TestGetInterfacesWithFilter(t *testing.T) {
	f := setupNewTestFixture()
	interfaceSettings := NewInterfaceSettings(f.settingsFile)

	testCases := []struct {
		name       string
		filterFunc func(Interface) bool
		expected   []Interface
	}{
		{
			name:       "All Enabled",
			filterFunc: (func(i Interface) bool { return i.Enabled }),
			expected:   []Interface{f.lanOneExpected, f.lanTwoExpected, f.vlanOneExpected, f.wanOneExpected},
		},
		{
			name:       "All Disabled",
			filterFunc: (func(i Interface) bool { return !i.Enabled }),
			expected:   []Interface{f.vlanTwoExpected},
		},
		{
			name:       "All WANs",
			filterFunc: (func(i Interface) bool { return i.IsWAN }),
			expected:   []Interface{f.wanOneExpected},
		},
		{
			name:       "All VLANs",
			filterFunc: (func(i Interface) bool { return i.IsVirtual }),
			expected:   []Interface{f.vlanOneExpected, f.vlanTwoExpected},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := interfaceSettings.GetInterfacesWithFilter(tc.filterFunc)
			if !assert.ElementsMatch(t, tc.expected, actual) {
				t.Errorf("Test(%s) failed, did not receive expected elements\n", tc.name)
			}
		})
	}
}

// Test if the interface struct works as expected with UnmarshalSettingsAtPath
func TestUnmarshalInterfaceFromSettings(t *testing.T) {
	f := setupNewTestFixture()
	expected := []Interface{f.lanOneExpected, f.lanTwoExpected, f.vlanOneExpected, f.vlanTwoExpected, f.wanOneExpected}

	var actual []Interface
	err := f.settingsFile.UnmarshalSettingsAtPath(&actual, "network", "interfaces")
	assert.Nil(t, err)
	assert.ElementsMatch(t, actual, expected)
}

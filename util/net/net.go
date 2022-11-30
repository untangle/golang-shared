package interfaces

import (
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

type InterfaceSettings struct {
	file     *settings.SettingsFile
	jsonPath []string
}

func (intfSettings *InterfaceSettings) UnmarshalJson(interfaces *[]Interface) error {
	return intfSettings.file.UnmarshalSettingsAtPath(&interfaces, intfSettings.GetJsonPath()...)
}

func (intfSettings *InterfaceSettings) GetJsonPath() []string {
	return intfSettings.jsonPath
}

func (intfSettings *InterfaceSettings) SetJsonPath(jsonPath ...string) {
	intfSettings.jsonPath = jsonPath
}

const (
	defaultSettingsFile = "/etc/config/settings.json"
	defaultJsonParent   = "network"
	defaultJsonChild    = "interfaces"
)

// GetInterfaces returns a list of interfaces, filtered by any propeties passed in
// @param - filter func(InterfaceDetail) bool - a function filter to filter results if needed
// @return - []InterfaceDetail - an array of InterfaceDetail types
func GetInterfaces(intfSettings InterfaceSettings, filter func(Interface) bool) []Interface {
	var interfaces []Interface
	err := intfSettings.UnmarshalJson(&interfaces)
	if err != nil {
		logger.Warn("Unable to read network settings: %s\n", err.Error())
		return nil
	}

	if filter != nil {
		var filteredInterfaces []Interface
		for _, intf := range interfaces {
			if filter(intf) {
				logger.Warn("Adding intf %s\n", intf.Device)
				filteredInterfaces = append(filteredInterfaces, intf)
			}
		}

		return filteredInterfaces
	} else {
		return interfaces
	}
}

// Returns local interfaces. That is, those that aren't a WAN, are enabled,
// and have either an IPv4 or IPv6 address
func GetLocalInterfacesFromSettings(intfSettings InterfaceSettings) []Interface {
	return GetInterfaces(intfSettings, (func(intf Interface) bool {
		hasIp := intf.V4StaticAddress != "" || intf.V6StaticAddress != ""
		return !intf.IsWAN && intf.Enabled && hasIp
	}))
}

// Calls GetLocalInterfacesFromPath with default settings.json path
func GetLocalInterfaces() []Interface {
	return GetLocalInterfacesFromSettings(GetDefaultInterfaceSettings())
}

func GetDefaultInterfaceSettings() InterfaceSettings {
	intfSettings := InterfaceSettings{
		file:     settings.NewSettingsFile(defaultSettingsFile),
		jsonPath: []string{},
	}
	intfSettings.SetJsonPath(defaultJsonParent, defaultJsonChild)
	return intfSettings
}

func GetInterfaceSettingsFromPath(path string) InterfaceSettings {
	intfSettings := InterfaceSettings{
		file:     settings.NewSettingsFile(path),
		jsonPath: []string{},
	}
	intfSettings.SetJsonPath(defaultJsonParent, defaultJsonChild)
	return intfSettings
}

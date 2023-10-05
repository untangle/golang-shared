package interfaces

import (
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

// InterfaceSettings is an object for manipulating the network/interfaces list in the
// settings.json file.
type InterfaceSettings struct {
	// settings file to read settings from.
	file *settings.SettingsFile

	// path within the settings.json file for the interfaces.
	jsonPath []string
}

func (ifaces *InterfaceSettings) UnmarshalJson(interfaces *[]Interface) error {
	return ifaces.file.UnmarshalSettingsAtPath(&interfaces, ifaces.GetJsonPath()...)
}

func (ifaces *InterfaceSettings) GetJsonPath() []string {
	return ifaces.jsonPath
}

func (ifaces *InterfaceSettings) SetJsonPath(jsonPath ...string) {
	ifaces.jsonPath = jsonPath
}

const (
	defaultSettingsFile = "/etc/config/settings.json"
	defaultJsonParent   = "network"
	defaultJsonChild    = "interfaces"
)

// NewInterfaceSettings returns an InterfaceSettings object which uses
// the file as it's underlying settings file.
func NewInterfaceSettings(file *settings.SettingsFile) *InterfaceSettings {
	return &InterfaceSettings{file: file, jsonPath: []string{defaultJsonParent, defaultJsonChild}}
}

// GetInterfacesWithFilter returns a slice of all interfaces for which
// the filter function returns true.
func (ifaces *InterfaceSettings) GetInterfacesWithFilter(filter func(Interface) bool) (interfaces []Interface) {
	err := ifaces.UnmarshalJson(&interfaces)
	if err != nil {
		logger.Warn("Unable to read network settings: %s\n", err.Error())
		return nil
	}

	if filter != nil {
		var filteredInterfaces []Interface
		for _, intf := range interfaces {
			if filter(intf) {
				filteredInterfaces = append(filteredInterfaces, intf)
			}
		}

		return filteredInterfaces
	} else {
		return interfaces
	}
}

// GetAllInterfaces returns all interfaces in the ifaces object.
func (ifaces *InterfaceSettings) GetAllInterfaces() []Interface {
	return ifaces.GetInterfacesWithFilter(func(Interface) bool { return true })
}

// GetLocalInterfaces returns all interfaces that are not WAN, are
// enabled, and have an assigned IP address of some kind.
func (ifaces *InterfaceSettings) GetLocalInterfaces() (interfaces []Interface) {
	return ifaces.GetInterfacesWithFilter(
		func(intf Interface) bool {
			hasIP := intf.V4StaticAddress != "" || intf.V6StaticAddress != ""
			return !intf.IsWAN && intf.Enabled && hasIP
		})

}

// GetVLANInterfaces returns all interfaces that are VLANs
func (ifaces *InterfaceSettings) GetVLANInterfaces() (interfaces []Interface) {
	return ifaces.GetInterfacesWithFilter(
		func(intf Interface) bool {
			return intf.Enabled && intf.IsVirtual && !intf.IsWAN
		})
}

// GetInterfaces returns a list of interfaces, filtered by any propeties passed in
// @param - filter func(InterfaceDetail) bool - a function filter to filter results if needed
// @return - []InterfaceDetail - an array of InterfaceDetail types
func GetInterfaces(intfSettings InterfaceSettings, filter func(Interface) bool) (interfaces []Interface) {
	return intfSettings.GetInterfacesWithFilter(filter)
}

// Returns local interfaces. That is, those that aren't a WAN, are enabled,
// and have either an IPv4 or IPv6 address
func GetLocalInterfacesFromSettings(intfSettings InterfaceSettings) []Interface {
	return intfSettings.GetLocalInterfaces()
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

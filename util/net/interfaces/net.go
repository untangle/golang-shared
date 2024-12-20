package interfaces

import (
	"net"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

type InterfaceFilter func(Interface) bool

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
	defaultJsonParent = "network"
	defaultJsonChild  = "interfaces"
)

// NewInterfaceSettings returns an InterfaceSettings object which uses
// the file as it's underlying settings file.
func NewInterfaceSettings(file *settings.SettingsFile) *InterfaceSettings {
	return &InterfaceSettings{file: file, jsonPath: []string{defaultJsonParent, defaultJsonChild}}
}

// GetInterfacesWithFilter returns a slice of all interfaces for which
// the filter function returns true.
func (ifaces *InterfaceSettings) GetInterfacesWithFilter(filter InterfaceFilter) (interfaces []Interface) {
	if err := ifaces.UnmarshalJson(&interfaces); err != nil {
		logger.Warn("Unable to read network settings: %s\n", err.Error())
		return nil
	}

	// Don't bother with filtering if no filter was given
	if filter == nil {
		return interfaces
	}

	var filteredInterfaces []Interface
	for _, intf := range interfaces {
		if filter(intf) {
			filteredInterfaces = append(filteredInterfaces, intf)
		}
	}

	return filteredInterfaces
}

// GetAllInterfaces returns all interfaces in the ifaces object.
func (ifaces *InterfaceSettings) GetAllInterfaces() []Interface {
	return ifaces.GetInterfacesWithFilter(func(Interface) bool { return true })
}

// GetLocalInterfaces returns all interfaces that are not WAN, are
// enabled, and have an assigned IP address of some kind.
// We consider mgmt interfaces to be local interfaces. Caller should
// filter out mgmt interfaces if needed.
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
		GetVLANFilter())
}

// Returns an InterfaceFilter function used to get VLANs
func GetVLANFilter() InterfaceFilter {
	return func(intf Interface) bool {
		return intf.Enabled && intf.VlanID != ""
	}
}

// GetInterfaces returns a list of interfaces, filtered by any propeties passed in
// @param - filter func(InterfaceDetail) bool - a function filter to filter results if needed
// @return - []InterfaceDetail - an array of InterfaceDetail types
func GetInterfaces(intfSettings InterfaceSettings, filter InterfaceFilter) (interfaces []Interface) {
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
	settingsFile, err := settings.GetSettingsFileSingleton()
	if err != nil {
		logger.Err(err.Error())
	}
	intfSettings := InterfaceSettings{
		file:     settingsFile,
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

// Check Ip belong to device local IP address
func IsDeviceLocalIp(ip net.IP) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && ip.Equal(ipNet.IP) {
			return true
		}
	}
	return false
}

package interfaces

import (
	"net"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const defaultSettingsFile = "/etc/config/settings.json"

// GetInterfaces returns a list of interfaces, filtered by any propeties passed in
// @param - filter func(InterfaceDetail) bool - a function filter to filter results if needed
// @return - []InterfaceDetail - an array of InterfaceDetail types
func GetInterfaces(path string, filter func(Interface) bool) []Interface {
	var interfaces []Interface
	// Maybe use SettingsFile here instead of settings. Then this function takes a file location as a parameter. Then for testing can skip "getlocalinterfaces" and call your mocked settings.json directly?
	settingsFile := settings.NewSettingsFile(path)
	err := settingsFile.UnmarshalSettingsAtPath(&interfaces, "network", "interfaces")
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

// Returns local interfaces. That is, those that aren't a WAN, are enabled,
// and have either an IPv4 or IPv6 address
func GetLocalInterfacesFromPath(path string) []Interface {
	return GetInterfaces(path, (func(intf Interface) bool {
		hasIp := intf.V4StaticAddress != "" || intf.V6StaticAddress != ""
		return !intf.IsWAN && intf.Enabled && hasIp
	}))
}

// Calls above function, GetLocalInterfacesFromPath, assuming default settings.json path
func GetLocalInterfaces() []Interface {
	return GetLocalInterfacesFromPath(defaultSettingsFile)
}

// Grabs a single local interface from an IP. If the passed IP is within the
// interface's network, that interface is returned. Otherwise an error is
// returned.
func GetLocalInterfaceFromIpAndPath(addr net.IP, path string) (*Interface, error) {
	currMask := 0
	localInterfaces := GetLocalInterfacesFromPath(path)
	var currIntf *Interface
	// we're assigning to localInterfaces at end of loop, so range won't work
	for i := 0; i < len(localInterfaces); i++ {
		intf := localInterfaces[i]
		// Grab first subnet mask from interface
		network, err := intf.GetNetwork()
		// don't need to return on error. Some interfaces have no network.
		if err == nil {
			ones, _ := network.Mask.Size()
			if network.Contains(addr) && ones > currMask {
				// if we have a mask, compare to current. Keep larger.
				currMask = ones
				currIntf = &intf
			}
		}
	}
	return currIntf, nil
}

// Calls above function, GetLocalInterfaceFromIpAndPath, assuming default settings.json path
func GetLocalInterfaceFromIp(addr net.IP) (*Interface, error) {
	return GetLocalInterfaceFromIpAndPath(addr, defaultSettingsFile)
}

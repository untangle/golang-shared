package interfaces

import (
	"fmt"
	"net"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

// GetInterfaces returns a list of interfaces, filtered by any propeties passed in
// @param - filter func(InterfaceDetail) bool - a function filter to filter results if needed
// @return - []InterfaceDetail - an array of InterfaceDetail types
func GetInterfaces(filter func(Interface) bool) []Interface {
	var interfaces []Interface
	if err := settings.UnmarshalSettingsAtPath(&interfaces, "network", "interfaces"); err != nil {
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
func GetLocalInterfaces() []Interface {
	return GetInterfaces((func(intf Interface) bool {
		return !intf.IsWAN && intf.Enabled && (intf.V4StaticAddress != "" || intf.V6StaticAddress != "")
	}))
}

// Grabs a single local interface from an IP. If the passed IP is within the
// interface's network, that interface is returned. Otherwise an error is
// returned.
func GetLocalInterfaceFromIp(addr net.IP) (Interface, error) {
	var intf Interface
	for _, intf = range GetLocalInterfaces() {
		if intf.NetworkHasIP(addr) {
			return intf, nil
		}
	}
	return intf, fmt.Errorf("address '%s' not in local interfaces", addr)
}

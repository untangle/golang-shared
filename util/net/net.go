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

func GetLocalInterfaces() []Interface {
	return GetInterfaces((func(intf Interface) bool {
		return !intf.IsWAN && intf.Enabled && intf.V4StaticAddress != ""
	}))
}

func GetLocalInterfaceFromIp(cidrAddr net.IP) (Interface, error) {
	var intf Interface
	for _, intf = range GetLocalInterfaces() {
		if intf.NetworkHasIP(cidrAddr) {
			return intf, nil
		}
	}
	return intf, fmt.Errorf("CIDR Address %s not in local interfaces", cidrAddr)
}

package interfaces

import (
	"fmt"
	"strconv"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

// GetInterfaces returns a list of interfaces, filtered by any propeties passed in
// @param - filter func(InterfaceDetail) bool - a function filter to filter results if needed
// @return - []InterfaceDetail - an array of InterfaceDetail types
func GetInterfaces(filter func(Interface) bool) []Interface {
	var interfaces []Interface
	if err := settings.UnmarshalSettingsAtPath(&interfaces, "network", "interfaces"); err != nil {
		logger.Err("chap chap\n\n\n")
		logger.Warn("Unable to read network settings: %s\n", err.Error())
		return nil
	}

	logger.Info("uggggggggggggggggggggggggggggggggggggggggggggggggggghhhhhhhhhhhhh\n")
	for _, str := range interfaces {
		logger.Err("yoooo %s", str)
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

func GetLocalNetworks(localIntfs []Interface) []string {
	var localNetworks []string = nil

	for _, intf := range localIntfs {
		prefix := strconv.FormatFloat(intf.V4StaticPrefix, 'f', -1, 64)
		localNetwork := fmt.Sprintf("%s/%s", intf.V4StaticAddress, prefix)
		logger.Debug("Found local network %s\n", localNetwork)
		localNetworks = append(localNetworks, localNetwork)
	}

	return localNetworks
}

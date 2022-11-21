package interfaces

import (
	"fmt"
	"net"
)

// Interface corresponds to the interface JSON structure in the settings.json
// file.
type Interface struct {
	ConfigType        string `json:"configType"`
	Device            string `json:"device"`
	DHCPEnabled       bool
	DHCPLeaseDuration int
	DHCPOptions       interface{}
	DHCPRangeStart    string
	DHCPRangeEnd      string
	DownloadKbps      int
	Enabled           bool   `json:"enabled"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	V4StaticAddress   string `json:"v4StaticAddress"`
	V4StaticPrefix    uint8  `json:"v4StaticPrefix"`
	IsWAN             bool   `json:"wan"`
}

func (intf *Interface) GetCidrNotation() string {
	return fmt.Sprintf("%s/%d", intf.V4StaticAddress, intf.V4StaticPrefix)
}

func (intf *Interface) GetNetwork() (net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(intf.GetCidrNotation())
	return *ipNet, err
}

func (intf *Interface) NetworkHasIP(ip net.IP) bool {
	_, intfNetwork, err := net.ParseCIDR(intf.GetCidrNotation())
	if err == nil && intfNetwork.Contains(ip) {
		return true
	}
	return false
}

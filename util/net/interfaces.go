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
	V6StaticAddress   string `json:"v6StaticAddress"`
	v6StaticPrefix    uint8  `json:"v6StaticPrefix"`
	IsWAN             bool   `json:"wan"`
}

// Get CIDR notation using the IP address and static prefix. If there's an
// IPv4 address, that is returned. If not, an IPv6 address is grabbed. If
// neither exists, an error is returned.
func (intf *Interface) GetCidrNotation() (string, error) {
	if intf.V4StaticAddress != "" {
		return fmt.Sprintf("%s/%d", intf.V4StaticAddress, intf.V4StaticPrefix), nil
	} else if intf.V6StaticAddress != "" {
		return fmt.Sprintf("%s/%d", intf.V6StaticAddress, intf.v6StaticPrefix), nil
	} else {
		return "", fmt.Errorf("interface '%s' does not have a V4StaticAddress or V6StaticAddress", intf.Name)
	}
}

// Returns the net.IPNet object corresponding to the interface obtained
// using GetCidrNotation()
func (intf *Interface) GetNetwork() (*net.IPNet, error) {
	cidr, err := intf.GetCidrNotation()
	if err == nil {
		_, ipNet, err := net.ParseCIDR(cidr)
		return ipNet, err
	} else {
		var ipNet net.IPNet
		return &ipNet, err
	}
}

// Checks if a given net.IP is within the interface's network. Used to map
// IP addresses to interface's for discovery
func (intf *Interface) NetworkHasIP(ip net.IP) bool {
	intfNetwork, err := intf.GetNetwork()
	if err == nil && intfNetwork.Contains(ip) {
		return true
	}
	return false
}

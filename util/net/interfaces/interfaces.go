package interfaces

import (
	"fmt"
	"net"
)

const (
	// ConfigTypeBridged - If the interface is a bridged
	// interface, this is it's ConfigType.
	ConfigTypeBridged = "BRIDGED"
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

	Enabled         bool          `json:"enabled"`
	Name            string        `json:"name"`
	InterfaceID     int           `json:"interfaceID"`
	Type            string        `json:"type"`
	V4StaticAddress string        `json:"v4StaticAddress"`
	V4StaticPrefix  uint8         `json:"v4StaticPrefix"`
	V4Aliases       []V4IpAliases `json:"v4Aliases"`
	V6StaticAddress string        `json:"v6StaticAddress"`
	V6StaticPrefix  uint8         `json:"v6StaticPrefix"`
	V6Aliases       []V6IpAliases `json:"v6Aliases"`
	IsWAN           bool          `json:"wan"`

	// VLAN specific settings
	BoundInterfaceID uint16 `json:"boundInterfaceId"`
	BridgedTo        uint16 `json:"bridgedTo"`
	IsVirtual        bool   `json:"virtual"`
	VlanID           string `json:"vlanid"`
}
type V4IpAliases struct {
	V4Address string `json:"v4Address"`
	V4Prefix  uint32 `json:"v4Prefix"`
}
type V6IpAliases struct {
	V6Address string `json:"v6Address"`
	V6Prefix  string `json:"v6Prefix"`
}

func (intf *Interface) IsBridged() bool {
	return intf.ConfigType == ConfigTypeBridged
}

// Get IPV4 static and aliases addresses
func (intf *Interface) GetIpV4Network() []*net.IPNet {
	var networks []*net.IPNet
	if intf.IsBridged() {
		return networks
	}
	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", intf.V4StaticAddress, intf.V4StaticPrefix))
	if err == nil {
		networks = append(networks, ipNet)
	}
	for _, alias := range intf.V4Aliases {
		_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", alias.V4Address, alias.V4Prefix))
		if err == nil {
			networks = append(networks, ipNet)
		}
	}
	return networks
}

// Get IPV6 static and aliases addresses
func (intf *Interface) GetIpV6Network() []*net.IPNet {
	var networks []*net.IPNet
	if intf.IsBridged() {
		return networks
	}
	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", intf.V6StaticAddress, intf.V6StaticPrefix))
	if err == nil {
		networks = append(networks, ipNet)
	}
	for _, alias := range intf.V6Aliases {
		_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%s", alias.V6Address, alias.V6Prefix))
		if err == nil {
			networks = append(networks, ipNet)
		}
	}
	return networks
}

// Get IPV4 and IPV6 static addresses and aliases addresses
func (intf *Interface) GetNetworks() []*net.IPNet {
	var networks []*net.IPNet
	ipV4Nets := intf.GetIpV4Network()
	if len(ipV4Nets) != 0 {
		networks = append(networks, ipV4Nets...)
	}
	ipV6Nets := intf.GetIpV6Network()
	if len(ipV6Nets) != 0 {
		networks = append(networks, ipV6Nets...)
	}
	return networks
}

// MostSpecificPrefixMatch returns from the networks and addr the
// network with the longest matching prefix.
func MostSpecificPrefixMatch(nets []*net.IPNet, addr net.IP) *net.IPNet {
	currMask := 0
	var maxMatching *net.IPNet = nil

	for _, network := range nets {
		// do most specific prefix match on networks belonging to this interface
		ones, _ := network.Mask.Size()
		if network.Contains(addr) && ones > currMask {
			// if we have a mask, compare to current. Keep larger.
			currMask = ones
			maxMatching = network
		}
	}
	return maxMatching
}

func (intf *Interface) HasContainingNetwork(addr net.IP) *net.IPNet {
	return MostSpecificPrefixMatch(intf.GetNetworks(), addr)
}

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
	V6StaticPrefix    uint8  `json:"v6StaticPrefix"`
	IsWAN             bool   `json:"wan"`
}

func (intf *Interface) GetIpV4Network() (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", intf.V4StaticAddress, intf.V4StaticPrefix))
	return ipNet, err
}

func (intf *Interface) GetIpV6Network() (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", intf.V6StaticAddress, intf.V6StaticPrefix))
	return ipNet, err
}

func (intf *Interface) GetNetworks() []net.IPNet {
	var networks []net.IPNet
	ipV4Net, v4Err := intf.GetIpV4Network()
	if v4Err == nil {
		networks = append(networks, *ipV4Net)
	}
	ipV6Net, v6Err := intf.GetIpV6Network()
	if v6Err == nil {
		networks = append(networks, *ipV6Net)
	}
	return networks
}

func (intf *Interface) HasContainingNetwork(addr net.IP) net.IPNet {
	currMask := 0
	var maxMatching net.IPNet
	for _, network := range intf.GetNetworks() {
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

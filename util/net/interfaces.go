package interfaces

import (
	"fmt"
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
	Enabled           bool        `json:"enabled"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	Aliases           []IpAliases `json:"v4Aliases"`
	V4StaticAddress   string      `json:"v4StaticAddress"`
	V4StaticPrefix    uint8       `json:"v4StaticPrefix"`
	IsWAN             bool        `json:"wan"`
}
type IpAliases struct {
	V4Address string `json:"v4Address"`
	V4Prefix  uint32 `json:"v4Prefix"`
}

func (intf *Interface) GetCidrNotation() string {
	return fmt.Sprintf("%s/%d", intf.V4StaticAddress, intf.V4StaticPrefix)
}

func (alias *IpAliases) GetCidrNotation() string {
	return fmt.Sprintf("%s/%d", alias.V4Address, alias.V4Prefix)
}

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

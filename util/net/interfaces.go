package interfaces

import (
	"fmt"
	"strconv"
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
	Enabled           bool    `json:"enabled"`
	Name              string  `json:"name"`
	Type              string  `json:"type"`
	V4StaticAddress   string  `json:"v4StaticAddress"`
	V4StaticPrefix    float64 `json:"v4StaticPrefix"`
	IsWAN             bool    `json:"wan"`
}

func (intf *Interface) GetCidrNotation() string {
	prefix := strconv.FormatFloat(intf.V4StaticPrefix, 'f', -1, 64)
	return fmt.Sprintf("%s/%s", intf.V4StaticAddress, prefix)
}

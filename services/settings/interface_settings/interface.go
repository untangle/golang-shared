package interface_settings

type Interface struct {
	/* ConfigType        string `json:"configType"`
	Device            string `json:"device"`
	DhcpEnabled       bool   `json:"dhcpEnabled"`
	DhcpLeaseDuration int    `json:"dhcpLeaseDuration"`
	DhcpRangeEnd      string `json:"dhcpRangeEnd"`
	DhcpRangeStart    string `json:"dhcpRangeStart"`
	DhcpRelayEnabled  bool   `json:"dhcpRelayEnabled"`
	DownloadKbps      int    `json:"downloadKbps"`
	Enabled           bool   `json:"enabled"`
	EthAutoneg        bool   `json:"ethAutoneg"`
	EthDuplex         string `json:"ethDuplex"`
	EthSpeed          int    `json:"ethSpeed"`
	Mtu               int    `json:"mtu"`
	Name              string `json:"name"`
	NatIngress        bool   `json:"natIngress"`
	QosEnabled        bool   `json:"qosEnabled"`
	Type              string `json:"type"`
	UploadKbps        int    `json:"uploadKbps"`
	V4ConfigType      string `json:"v4ConfigType"`
	V4StaticAddress   string `json:"v4StaticAddress"`
	V4StaticPrefix    int    `json:"v4StaticPrefix"`
	V6AssignHint      string `json:"v6AssignHint"`
	V6AssignPrefix    int    `json:"v6AssignPrefix"`
	V6ConfigType      string `json:"v6ConfigType"`
	Virtual           bool   `json:"virtual"`
	Wan               bool   `json:"wan"`
	WanWeight         int    `json:"wanWeight"`
	InterfaceId       int    `json:"interfaceId"`
	DhcpRelayAddress  string `json:"dhcpRelayAddress"` */
	V4PPPoEPassword string `json:"v4PPPoEPassword"`
}

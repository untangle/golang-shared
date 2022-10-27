package pluginloads

// Register plugins by making sure their init() functions are run
import (
	_ "github.com/untangle/discoverd/plugins/arp"
	_ "github.com/untangle/discoverd/plugins/discovery"
	_ "github.com/untangle/discoverd/plugins/lldp"
	_ "github.com/untangle/discoverd/plugins/nmap"
)

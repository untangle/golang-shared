package discovery_settings

import (
	"github.com/untangle/golang-shared/services/discovery"
)

// base for collector settings, inherited by all collector settings
type CollectorSettingsBase struct {
	Type         discovery.CollectorName `json:"type"`
	Enabled      bool                    `json:"enabled"`
	AutoInterval uint                    `json:"autoInterval"`
}

// settings for lldp collector
type LldpSettings struct {
	CollectorSettingsBase
}

// settings for neighbour collector
type NeighbourSettings struct {
	CollectorSettingsBase
}

// settings for nmap collector
type NmapSettings struct {
	CollectorSettingsBase
}

// settings for discovery plugin
type DiscoveryPluginSettings struct {
	Enabled bool `json:"enabled"`
}

// whole settings object, contains discovery settings and an array of settings for individual collectors
type discoverySettingsObject struct {
	DiscoveryPluginSettings
	Plugins []interface{} `json:"plugins"`
}

type iCollectorSettings interface {
	IsValid() bool
	readBytes(bytes []byte) bool
}

package settings

import (
	"encoding/json"

	"github.com/untangle/golang-shared/services/discovery"
)

const (
	// keep those in sync with the front end validation
	minAutoIncrement uint = 60            // 1 hour, 60 minutes
	maxAutoIncrement uint = 365 * 24 * 60 // 365 days * 24 hours * 60 minutes
)

type CollectorSettingsBase struct {
	Type         discovery.CollectorName `json:"type"`
	Enabled      bool                    `json:"enabled"`
	AutoInterval uint                    `json:"autoInterval"`
}

type LldpSettings struct {
	CollectorSettingsBase
}

type NeighbourSettings struct {
	CollectorSettingsBase
}

type NmapSettings struct {
	CollectorSettingsBase
}

type DiscoveryPluginSettings struct {
	Enabled bool `json:"enabled"`
}

type discoverySettingsObject struct {
	DiscoveryPluginSettings
	Plugins []interface{} `json:"plugins"`
}

// ValidateDiscoverySettings - ensures the settings we received are in a valid format
//  returns true if the object is valid, false otherwise
func ValidateDiscoverySettings(settingsObjBytes []byte) bool {

	// first we unmarshal the whole settings object
	discoverySettings := discoverySettingsObject{}
	if err := json.Unmarshal(settingsObjBytes, &discoverySettings); err != nil {
		logger.Info("ValidateDiscoverySettings unable to unmarshall discoverySettingsObject with err %v\n", err)
		return false
	}

	// then we check each element in the plugins array
	for _, iPlugin := range discoverySettings.Plugins {
		if !validateOneCollector(iPlugin) {
			return false
		}
	}

	return true
}

// validateOneCollector - validates the settings for individual collectors
func validateOneCollector(pluginInterface interface{}) bool {
	pluginBytes, err := json.Marshal(pluginInterface)
	if err != nil {
		logger.Info("ValidateDiscoverySettings unable to marshal plugin settings interface with err %v\n", err)
		return false
	}

	// convert to CollectorSettingsBase so we can get the plugin type and cast it accordingly
	basePlugin := CollectorSettingsBase{}
	if err := json.Unmarshal(pluginBytes, &basePlugin); err != nil {
		logger.Info("ValidateDiscoverySettings could not unmarshal basePlugin settings with err %v\n", err)
		return false
	}

	// depending on the collector type we try to cast to the actual collector settings struct and validate it
	switch basePlugin.Type {
	case discovery.Arp: // settings are still saved under "arp" field
		neighbourSettings := NeighbourSettings{}
		if err := json.Unmarshal(pluginBytes, &neighbourSettings); err != nil {
			logger.Info("ValidateDiscoverySettings could not unmarshal neighbour settings with err %v\n", err)
			return false
		}
		if !neighbourSettings.IsValid() {
			return false
		}
	case discovery.Lldp:
		lldpSettings := LldpSettings{}
		if err := json.Unmarshal(pluginBytes, &lldpSettings); err != nil {
			logger.Info("ValidateDiscoverySettings could not unmarshal lldp settings with err %v\n", err)
			return false
		}
		if !lldpSettings.IsValid() {
			return false
		}
	case discovery.Nmap:
		nmapSettings := NmapSettings{}
		if err := json.Unmarshal(pluginBytes, &nmapSettings); err != nil {
			logger.Info("ValidateDiscoverySettings could not unmarshal nmap settings with err %v\n", err)
			return false
		}
		if !nmapSettings.IsValid() {
			return false
		}

	default:
		logger.Info("ValidateDiscoverySettings received unknown settings type %v\n", basePlugin.Type)
		return false
	}

	return true
}

func (s *LldpSettings) IsValid() bool {
	if !validateAutoInterval(s.AutoInterval) {
		logger.Info("ValidateDiscoverySettings LldpSettings AutoInterval should be between %v and %v but is %v\n", minAutoIncrement, maxAutoIncrement, s.AutoInterval)
		return false
	}
	if s.Type != discovery.Lldp {
		logger.Info("ValidateDiscoverySettings LldpSettings wrong type %v\n", s.Type)
		return false
	}
	return true
}

func (s *NeighbourSettings) IsValid() bool {
	if !validateAutoInterval(s.AutoInterval) {
		logger.Info("ValidateDiscoverySettings NeighbourSettings AutoInterval should be between %v and %v but is %v\n", minAutoIncrement, maxAutoIncrement, s.AutoInterval)
		return false
	}
	if s.Type != discovery.Arp { // settings are still saved under "arp" field
		logger.Info("ValidateDiscoverySettings NeighbourSettings wrong type %v\n", s.Type)
		return false
	}
	return true
}

func (s *NmapSettings) IsValid() bool {
	if !validateAutoInterval(s.AutoInterval) {
		logger.Info("ValidateDiscoverySettings NmapSettings AutoInterval should be between %v and %v but is %v\n", minAutoIncrement, maxAutoIncrement, s.AutoInterval)
		return false
	}
	if s.Type != discovery.Nmap {
		logger.Info("ValidateDiscoverySettings NmapSettings wrong type %v\n", s.Type)
		return false
	}
	return true
}

func validateAutoInterval(autoIntervalMinutes uint) bool {
	return minAutoIncrement <= autoIntervalMinutes && autoIntervalMinutes <= maxAutoIncrement
}

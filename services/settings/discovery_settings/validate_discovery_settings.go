package discovery_settings

import (
	"encoding/json"
	"os"

	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
)

const (
	// keep those in sync with the front end validation
	minAutoIncrement uint = 60            // 1 hour, 60 minutes
	maxAutoIncrement uint = 365 * 24 * 60 // 365 days * 24 hours * 60 minutes
)

// ValidateDiscoverySettings - ensures the settings we received are in a valid format
//  returns true if the object is valid, false otherwise
func ValidateDiscoverySettings(settingsObjBytes []byte) bool {
	err := os.WriteFile("/tmp/testBody.json", settingsObjBytes, 0655)
	logger.Err("%s\n", err)
	logger.Err("%v\n", string(settingsObjBytes))
	// first we unmarshal the whole discovery settings object
	discoverySettings := discoverySettingsObject{}
	if err := json.Unmarshal(settingsObjBytes, &discoverySettings); err != nil {
		logger.Err("ValidateDiscoverySettings unable to unmarshall discoverySettingsObject with err %v\n", err)
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
// returns true if the object is valid, false otherwise
func validateOneCollector(settingsInterface interface{}) bool {
	settingsBytes, err := json.Marshal(settingsInterface)
	if err != nil {
		logger.Info("validateOneCollector unable to marshal plugin settings interface with err %v\n", err)
		return false
	}

	// convert to CollectorSettingsBase so we can get the plugin type and cast it accordingly
	baseSettings := CollectorSettingsBase{}
	if err := json.Unmarshal(settingsBytes, &baseSettings); err != nil {
		logger.Info("validateOneCollector could not unmarshal CollectorSettingsBase settings with err %v\n", err)
		return false
	}

	// convert setting bytes to the settings corresponding to the collector type
	pluginSettings, ok := readCollectorBytes(baseSettings.Type, settingsBytes)
	if !ok {
		return false
	}

	// validate settings
	if !pluginSettings.IsValid() {
		return false
	}

	return true
}

// check if the base collector settings are valid, returns true if the object is valid, false otherwise
func (base *CollectorSettingsBase) IsValid(collectorType discovery.CollectorName) bool {
	if !validateAutoInterval(base.AutoInterval) {
		logger.Info(
			"CollectorSettingsBase.IsValid %v AutoInterval should be between %v and %v but is %v\n",
			collectorType, minAutoIncrement, maxAutoIncrement, base.AutoInterval,
		)
		return false
	}
	if base.Type != collectorType {
		logger.Info("CollectorSettingsBase.IsValid collector type mismatch '%v' vs '%v'\n", collectorType, base.Type)
		return false
	}
	return true
}

// validate autoInterval value, returns true if the object is valid, false otherwise
func validateAutoInterval(autoIntervalMinutes uint) bool {
	return minAutoIncrement <= autoIntervalMinutes && autoIntervalMinutes <= maxAutoIncrement
}

// validate lldp collector settings, returns true if the object is valid, false otherwise
func (s *LldpSettings) IsValid() bool {
	return s.CollectorSettingsBase.IsValid(discovery.Lldp)
}

// validate neighbour collector settings, returns true if the object is valid, false otherwise
func (s *NeighbourSettings) IsValid() bool {
	return s.CollectorSettingsBase.IsValid(discovery.Arp)
}

// validate nmap collector settings, returns true if the object is valid, false otherwise
func (s *NmapSettings) IsValid() bool {
	return s.CollectorSettingsBase.IsValid(discovery.Nmap)
}

package discovery_settings

import (
	"encoding/json"

	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
)

var (
	// map [collector type] function that returns a new settings object for said collector type
	collectorToSettingsMap = map[discovery.CollectorName]func() iCollectorSettings{
		discovery.Arp:       func() iCollectorSettings { return &NeighbourSettings{} },
		discovery.Neighbour: func() iCollectorSettings { return &NeighbourSettings{} },
		discovery.Lldp:      func() iCollectorSettings { return &LldpSettings{} },
		discovery.Nmap:      func() iCollectorSettings { return &NmapSettings{} },
	}
)

// if collectorName is valid, returns a new settings object for the collectorName and true, otherwise returns (nil, false)
func createSettingsForCollector(collectorName discovery.CollectorName) (
	settingsObj iCollectorSettings,
	exists bool,
) {
	createSettings, ok := collectorToSettingsMap[collectorName]
	if !ok {
		logger.Info("createSettingsForCollector received unknown settings type %v\n", collectorName)
		return nil, false
	}
	return createSettings(), true
}

// attempts to read a byte array and convert it to LldpSettings, returns true if the conversion was successful, false otherwise
func (s *LldpSettings) readBytes(bytes []byte) bool {
	if err := json.Unmarshal(bytes, s); err != nil {
		logger.Info("ValidateDiscoverySettings could not unmarshal LldpSettings with err %v\n", err)
		return false
	}
	return true
}

// attempts to read a byte array and convert it to LldpSettings, returns true if the conversion was successful, false otherwise
func (s *NeighbourSettings) readBytes(bytes []byte) bool {
	if err := json.Unmarshal(bytes, s); err != nil {
		logger.Info("ValidateDiscoverySettings could not unmarshal NeighbourSettings with err %v\n", err)
		return false
	}
	return true
}

// attempts to read a byte array and convert it to LldpSettings, returns true if the conversion was successful, false otherwise
func (s *NmapSettings) readBytes(bytes []byte) bool {
	if err := json.Unmarshal(bytes, s); err != nil {
		logger.Info("ValidateDiscoverySettings could not unmarshal NmapSettings with err %v\n", err)
		return false
	}
	return true
}

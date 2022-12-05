package discovery_settings

import (
	"encoding/json"

	"github.com/untangle/golang-shared/services/discovery"
	"github.com/untangle/golang-shared/services/logger"
)

var (
	// map [collector type] function that returns a new settings object for said collector type
	collectorToSettingsMap = map[discovery.CollectorName]func() iCollectorSettings{
		discovery.Neighbour: func() iCollectorSettings { return &NeighbourSettings{} },
		discovery.Lldp:      func() iCollectorSettings { return &LldpSettings{} },
		discovery.Nmap:      func() iCollectorSettings { return &NmapSettings{} },
	}
)

// if collectorName is valid, returns a new settings object for the collectorName and true, otherwise returns (nil, false)
func readCollectorBytes(collectorName discovery.CollectorName, bytes []byte) (iCollectorSettings, bool) {
	createSettings, ok := collectorToSettingsMap[collectorName]
	if !ok {
		logger.Info("readCollectorBytes received unknown settings type %v\n", collectorName)
		return nil, false
	}

	settings := createSettings()
	if err := json.Unmarshal(bytes, settings); err != nil {
		logger.Info("readCollectorBytes could not unmarshal %v with err %v\n", collectorName, err)
		return nil, false
	}

	return settings, true
}

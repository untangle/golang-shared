package policy

import (
	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

var logger = logService.GetLoggerInstance()

// Returns a double map of policy plugin settings. E.g. map["plugin"]map[policy]interface{} where
// plugin and policyare a strings. This will allow for easy access to policy settings for a plugin.
func getAllPolicyConfigurationSettings() (map[string]map[string]interface{}, error) {

	f := settings.GetSettingsFileSingleton()

	policySettings := &PolicySettingsType{}

	if err := f.UnmarshalSettingsAtPath(&policySettings, "policy_manager"); err != nil {
		logger.Info("getAllPolicyConfigurationSettings failed : %v\n", err)
		return nil, err
	}

	// Process into a map of maps
	pluginSettings := make(map[string]map[string]interface{})

	// Go through each Policy and find matching configurations.
	for _, p := range policySettings.Policies {
		if !p.Enabled {
			continue
		}
		for _, config := range p.Configurations {
			config := policySettings.findConfiguration(*config)
			if config == nil {
				// No matching configuration found, skip. Although this should never happen.
				continue
			}
			// Add the plugins into the map. Wish there was a better way to do this
			if config.TPSettings != nil {
				pluginSettings[p.Name]["threatprevention"] = config.TPSettings
			}
			if config.WFSettings != nil {
				pluginSettings[p.Name]["webfilter"] = config.WFSettings
			}
			if config.GEOSettings != nil {
				pluginSettings[p.Name]["geoip"] = config.GEOSettings
			}
			if config.AppControlSettings != nil {
				pluginSettings[p.Name]["application_control"] = config.AppControlSettings
			}
		}
	}

	return nil, nil

}

// Returns a map of policy plugin settings for a given plugin. E.g. map[policy]interface{} where policy is
// the policy name and interface{} is the plugin settings.
func GetPolicyPluginSettings(pluginName string) map[string]interface{} {

	var pluginSettings map[string]map[string]interface{}
	var err error

	if pluginSettings, err = getAllPolicyConfigurationSettings(); err != nil {
		logger.Info("GetPolicyPluginSettings: %v\n", err)
		return nil
	}

	if pluginSettings[pluginName] == nil {
		logger.Info("GetPolicyPluginSettings: %v\n", "Plugin not found")
		return nil
	}

	return pluginSettings[pluginName]
}

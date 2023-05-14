package policy_settings

import (
	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

var logger = logService.GetLoggerInstance()

// Returns a map of policy plugin settings for a given plugin. E.g. map[policy]interface{} where policy is
// the policy name and interface{} is the plugin settings.
func GetPolicyPluginSettings(settingsFile *settings.SettingsFile, pluginName string) map[string]interface{} {

	var pluginSettings map[string]map[string]interface{}
	var err error

	if pluginSettings, err = getAllPolicyConfigurationSettings(settingsFile); err != nil {
		return nil
	}

	if pluginSettings[pluginName] == nil {
		return nil
	}

	return pluginSettings[pluginName]
}

// Returns a double map of policy plugin settings. E.g. map["plugin"]map[policy]interface{} where
// plugin and policyare a strings. This will allow for easy access to policy settings for a plugin.
// Each plugin is still responsible for adding the default entry.
func getAllPolicyConfigurationSettings(settingsFile *settings.SettingsFile) (map[string]map[string]interface{}, error) {

	policySettings := &PolicySettingsType{}

	if err := settingsFile.UnmarshalSettingsAtPath(&policySettings, "policy_manager"); err != nil {
		return nil, err
	}

	// Process into a map of maps
	pluginSettings := make(map[string]map[string]interface{})
	pluginSettings["threatprevention"] = make(map[string]interface{})
	pluginSettings["webfilter"] = make(map[string]interface{})
	pluginSettings["geoip"] = make(map[string]interface{})
	pluginSettings["application_control"] = make(map[string]interface{})

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
			logger.Info("getAllPolicyConfigurationSettings: %v, %+v\n", p.Name, config)
			if config.TPSettings != nil {
				pluginSettings["threatprevention"][p.Name] = config.TPSettings
			}
			if config.WFSettings != nil {
				pluginSettings["webfilter"][p.Name] = config.WFSettings
			}
			if config.GEOSettings != nil {
				pluginSettings["geoip"][p.Name] = config.GEOSettings
			}
			if config.AppControlSettings != nil {
				pluginSettings["application_control"][p.Name] = config.AppControlSettings
			}
		}
	}

	return pluginSettings, nil

}

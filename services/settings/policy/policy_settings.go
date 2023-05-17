package policy

import (
	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	// Defines the name of the settings properties for policy manager
	PolicyConfigName = "policy_manager"
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

	policySettings := &PolicySettings{}

	if err := settingsFile.UnmarshalSettingsAtPath(&policySettings, PolicyConfigName); err != nil {
		return nil, err
	}

	// Populate the Configurations.
	for _, p := range policySettings.TempConfigurations.([]interface{}) {
		data := p.(map[string]interface{})

		newConfig := &PolicyConfiguration{}
		newConfig.AppSettings = make(map[string]interface{})
		for pName, pValue := range data {
			switch pName {
			case "description":
				if _, ok := pValue.(string); ok {
					newConfig.Description = pValue.(string)
				}
			case "name":
				if _, ok := pValue.(string); ok {
					newConfig.Name = pValue.(string)
				}
			case "id":
				if _, ok := pValue.(string); ok {
					newConfig.ID = pValue.(string)
				}
			default: // Everything else is an app setting
				newConfig.AppSettings[pName] = pValue
			}
		}
		policySettings.Configurations = append(policySettings.Configurations, newConfig)
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
			logger.Debug("getAllPolicyConfigurationSettings: %v, %+v\n", p.Name, config)

			for name, settings := range config.AppSettings {
				if pluginSettings[name] == nil {
					pluginSettings[name] = make(map[string]interface{})
				}
				pluginSettings[name][p.Name] = settings
			}
		}
	}

	return pluginSettings, nil

}

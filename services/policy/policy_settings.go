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

	logger.Info("Found policy manager settings %+v\n", policySettings)
	logger.Info("There are %i polices\n", len(policySettings.Policies))
	logger.Info("There are %i flows\n", len(policySettings.Flows))
	logger.Info("There are %i configurations\n", len(policySettings.Configurations))
	logger.Info("Name of policy: %v\n", policySettings.Policies[0].Name)

	for _, policy := range policySettings.Policies {
		logger.Info("Parsing policy: %v\n", policy.Name)
		for _, config := range policy.Configuration {
			logger.Info("Plugin: %+v", config.PluginSettings)
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

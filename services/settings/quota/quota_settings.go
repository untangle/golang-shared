package quota

import (
	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	// Defines the name of the settings properties for quota manager
	QuotaConfigName   = "quota_manager"
	DefaultSettingUUID = "00000000-0000-0000-0000-000000000000"
)

var logger = logService.GetLoggerInstance()

// Returns a map of quota plugin settings for a given plugin. E.g. map[quota]interface{} where quota is
// the quota name and interface{} is the plugin settings.
func GetQuotaPluginSettings(settingsFile *settings.SettingsFile, pluginName string) (map[string]interface{}, error) {

	var pluginSettings map[string]map[string]interface{}
	var defaultPluginSettings interface{}
	var err error

	if pluginSettings, err = getAllQuotaConfigurationSettings(settingsFile); err != nil {
		return nil, err
	}

	// Add default settings into map with key default.
	if err := settingsFile.UnmarshalSettingsAtPath(&defaultPluginSettings, pluginName); err != nil {
		return nil, err
	}

	if _, ok := pluginSettings[pluginName]; !ok {
		pluginSettings[pluginName] = map[string]any{}
	}
	pluginSettings[pluginName][DefaultSettingUUID] = defaultPluginSettings
	return pluginSettings[pluginName], nil
}

// Returns a double map of quota plugin settings. E.g. map["plugin"]map[quota]interface{} where
// plugin and quotaare a strings. This will allow for easy access to quota settings for a plugin.
// Each plugin is still responsible for adding the default entry.
func getAllQuotaConfigurationSettings(settingsFile *settings.SettingsFile) (map[string]map[string]interface{}, error) {

	quotaSettings := &QuotaSettings{}

	if err := settingsFile.UnmarshalSettingsAtPath(&quotaSettings, QuotaConfigName); err != nil {
		return nil, err
	}

	// Process into a map of maps
	pluginSettings := make(map[string]map[string]interface{})

	// Go through each Quota and find matching configurations.
	for _, p := range quotaSettings.Policies {
		if !p.Enabled {
			continue
		}
		for _, config := range p.Configurations {
			config := quotaSettings.findConfiguration(config)
			if config == nil {
				logger.Warn("Can't find configuration in settings: %s(%s)\n",
					config.ID,
					config.Name)
				// No matching configuration found, skip. Although this should never happen.
				continue
			}
			// Add the plugins into the map. Wish there was a better way to do this
			logger.Debug("getAllQuotaConfigurationSettings: %v, %+v\n", p.Name, config)

			for name, settings := range config.AppSettings {
				if pluginSettings[name] == nil {
					pluginSettings[name] = make(map[string]interface{})
				}
				pluginSettings[name][p.ID] = settings
			}
		}
	}
	return pluginSettings, nil
}

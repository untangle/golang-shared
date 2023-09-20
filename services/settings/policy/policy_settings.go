package policy

import (
	logService "github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	// Defines the name of the settings properties for policy manager
	PolicyConfigName   = "policy_manager"
	DefaultSettingUUID = "00000000-0000-0000-0000-000000000000"
)

var logger = logService.GetLoggerInstance()

// PolicySettings is the main data structure for Policy Management.
// It contains an array of PolicyConfigurations, an array of PolicyFlowCategory's
// and an array of Policy which reference the Configurations and FlowCategories by id.
// Those arrays are loaded from the json primarily by mapstructure.
// facilitate lookup.
type PolicySettings struct {
	Enabled         bool                   `json:"enabled"`
	Configurations  []*PolicyConfiguration `json:"configurations"`
	Objects         []*Group               `json:"objects"`
	ObjectGroups    []*Object              `json:"object_groups"`
	Conditions      []*Object              `json:"conditions"`
	ConditionGroups []*Object              `json:"condition_groups"`
	Rules           []*Object              `json:"rules"`
	Policies        []*Policy              `json:"policies"`

	//DEPRECATED
	Flows  []*PolicyFlow `json:"flows,omitempty" `
	Groups []*Group      `json:"groups,omitempty"`
}

func (p *PolicySettings) findConfiguration(c string) *PolicyConfiguration {
	for _, config := range p.Configurations {
		if config.ID == c {
			return config
		}
	}
	return nil
}

// Returns the policy flow given the ID.
func (p *PolicySettings) FindFlow(id string) *Object {
	for _, flow := range p.Flows {
		if flow.ID == id {
			return flow
		}
	}
	return nil
}

// FindConfigsWithEnabled returns the configs with enabled status.
func (p *PolicySettings) FindConfigsWithEnabled(pol *Policy, enabled bool) []string {
	configs := []string{}
	for _, configID := range pol.Configurations {
		config := p.findConfiguration(configID)
		if config != nil && config.AppSettings != nil {
			for pluginName, pluginSettings := range config.AppSettings {
				if pluginSettings.(map[string]interface{})["enabled"] == enabled {
					configs = append(configs, pluginName)
				}
			}
		}
	}
	return configs
}

// Returns a list of disabled app services for a given policy ID.
func (p *PolicySettings) FindDisabledConfigs(pol *Policy) []string {
	return p.FindConfigsWithEnabled(pol, false)
}

// FindEnabledConfigs returns enabled configs for this policy
func (p *PolicySettings) FindEnabledConfigs(pol *Policy) []string {
	return p.FindConfigsWithEnabled(pol, true)
}

// Returns a map of policy plugin settings for a given plugin. E.g. map[policy]interface{} where policy is
// the policy name and interface{} is the plugin settings.
func GetPolicyPluginSettings(settingsFile *settings.SettingsFile, pluginName string) (map[string]interface{}, error) {

	var pluginSettings map[string]map[string]interface{}
	var defaultPluginSettings interface{}
	var err error

	if pluginSettings, err = getAllPolicyConfigurationSettings(settingsFile); err != nil {
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

// Returns a double map of policy plugin settings. E.g. map["plugin"]map[policy]interface{} where
// plugin and policyare a strings. This will allow for easy access to policy settings for a plugin.
// Each plugin is still responsible for adding the default entry.
func getAllPolicyConfigurationSettings(settingsFile *settings.SettingsFile) (map[string]map[string]interface{}, error) {

	policySettings := &PolicySettings{}

	if err := settingsFile.UnmarshalSettingsAtPath(&policySettings, PolicyConfigName); err != nil {
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
			config := policySettings.findConfiguration(config)
			if config == nil {
				logger.Warn("Can't find configuration in settings: %s(%s)\n",
					config.ID,
					config.Name)
				// No matching configuration found, skip. Although this should never happen.
				continue
			}
			// Add the plugins into the map. Wish there was a better way to do this
			logger.Debug("getAllPolicyConfigurationSettings: %v, %+v\n", p.Name, config)

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

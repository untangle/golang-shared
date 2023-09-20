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

// FindConfiguration searches this PolicySetting to load a configuration by ID
func (p *PolicySettings) FindConfiguration(configID string) *PolicyConfiguration {
	for _, config := range p.Configurations {
		if config.ID == configID {
			return config
		}
	}
	return nil
}

// FindRule searches this PolicySetting to load a rule by ID
func (p *PolicySettings) FindRule(ruleID string) *Object {
	for _, rule := range p.Rules {
		if rule.ID == ruleID {
			return rule
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

// FindConfigsWithEnabled returns the names of configs with enabled status that are attached to this policy
func (p *PolicySettings) FindConfigsWithEnabled(pol *Policy, enabled bool) []string {
	configNames := []string{}
	for _, ruleID := range pol.Rules {
		rule := p.FindRule(ruleID)
		if rule != nil {
			// We also have to check if this config is enabled/disabled...
			configDetails := p.FindConfiguration(rule.Action.UUID)
			if configDetails.Enabled == enabled {
				configNames = append(configNames, rule.Action.Key)
			}
		}
	}
	return configNames
}

// Returns a list of disabled app services for a given policy ID.
func (p *PolicySettings) FindDisabledConfigs(pol *Policy) []string {
	return p.FindConfigsWithEnabled(pol, false)
}

// FindEnabledConfigs returns enabled configs for this policy
func (p *PolicySettings) FindEnabledConfigs(pol *Policy) []string {
	return p.FindConfigsWithEnabled(pol, true)
}

// GetPolicyPluginSettings Returns a map of policy plugin settings for a given plugin.
// E.g. map[policy]interface{} where policy is
// the policy name and interface{} is the plugin settings.
// This returns default settings as well
func GetPolicyPluginSettings(settingsFile *settings.SettingsFile, pluginName string) (map[string]interface{}, error) {

	var pluginSettings map[string]map[string]interface{}
	var defaultPluginSettings interface{}
	var err error

	if pluginSettings, err = GetAllPolicyConfigs(settingsFile); err != nil {
		return nil, err
	}

	// Add default settings into map with key default.
	// This needs plugin metadata to figure out that 'mfw-template-XXX' is the same as the top level settings name
	if err := settingsFile.UnmarshalSettingsAtPath(&defaultPluginSettings, SettingsMetaLookup[pluginName].SettingsName); err != nil {
		return nil, err
	}

	if _, ok := pluginSettings[string(SettingsMetaLookup[pluginName].Type)]; !ok {
		pluginSettings[string(SettingsMetaLookup[pluginName].Type)] = map[string]any{}
	}
	pluginSettings[string(SettingsMetaLookup[pluginName].Type)][DefaultSettingUUID] = defaultPluginSettings
	return pluginSettings[string(SettingsMetaLookup[pluginName].Type)], nil
}

// GetAllPolicyConfigs Returns a double map of policy plugin settings. E.g. map["plugin"]map[policy]interface{} where
// plugin and policyare a strings. This will allow for easy access to policy settings for a plugin.
// Each plugin is still responsible for adding the default entry.
func GetAllPolicyConfigs(settingsFile *settings.SettingsFile) (map[string]map[string]interface{}, error) {

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

		for _, ruleID := range p.Rules {
			ruleDetails := policySettings.FindRule(ruleID)
			if ruleDetails != nil {
				// if action type is not SET_CONFIGURATION then do what?
				if ruleDetails.Action.Type == "SET_CONFIGURATION" {
					config := policySettings.FindConfiguration(ruleDetails.Action.UUID)

					if config == nil {
						logger.Warn("Can't find configuration in settings: %s(%s)\n",
							ruleDetails.Action.Key,
							ruleDetails.Action.UUID)
						// No matching configuration found, skip. Although this should never happen.
						continue
					}

					logger.Debug("Attaching action details of %v", ruleDetails.Action)

					if pluginSettings[string(config.Type)] == nil {
						pluginSettings[string(config.Type)] = make(map[string]interface{})
					}
					pluginSettings[string(config.Type)][config.ID] = config
				}
			}
		}
	}
	return pluginSettings, nil
}

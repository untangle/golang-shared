package policy

import (
	"github.com/untangle/golang-shared/services/settings"
	utilNet "github.com/untangle/golang-shared/util/net"
)

const (
	// Defines the name of the settings properties for policy manager
	PolicyConfigName   = "policy_manager"
	DefaultSettingUUID = "00000000-0000-0000-0000-000000000000"
)

// PolicySettings is the main data structure for Policy Management.
// It contains an array of PolicyConfigurations, an array of PolicyFlowCategory's
// and an array of Policy which reference the Configurations and FlowCategories by id.
// Those arrays are loaded from the json primarily by mapstructure.
// facilitate lookup.
type PolicySettings struct {
	Enabled         bool                   `json:"enabled"`
	Configurations  []*PolicyConfiguration `json:"configurations"`
	Objects         []*Object              `json:"objects"`
	ObjectGroups    []*Object              `json:"object_groups"`
	Conditions      []*Object              `json:"conditions"`
	ConditionGroups []*Object              `json:"condition_groups"`
	Rules           []*Object              `json:"rules"`
	Quotas          []*Object              `json:"quotas"`
	Policies        []*Policy              `json:"policies"`
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
	// This needs plugin metadata to figure out that 'mfw-config-XXX' is the same as the top level settings name
	if err := settingsFile.UnmarshalSettingsAtPath(&defaultPluginSettings, SettingsMetaLookup[pluginName].SettingsName); err != nil {
		return nil, err
	}

	if _, ok := pluginSettings[string(SettingsMetaLookup[pluginName].Type)]; !ok {
		pluginSettings[string(SettingsMetaLookup[pluginName].Type)] = map[string]any{}
	}
	pluginSettings[string(SettingsMetaLookup[pluginName].Type)][DefaultSettingUUID] = &Object{ID: DefaultSettingUUID, Settings: defaultPluginSettings, Type: SettingsMetaLookup[pluginName].Type}
	return pluginSettings[string(SettingsMetaLookup[pluginName].Type)], nil
}

// GetAllPolicyConfigs Returns a double map of policy plugin settings. E.g. map["plugin"]map[policy]interface{} where
// plugin and policy are a strings. This will allow for easy access to policy settings for a plugin.
// Each plugin is still responsible for adding the default entry.
func GetAllPolicyConfigs(settingsFile *settings.SettingsFile) (map[string]map[string]interface{}, error) {

	policySettings := &PolicySettings{}

	if err := settingsFile.UnmarshalSettingsAtPath(&policySettings, PolicyConfigName); err != nil {
		return nil, err
	}

	// Process into a map of maps
	pluginSettings := make(map[string]map[string]interface{})

	// Just pull policy configs from the configurations elements
	for _, config := range policySettings.Configurations {
		if pluginSettings[string(config.Type)] == nil {
			pluginSettings[string(config.Type)] = make(map[string]interface{})
		}
		pluginSettings[string(config.Type)][config.ID] = config
	}

	return pluginSettings, nil
}

// ItemsStringList returns the Items of the object as a slice of
// strings if they can be interpreted this way, or an empty slice and
// false if not.
func (o *Object) ItemsStringList() ([]string, bool) {
	val, ok := o.Items.([]string)
	return val, ok
}

// ItemsIPSpecList returns the Items of an object as a slice of
// utilNet.IPSpecifierString and true if they can be interpreted this way,
// or an empty slice and false otherwise.
func (o *Object) ItemsIPSpecList() ([]utilNet.IPSpecifierString, bool) {
	val, ok := o.Items.([]utilNet.IPSpecifierString)
	return val, ok
}

// ItemsServiceEndpointList returns the Items of an object as a slice of
// ServiceEndpoint and true if they can be interpreted this way, nil
// and false otherwise.
func (o *Object) ItemsServiceEndpointList() ([]ServiceEndpoint, bool) {
	val, ok := o.Items.([]ServiceEndpoint)
	return val, ok
}

// ItemsApplicationObject returns the Items of an object as a
// ApplicatonObject and true if they can be interpreted this way, nil
// and false otherwise.
func (o *Object) ItemsApplicationObject() (ApplicationObject, bool) {
	if val, ok := o.Items.([]ApplicationObject); ok {
		if len(val) > 0 && (len(val[0].Port) > 0 || len(val[0].IPAddrList) > 0) {
			return val[0], true
		}
	}
	// Returning an empty object prevents the objects loading from failing
	return ApplicationObject{}, false
}

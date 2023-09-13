package policy

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	utilNet "github.com/untangle/golang-shared/util/net"
)

// PolicySettings is the main data structure for Policy Management.
// It contains an array of PolicyConfigurations, an array of PolicyFlowCategory's
// and an array of Policy which reference the Configurations and FlowCategories by id.
// Those arrays are loaded from the json primarily by mapstructure.
// facilitate lookup.
type PolicySettings struct {
	Enabled        bool                   `json:"enabled"`
	Flows          []*PolicyFlow          `json:"flows"`
	Configurations []*PolicyConfiguration `json:"configurations"`
	Policies       []*Policy              `json:"policies"`
	Groups         []*Group               `json:"groups"`
}

// Object is a way to generically re-use the idea of something that is
// identified by ID, with associated metadata of name and description,
// with possible accompanying Items.
type Object struct {
	Name        string     `json:"name"`
	Type        ObjectType `json:"type"`
	Description string     `json:"description"`
	ID          string     `json:"id"`
	Enabled     bool       `json:"enabled,omitempty"`
	Items       any        `json:"items,omitempty"`

	// Other Object Types that use conditions
	Conditions []*PolicyCondition `json:"conditions,omitempty"`

	// Policy Object
	Rules []string `json:"rules,omitempty"`

	// DEPRECATED
	Configurations []string `json:"configurations,omitempty"`
	Flows          []string `json:"flows,omitempty"`
}

// Group is a deprecated concept, please use Object.
// Deprecated: Group is deprecated, use Object instead. See MFW-3517.
type Group = Object

// Policies are the root of our policy configurations. It includes pointers to substructure.
type Policy = Object

// ServiceEndpoint is a particular group type, a group may be
// identified by a list of these.
type ServiceEndpoint struct {
	Protocol uint `json:"protocol"`
	Port     uint `json:"port"`
}

// utility function for setting a list in the Group.Items field. We
// use a trick where json.Unmarshal will look at an 'any' value and if
// it has a pointer to a specific type, unmarshall into that
// type. However, we don't want the pointer later on, we just want the
// slice. setting g.Items to []T{} where T is a type we want does not
// work.
func setList[T any](g *Group) func() {
	list := []T{}
	g.Items = &list
	return func() {
		g.Items = list
	}
}

// UnmarshalJSON is a custom json unmarshaller for Objects.
func (obj *Object) UnmarshalJSON(data []byte) error {
	var typeField GroupTypeField

	if err := json.Unmarshal(data, &typeField); err != nil {
		return fmt.Errorf("unable to unmarshal group: %w", err)
	}

	switch typeField.Type {
	// If type field is empty - then we need to use a different type of alias to marshal (just direct object alias?)
	case "":
		type aliasObject Object

		if err := json.Unmarshal(data, (*aliasObject)(obj)); err != nil {
			return fmt.Errorf("unable to unmarshal generic object: %w", err)
		}
		return nil

	case IPAddrListType, IPObjectType:
		defer setList[utilNet.IPSpecifierString](obj)()
	case GeoIPListType, GeoIPObjectType:
		defer setList[string](obj)()
	case ServiceEndpointType, ServiceEndpointObjectType:
		defer setList[ServiceEndpoint](obj)()
	case InterfaceType, InterfaceObjectType:
		defer setList[uint](obj)()
	case ConditionType:
		defer setList[*PolicyCondition](obj)()
	case ConditionGroupType:
		defer setList[string](obj)()
	case ThreatPreventionType:
		defer setList[uint](obj)()
	case WebFilterCategoryType:
		defer setList[uint](obj)()
	default:
		return fmt.Errorf("error unmarshalling policy group: invalid group type: %s", typeField.Type)
	}

	// alias to make use of tags but avoid recursion
	type aliasGroup Group

	// unmarshal PolicyConfiguration using struct tags
	return json.Unmarshal(data, (*aliasGroup)(obj))
}

// ItemsStringList returns the Items of the group as a slice of
// strings if they can be interpreted this way, or an empty slice and
// false if not.
func (g *Group) ItemsStringList() ([]string, bool) {
	val, ok := g.Items.([]string)
	return val, ok
}

// ItemsIPSpecList returns the Items of a group as a slice of
// utilNet.IPSpecifierString and true if they can be interpreted this way,
// or an empty slice and false otherwise.
func (g *Group) ItemsIPSpecList() ([]utilNet.IPSpecifierString, bool) {
	val, ok := g.Items.([]utilNet.IPSpecifierString)
	return val, ok
}

// ItemsServiceEndpointList returns the Items of a group as a slice of
// ServiceEndpoint and true if they can be interpreted this way, nil
// and false otherwise.
func (g *Group) ItemsServiceEndpointList() ([]ServiceEndpoint, bool) {
	val, ok := g.Items.([]ServiceEndpoint)
	return val, ok
}

// PolicyFlow contains policy flow configuration.
type PolicyFlow struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Conditions  []*PolicyCondition `json:"conditions"`
}

// PolicyCondition contains policy condition configuration.
type PolicyCondition struct {
	Op      string   `json:"op"`
	CType   string   `json:"type"`
	Value   []string `json:"value,omitempty"`
	GroupID string   `json:"groupId"`
}

// Unmarshal policy condition so that types of values can be checked
func (pCondition *PolicyCondition) UnmarshalJSON(data []byte) error {
	// unmarshal like normal first
	type aliasPolicyCondition PolicyCondition
	if err := json.Unmarshal(data, (*aliasPolicyCondition)(pCondition)); err != nil {
		return err
	}

	// check that pCondition.Value is formatted correctly for the CType
	for i, value := range pCondition.Value {
		switch pCondition.CType {
		case "CLIENT_ADDRESS", "SERVER_ADDRESS":
			// Check that address is in CIDR format (w/ mask)
			if _, _, err := net.ParseCIDR(value); err != nil {
				// If address is a valid IP, but without a mask, just add the default
				if ip := net.ParseIP(value); ip != nil {
					if ip.To4() != nil {
						pCondition.Value[i] = fmt.Sprintf("%s%s", value, "/32")
					} else {
						pCondition.Value[i] = fmt.Sprintf("%s%s", value, "/64")
					}
				} else {
					return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
				}
			}
		case "CLIENT_PORT", "SERVER_PORT":
			if _, err := strconv.ParseUint(value, 10, 32); err != nil {
				return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
			}
		case "DAY_OF_WEEK", "DEST_ADDRESS", "GEOIP_LOCATION", "INTERFACE", "SERVICE_ENDPOINT", "SOURCE_ADDRESS", "PROTOCOL_TYPE", "TIME_OF_DAY", "VLAN_ID", "THREATPREVENTION":
			// These are not yet implemented and need to have a designated format
		default:
			return fmt.Errorf("error while unmarshalling policy condition: invalid type: %s", pCondition.CType)
		}
	}

	return nil
}

// PolicyConfiguration contains policy configuration.
type PolicyConfiguration struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AppSettings map[string]interface{} `json:"-"` // map of plugin settings, key is the plugin name.
}

func (pConfig PolicyConfiguration) MarshalJSON() ([]byte, error) {
	// alias to make use of tags but avoid recursion
	type aliasPolicyConfiguration PolicyConfiguration

	// marshal PolicyConfiguration using struct tags
	fieldJSON, err := json.Marshal((*aliasPolicyConfiguration)(&pConfig))
	if err != nil || len(pConfig.AppSettings) == 0 {
		// return if there was an error or nothing else needs to be marshalled
		return fieldJSON, err
	}

	// marshal AppSettings separately
	dynamicJSON, err := json.Marshal(pConfig.AppSettings)
	if err != nil {
		return nil, err
	}
	// replace opening '{' with ',' as fields from dynamicJSON will be added
	dynamicJSON[0] = ','
	return append(fieldJSON[:len(fieldJSON)-1], dynamicJSON...), nil

}

func (pConfig *PolicyConfiguration) UnmarshalJSON(data []byte) error {
	// alias to make use of tags but avoid recursion
	type aliasPolicyConfiguration PolicyConfiguration

	// unmarshal PolicyConfiguration using struct tags
	if err := json.Unmarshal(data, (*aliasPolicyConfiguration)(pConfig)); err != nil {
		return err
	}

	// unmarshal remaining fields in JSON and put them in AppSettings
	dataMap := make(map[string]any)
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return err
	}
	// delete fields that are not part of AppSettings
	delete(dataMap, "description")
	delete(dataMap, "name")
	delete(dataMap, "id")

	pConfig.AppSettings = dataMap

	return nil
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
func (p *PolicySettings) FindFlow(id string) *PolicyFlow {
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

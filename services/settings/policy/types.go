package policy

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/untangle/golang-shared/util/net"
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

// GroupType is the type of group that a Group is, used to demux the
// Items field.
type GroupType string

const (
	// GeoIPListType means that the Items of a Group are geoip countries.
	GeoIPListType GroupType = "GeoIPLocation"

	// IPAddrListType means that the Items of the Group are ip
	// specifications (ranges, CIDRs, or single IPs).
	IPAddrListType GroupType = "IPAddrList"

	// ServiceEndpointType means that the Items of a Group are
	// service endpoints.
	ServiceEndpointType GroupType = "ServiceEndpoint"
)

// Group is a way to generically re-use certain lists of attributes
// that may be true for a session.
type Group struct {
	Name        string    `json:"name"`
	Type        GroupType `json:"type"`
	Description string    `json:"description"`
	ID          string    `json:"id"`
	Items       any       `json:"items"`
}

// ServiceEndpoint is a particular group type, a group may be
// identified by a list of these.
type ServiceEndpoint struct {
	Protocol    string `json:"protocol"`
	IPSpecifier string `json:"ipspecifier"`
	Port        uint   `json:"port"`
}

// UnmarshalJSON is a custom json unmarshaller for a Group.
func (g *Group) UnmarshalJSON(data []byte) error {

	type GroupTypeField struct {
		Type GroupType `json:"type"`
	}
	var typeField GroupTypeField

	if err := json.Unmarshal(data, &typeField); err != nil {
		return fmt.Errorf("unable to unmarshal group: %w", err)
	}

	switch typeField.Type {
	case IPAddrListType:
		list := []net.IPSpecifierString{}
		g.Items = &list
		defer func() { g.Items = list }()
	case GeoIPListType:
		list := []string{}
		g.Items = &list
		defer func() { g.Items = list }()
	case ServiceEndpointType:
		list := []ServiceEndpoint{}
		g.Items = &list
		defer func() { g.Items = list }()
	default:
		return fmt.Errorf("error unmarshalling policy group: invalid group type: %s", typeField.Type)
	}

	// alias to make use of tags but avoid recursion
	type aliasGroup Group

	// unmarshal PolicyConfiguration using struct tags
	return json.Unmarshal(data, (*aliasGroup)(g))
}

// UnmarshalJSON required since "port" is a string in settings but a uint in the ServiceEndpoint
func (se *ServiceEndpoint) UnmarshalJSON(data []byte) error {
	var rawvalue map[string]interface{}

	if err := json.Unmarshal(data, &rawvalue); err != nil {
		return fmt.Errorf("unable to unmarshal endpoint: %w", err)
	}

	config := mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           se,
		WeaklyTypedInput: true,
	}

	d, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return fmt.Errorf("unable to decode service endpoint: %w", err)
	}
	if err = d.Decode(rawvalue); err != nil {
		return fmt.Errorf("error while trying to unmarshal policy group service endpoint values: %w",
			err)
	}
	return nil
}

// MarshalJSON for a ServiceEndpoint
func (se ServiceEndpoint) MarshalJSON() ([]byte, error) {
	type alias ServiceEndpoint

	// use struct to write port field as a string
	type fieldOverwrite struct {
		*alias
		Port string `json:"port"`
	}

	overwrite := &fieldOverwrite{
		alias: (*alias)(&se),
		Port:  fmt.Sprint(se.Port),
	}

	return json.Marshal(overwrite)
}

// ItemsStringList returns the Items of the group as a slice of
// strings if they can be interpreted this way, or an empty slice and
// false if not.
func (g *Group) ItemsStringList() ([]string, bool) {
	val, ok := g.Items.([]string)
	return val, ok
}

// ItemsIPSpecList returns the Items of a group as a slice of
// net.IPSpecifierString and true if they can be interpreted this way,
// or an empty slice and false otherwise.
func (g *Group) ItemsIPSpecList() ([]net.IPSpecifierString, bool) {
	val, ok := g.Items.([]net.IPSpecifierString)
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
	Value   []string `json:"value"`
	GroupID string   `json:"groupId"`
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

// Policies are the root of our policy configurations. It includes pointers to substructure.
type Policy struct {
	Defaults       bool     `json:"defaults"`
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Enabled        bool     `json:"enabled"`
	Configurations []string `json:"configurations"`
	Flows          []string `json:"flows"`
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

// Returns a list of disabled app services for a given policy ID.
func (p *PolicySettings) FindDisabledConfigs(pol *Policy) []string {

	disabledConfigs := []string{}
	for _, configID := range pol.Configurations {
		config := p.findConfiguration(configID)
		if config != nil && config.AppSettings != nil {
			for pluginName, pluginSettings := range config.AppSettings {
				if pluginSettings.(map[string]interface{})["enabled"] == false {
					disabledConfigs = append(disabledConfigs, pluginName)
				}
			}
		}
	}
	return disabledConfigs
}

package policy

import (
	"encoding/json"
	"fmt"

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

	// WebFilterCategoryType means that the Items of the Group are web filter categories.
	WebFilterCategoryType GroupType = "WebFilterCategory"

	// ClassifyType means that the Items of the Group are classifications.
	ClassifyType GroupType = "Classify" // TBD Not implemented.

	// ThreatPreventionType means that the Items of the Group are threat prevention score.
	ThreatPreventionType GroupType = "ThreatPrevention"
	// InterfaceType is a group type where all items are interface
	// IDs (integers)
	InterfaceType GroupType = "Interface"
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
		return g.parseIPSpecList(rawvalue)
	case GeoIPListType, WebFilterCategoryType:
		return g.parseStringList(rawvalue)
	case ServiceEndpointType:
		return g.parseServiceEndpointList(rawvalue)
	case ThreatPreventionType:
		return g.parseIntList(rawvalue)
	case ClassifyType: // TBD Not implemented.
		return nil
	default:
		return fmt.Errorf("error unmarshalling policy group: invalid group type: %s", theType)
	}
}

func parseIntList[T any](raw map[string]interface{}, conv func(int) T) ([]T, error) {
	items, ok := raw["items"]
	if !ok {
		return nil, fmt.Errorf("error unmarshalling policy group: no items")
	}
	values, ok := items.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"error unmarshalling policy group: not a string list: %v(type is: %T)",
			values, values)
	}
	intList := make([]T, len(values), len(values))
	for i, val := range values {
		intConversion, ok := val.(int)
		if !ok {
			return nil, fmt.Errorf(
				"not a valid string in list of strings: %v(type is: %T)",
				val, val)
		}
		intList[i] = conv(intConversion)
	}
	return intList, nil
}

func parseStringableList[T any](raw map[string]interface{}, conv func(string) T) ([]T, error) {
	items, ok := raw["items"]
	if !ok {
		return nil, fmt.Errorf("error unmarshalling policy group: no items")
	}
	values, ok := items.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"error unmarshalling policy group: not a string list: %v(type is: %T)",
			values, values)
	}
	stringList := make([]T, len(values), len(values))
	for i, val := range values {
		stringConversion, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf(
				"not a valid string in list of strings: %v(type is: %T)",
				val, val)
		}
		stringList[i] = conv(stringConversion)
	}
	return stringList, nil
}

func (g *Group) parseIntList(raw map[string]interface{}) error {
	if items, err := parseIntList(raw, func(x int) int { return x }); err != nil {
		return err
	} else {
		g.Items = items
	}
	return nil
}

func (g *Group) parseStringList(raw map[string]interface{}) error {
	if items, err := parseStringableList(raw, func(x string) string { return x }); err != nil {
		return err
	} else {
		g.Items = items
	}
	return nil
}

func (g *Group) parseIPSpecList(raw map[string]interface{}) error {
	if items, err := parseStringableList(raw,
		func(x string) net.IPSpecifierString { return net.IPSpecifierString(x) }); err != nil {
		return err
	} else {
		g.Items = items
	}
	return nil
}

func (g *Group) parseServiceEndpointList(raw map[string]interface{}) error {
	items, ok := raw["items"]
	if !ok {
		return fmt.Errorf("error unmarshalling policy group: no items")
	}
	values, ok := items.([]interface{})
	if !ok {
		return fmt.Errorf(
			"error unmarshalling policy group: not a list of objects: %v(%T)",
			values, items)
		defer setList[net.IPSpecifierString](g)()
	case GeoIPListType:
		defer setList[string](g)()
	case ServiceEndpointType:
		defer setList[ServiceEndpoint](g)()
	case InterfaceType:
		defer setList[uint](g)()
	default:
		return fmt.Errorf("error unmarshalling policy group: invalid group type: %s", typeField.Type)
	}

	// alias to make use of tags but avoid recursion
	type aliasGroup Group

	// unmarshal PolicyConfiguration using struct tags
	return json.Unmarshal(data, (*aliasGroup)(g))
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
	Value   []string `json:"value,omitempty"`
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

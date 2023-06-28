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
	Enabled            bool          `json:"enabled"`
	Flows              []*PolicyFlow `json:"flows"`
	TempConfigurations interface{}   `json:"configurations"` // Config is dynamic so need temp place to store it.
	Configurations     []*PolicyConfiguration
	Policies           []*Policy `json:"policies"`
	Groups             []*Group  `json:"groups"`
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
	Type        GroupType
	Description string
	ID          string
	Items       any
}

// ServiceEndpoint is a particular group type, a group may be
// identified by a list of these.
type ServiceEndpoint struct {
	Protocol uint `json:"protocol"`
	Port     uint `json:"port"`
}

// UnmarshalJSON is a custom json unmarshaller for a Group.
func (g *Group) UnmarshalJSON(data []byte) error {
	var rawvalue struct {
		Type        GroupType `json:"type"`
		Description string    `json:"description"`
		ID          string    `json:"id"`
	}

	if err := json.Unmarshal(data, &rawvalue); err != nil {
		return fmt.Errorf("unable to unmarshal group: %w", err)
	}

	g.ID = rawvalue.ID
	g.Description = rawvalue.Description
	g.Type = rawvalue.Type

	switch g.Type {
	case IPAddrListType:
		return g.parseIPSpecList(data)
	case GeoIPListType:
		return g.parseStringList(data)
	case ServiceEndpointType:
		return g.parseServiceEndpointList(data)
	default:
		return fmt.Errorf("error unmarshalling policy group: invalid group type: %s", g.Type)
	}
}

func parseList[T any](raw []byte) ([]T, error) {
	var value struct {
		Items []T `json:"items"`
	}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return value.Items, nil
}

func (g *Group) parseStringList(raw []byte) error {
	if items, err := parseList[string](raw); err != nil {
		return err
	} else {
		g.Items = items
	}
	return nil
}

func (g *Group) parseIPSpecList(raw []byte) error {
	if items, err := parseList[net.IPSpecifierString](raw); err != nil {
		return err
	} else {
		g.Items = items
	}
	return nil
}

func (g *Group) parseServiceEndpointList(raw []byte) error {
	if items, err := parseList[ServiceEndpoint](raw); err != nil {
		return err
	} else {
		g.Items = items
	}
	return nil
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
	ID          string
	Name        string
	Description string
	AppSettings map[string]interface{} // map of plugin settings, key is the plugin name.
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

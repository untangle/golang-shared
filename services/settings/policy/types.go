package policy

import (
	"net"
)

// PolicySettings is the main data structure for Policy Management.
// It contains an array of PolicyConfigurations, an array of PolicyFlowCategory's
// and an array of Policy which reference the Configurations and FlowCategories by id.
// Those arrays are loaded from the json primarily by mapstructure.
// PolicyManager also maintains map[string]'s based on those arrays to
// facilitate lookup.
type PolicySettings struct {
	Enabled            bool          `json:"enabled"`
	Flows              []*PolicyFlow `json:"flows"`
	TempConfigurations interface{}   `json:"configurations"` // Config is dynamic so need temp place to store it.
	Configurations     []*PolicyConfiguration
	Policies           []*Policy `json:"policies"`
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
	Op    string   `json:"op"`
	CType string   `json:"type"`
	Value []string `json:"value"`
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

// Equal compares the leftSide to the PolicyConditionType
// and returns true if the leftSide matches the condition.
func (c *PolicyCondition) Equals(leftSide *string) bool {
	switch c.CType {
	case "CLIENT_ADDRESS", "SERVER_ADDRESS":
		// convert to net.IPNet
		for _, v := range c.Value {
			_, cidr, _ := net.ParseCIDR(v)
			if cidr.Contains(net.ParseIP(*leftSide)) {
				return true
			}
		}
	case "CLIENT_PORT", "SERVER_PORT", "FAMILY":
		for _, v := range c.Value {
			if *leftSide == v {
				return true
			}
		}
	}
	return false
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

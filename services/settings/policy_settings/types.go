package policy_settings

// PolicySettingsType is the main data structure for Policy Management.
// It contains an array of PolicyConfigurations, an array of PolicyFlowCategory's
// and an array of Policy which reference the Configurations and FlowCategories by id.
// Those arrays are loaded from the json primarily by mapstructure.
// PolicyManager also maintains map[string]'s based on those arrays to
// facilitate lookup.
type PolicySettingsType struct {
	Enabled        bool                       `json:"enabled"`
	Flows          []*PolicyFlowType          `json:"flows"`
	Configurations []*PolicyConfigurationType `json:"configurations"`
	Policies       []*PolicyType              `json:"policies"`
}

type PolicyFlowType struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Conditions  []*PolicyConditionType `json:"conditions"`
}

type PolicyConditionType struct {
	Op    string   `json:"op"`
	CType string   `json:"type"`
	Value []string `json:"value"`
}

type PolicyConfigurationType struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	TPSettings         interface{} `json:"threatprevention",optional:"true"`
	WFSettings         interface{} `json:"webfilter",optional:"true"`
	GEOSettings        interface{} `json:"geoip",optional:"true"`
	AppControlSettings interface{} `json:"application_control",optional:"true"`
}

type PolicyType struct {
	Defaults       bool      `json:"defaults"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Enabled        bool      `json:"enabled"`
	Configurations []*string `json:"policyConfigurations"`
	Flows          []*string `json:"flowCategories"`
}

func (p *PolicyType) GetFlows() []*string {
	return p.Flows
}

func (p PolicyType) GetName() string {
	return p.Name
}

func (p *PolicyType) IsEnabled() bool {
	return p.Enabled
}

func (p *PolicyFlowType) GetConditions() []*PolicyConditionType {
	return p.Conditions
}

func (c *PolicyConditionType) GetType() string {
	return c.CType
}

func (c *PolicyConditionType) GetOp() string {
	return c.Op
}

func (c *PolicyConditionType) GetValue() string {
	return c.Value
}

func (p *PolicySettingsType) findConfiguration(c string) *PolicyConfigurationType {
	for _, config := range p.Configurations {
		if config.ID == c {
			return config
		}
	}
	return nil
}

// Returns the policy flow given the ID.
func (p *PolicySettingsType) FindFlow(id string) *PolicyFlowType {
	for _, flow := range p.Flows {
		if flow.ID == id {
			return flow
		}
	}
	return nil
}

var policyConditionTypeMap = map[string]bool{
	"CLIENT_ADDRESS": true,
	"CLIENT_PORT":    true,
	"DAY_OF_WEEK":    true,
	"DEST_ADDRESS":   true,
	"INTERFACE":      true,
	"SERVER_ADDRESS": true,
	"SERVER_PORT":    true,
	"SOURCE_ADDRESS": true,
	"PROTOCOL_TYPE":  true,
	"TIME_OF_DAY":    true,
	"VLAN_ID":        true,
}

// Valid PolicyCondition Ops - there may be more at some point
// == implies an OR operation between the different entries in the value arrray
// != implies an AND operation between the different entries in the value array
// all other operations assume a single entry in the value array (or string)
var policyConditionOpsMap = map[string]bool{
	"==": true,
	"!=": true,
	"<":  true,
	">":  true,
	"<=": true,
	">=": true,
}

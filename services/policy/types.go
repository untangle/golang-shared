package policy

// policyManager config.
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
	Op    string `json:"op"`
	CType string `json:"type"`
	Value string `json:"value"`
}

type PolicyConfigurationType struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	PluginSettings interface{} `json:"service"`
}

type PolicyType struct {
	Defaults      bool                       `json:"defaults"`
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Enabled       bool                       `json:"enabled"`
	Configuration []*PolicyConfigurationType `json:"policyConfigurations"`
	Flows         []*PolicyFlowType          `json:"flowCategories"`
}

func (p *PolicyType) GetFlows() []*PolicyFlowType {
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

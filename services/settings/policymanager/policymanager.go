package policy

import (
	"fmt"
	"sync"

	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

// PolicyManager is the main data structure for Policy Management.
// It contains an array of PolicyConfigurations, an array of PolicyFlowCategory's
// and an array of Policy which reference the Configurations and FlowCategories by id.
// Those arrays are loaded from the json primarily by mapstructure.
// PolicyManager also maintains map[string]'s based on those arrays to
// facilitate lookup.
type PolicyManager struct {
	// Fields populated using mapstructure
	Enabled bool `json:"enabled"`

	// These arrays are loaded directly by mapstructure Decode
	ConfigurationArray []*PolicyConfiguration `json:"configurations"`
	FlowArray          []*PolicyFlowCategory  `json:"flows"`
	PolicyArray        []*Policy              `json:"policies"`

	// These maps resolved after loading the arrays above
	configurations map[string]*PolicyConfiguration
	flowCategories map[string]*PolicyFlowCategory
	policies       map[string]*Policy

	policySettingsLock sync.RWMutex
	settingsFile       *settings.SettingsFile
	settings           map[string]interface{}
	logger             *logger.Logger
}

// The PolicyFlowCategory captures a set of PolicyConditions
// that determine whether the containing Policy applies to
// traffic as it is seen.
// Each PolicyFlowCategory is identifed by its unique ID.
// The data model does not support charing PolicyConditions
// between Policy's although PolicyConfigurations can be shared.
type PolicyFlowCategory struct {
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	ConditionArray []PolicyCondition `json:"conditions"`
}

// PolicyCondition contains the criteria to test packets against
// to detemine whether the associated PolicyFlowCategor and Policy apply.
// Valid CTypes are specified in policysettings.go policyConditionTypeMap.
// VAlid Ops are specified in policysettings.go policyConditionOpsMap.
type PolicyCondition struct {
	CType string   `json:"type"`
	Op    string   `json:"op"`
	Value []string `json:"value"`
}

// PolicyConfiguration configures which plugins are applied when
// the PolicyFlowCategory/PolicyConditions are met
// Typical plugins would be threatprevention, geoip or webfilter
type PolicyConfiguration struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	TPSettings  interface{} `json:"threatprevention",optional:"true"`
	WFSettings  interface{} `json:"webfilter",optional:"true"`
	GEOSettings interface{} `json:"geoip",optional:"true"`
	// This probably doesn't belong here but keep it here for now
	// to support some settings files that have it
	DiscoverySettings interface{} `json:"discovery",optional:"true"`
}

// This is the main Policy object which is contained in an array in PolicyManager
// Each policy contains a set of configuration id's and a set of flow id's
// which are resolved by look up in the arrays in PolicyManager.
type Policy struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Enabled        bool     `json:"enabled"`
	Configurations []string `json:"configurations"`
	Flows          []string `json:"flows"`
}

// Returns a new policy instance
func NewPolicyManager(
	settingsFile *settings.SettingsFile,
	logger *logger.Logger) *PolicyManager {
	return &PolicyManager{
		logger:       logger,
		settingsFile: settingsFile,
	}
}

func (p *PolicyManager) LoadPolicyManagerSettings() error {
	p.policySettingsLock.Lock()
	defer p.policySettingsLock.Unlock()

	if err := p.readPolicyManagerSettings(); err != nil {
		return err
	}

	// Now populate the maps in PolicyManager and p.policies
	// to facilitate lookup at runtime
	p.configurations = make(map[string]*PolicyConfiguration, len(p.ConfigurationArray))
	for _, config := range p.ConfigurationArray {
		p.configurations[config.Id] = config
	}
	p.flowCategories = make(map[string]*PolicyFlowCategory, len(p.FlowArray))
	for _, flow := range p.FlowArray {
		p.flowCategories[flow.Id] = flow
	}
	p.policies = make(map[string]*Policy, len(p.PolicyArray))
	for _, policy := range p.PolicyArray {
		p.policies[policy.Id] = policy
	}
	return nil
}

// Basic validation to make sure that constrained fields are valid
// and id fields can be looked up in the appropriate maps
// Not sure whether we need/want to validate id, name, description as
// not empty.
func (p *PolicyManager) ValidatePolicies() error {
	p.policySettingsLock.RLock()
	defer p.policySettingsLock.RUnlock()

	for _, policy := range p.policies {
		for _, configId := range policy.Configurations {
			if _, ok := p.configurations[configId]; !ok {
				return fmt.Errorf("validatePolicies: found invalid configuration Id: %s in Policy %s",
					configId, policy.Id)
			}
		}
		for _, flowId := range policy.Flows {
			if _, ok := p.flowCategories[flowId]; !ok {
				return fmt.Errorf("validatePolicies: found invalid flow Id: %s in Policy %s",
					flowId, policy.Id)
			}
			for _, cond := range p.flowCategories[flowId].ConditionArray {
				if _, ok := policyConditionTypeMap[cond.CType]; !ok {
					return fmt.Errorf("validatePolicies: found invalid CType: %s in Policy %s, Flow %s",
						cond.CType, policy.Id, flowId)
				}
				if _, ok := policyConditionOpsMap[cond.Op]; !ok {
					return fmt.Errorf("validatePolicies: found invalid Op: %s in Policy %s, Flow %s",
						cond.Op, policy.Id, flowId)
				}
			}
		}
	}
	return nil
}

func (p *PolicyManager) readPolicyManagerSettings() error {
	if err := p.settingsFile.UnmarshalSettingsAtPath(&p, "policy_manager"); err != nil {
		p.logger.Err("Could not read Policy Manager Settings from", p.settingsFile)
		return err
	}
	return nil
}

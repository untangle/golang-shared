package policy

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
)

const (
	pluginName         = "policy"
	licensePluginName  = "untangle-node-policy"
	defaultSensitivity = 20

	policyPriority = 1
)

type PolicyManager struct {
	// Fields populated using mapstructure
	Id                 string                `json:"id"`
	Enabled            bool                  `json:"enabled"`
	NameField          string                `json:"name"`
	Description        string                `json:"description"`
	ConfigurationArray []PolicyConfiguration `json:"configurations"`
	FlowArray          []PolicyFlowCategory  `json:"flows"`
	PolicyArray        []Policy              `json:"policies"`

	// Fields resolved after loading the arrays above
	configurations map[string]*PolicyConfiguration
	flowCategories map[string]*PolicyFlowCategory
	policies       map[string]*Policy

	policySettingsLock sync.RWMutex
	settingsFile       *settings.SettingsFile
	settings           map[string]interface{}
	logger             *logger.Logger
}

type PolicyFlowCategory struct {
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	ConditionArray []PolicyCondition `json:"conditions"`
}

type PolicyCondition struct {
	CType string `json:"type"`
	Op    string `json:"op"`
	Value string `json:"value"`
	// There was a discussion about allowing value to be an array.
	// For now, I am assuming that it is only a string
	// with a comma-separated list if need be.
	// Having it be either an []string or string is not trivial AFAICT with mapstructure
	// although it is easy if we parse the map[string]interface{}
	// as initially done.
	value []string
}

type PolicyConfiguration struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// This doesn't parse completely with mapstructure
	// so we need to resolve this after the initial load
	Other   map[string]interface{} `json:",remain"`
	plugins []*PolicyPluginCategory
}

// TODO WIP
// PolicyPlugin needs to be an interface which
// can be implemented by a "geoip", "threatprevention", or "webfilter" configuration
type PolicyPluginCategory struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Conditions  []*PolicyPlugin `json:"conditions"`
}

// TODO WIP
type PolicyPlugin interface {
}

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
	decoderConfig := mapstructure.DecoderConfig{
		TagName:     "json",
		Result:      p,
		ErrorUnused: true,
		Squash:      true,
	}
	mapstructdecoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		p.logger.Err("policymanager: could not cerate mapstructure decoder", err)
		return err
	}
	if err := mapstructdecoder.Decode(&p.settings); err != nil {
		p.logger.Warn("policymanager: could not decode json:", err)
	}
	// Now populate the maps in PolicyManager and p.policies
	// to facilitate lookup at runtime
	p.configurations = make(map[string]*PolicyConfiguration, len(p.ConfigurationArray))
	for _, config := range p.ConfigurationArray {
		p.configurations[config.Id] = &config
	}
	p.flowCategories = make(map[string]*PolicyFlowCategory, len(p.FlowArray))
	for _, flow := range p.FlowArray {
		p.flowCategories[flow.Id] = &flow
	}
	p.policies = make(map[string]*Policy, len(p.PolicyArray))
	for _, policy := range p.PolicyArray {
		p.policies[policy.Id] = &policy
	}
	return p.validatePolicies()
}

func (p *PolicyManager) validatePolicies() error {
	for _, policy := range p.policies {
		for _, configId := range policy.Configurations {
			if _, ok := p.configurations[configId]; !ok {
				return fmt.Errorf("validatePolicies: found invalid configuration Id: %s", configId)
			}
		}
		for _, flowId := range policy.Flows {
			if _, ok := p.flowCategories[flowId]; !ok {
				return fmt.Errorf("validatePolicies: found invalid flow Id: %s", flowId)
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
	if err := p.settingsFile.UnmarshalSettingsAtPath(&p.settings, "policy_manager"); err != nil {
		p.logger.Err("Could not read Policy Manager Settings from", p.settingsFile)
		return err
	}
	return nil
}

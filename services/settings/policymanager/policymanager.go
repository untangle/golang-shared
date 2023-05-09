package policy

import (
	"sync"
	"syscall"

	"github.com/mitchellh/mapstructure"
	"github.com/untangle/golang-shared/services/alerts"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/untangle/golang-shared/services/settings"
	"github.com/untangle/golang-shared/util/net/interfaces"
)

const (
	pluginName         = "policy"
	licensePluginName  = "untangle-node-policy"
	defaultSensitivity = 20

	policyPriority = 1
)

type PolicyManager struct {
	Id                 string                `json:"id"`
	Enabled            bool                  `json:"enabled"`
	NameField          string                `json:"name"`
	Description        string                `json:"description"`
	ConfigurationArray []PolicyConfiguration `json:"configurations"`
	FlowArray          []PolicyFlowCategory  `json:"flows"`
	PolicyArray        []Policy              `json:"policies"`
	configurations     map[string]*PolicyConfiguration
	flowCategories     map[string]*PolicyFlowCategory
	conditions         map[string]*PolicyCondition
	policies           map[string]*Policy
	policySettingsLock sync.RWMutex
	settingsFile       *settings.SettingsFile
	settings           map[string]interface{}
	interfaceSettings  *interfaces.InterfaceSettings
	alertsPublisher    alerts.AlertPublisher
	logger             logger.Logger
}

type PolicyFlowCategory struct {
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:description"`
	ConditionArray []PolicyCondition `json:"conditions"`
	conditions     []*PolicyCondition
}

type PolicyCondition struct {
	CType string   `json:"type"`
	Op    string   `json:"op"`
	Value []string `json:"value"`
}

type PolicyConfiguration struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	plugins     []*PolicyPluginCategory
}

// PolicyPlugin needs to be an interface which
// can be implemented by a "geoip", "threatprevention", or "webfilter" configuration
type PolicyPluginCategory struct {
	id          string          `json:"id"`
	name        string          `json:"name"`
	description string          `json:"description"`
	conditions  []*PolicyPlugin `json:"conditions"`
}

type PolicyPlugin interface {
}

type Policy struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Enabled        bool     `json:"enabled"`
	Configurations []string `json:"configurations"`
	Flows          []string `json:"flows"`
	configurations map[string]*PolicyConfiguration
	flows          map[string]*PolicyFlowCategory
}

// Returns a new policy instance
func NewPolicyManager(
	settingsFile *settings.SettingsFile,
	logger logger.Logger) *PolicyManager {
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
	if p.settings["enabled"] != nil {
		p.Enabled = p.settings["enabled"].(bool)
	}
	p.conditions = make(map[string]*PolicyCondition)
	// Load all PolicyFlowCategories
	// and create a map based on id// findPolicy returns matching polices.
	// func (p *PolicyManager) findPolicy(s net.IP) *string {
	// 	for _, pol := range p.policySettings {
	// 		if p.matchPolicy(s, pol) {
	// 			return &pol.Name
	// 		}
	// 	}
	// 	return nil
	// }

	// func (p *PolicyManager) matchPolicy(sourceAddr net.IP, pol *policySettingsType) bool {
	// 	// Check if source matches.

	// 	for _, pSrc := range pol.Source {
	// 		if _, psource, err := net.ParseCIDR(pSrc); err == nil {
	// 			if psource.Contains(sourceAddr) {
	// 				p.logger.Debug("Policy match: %v", pol.Name)
	// 				return true
	// 			}
	// 		} else {
	// 			p.logger.Err("Error parsing policy source: %v\n", err)
	// 		}
	// 	}
	// 	return false
	// }

	p.flowCategories = make(map[string]*PolicyFlowCategory)
	if p.settings["flows"] != nil {
		// and also populate the conditions based on configured id's
		for _, v := range p.settings["flows"].([]interface{}) {
			// Populate
			fc, err := p.NewPolicyFlowCategory(v)
			if err != nil {
				return err
			}
			p.flowCategories[fc.Id] = fc
		}
	}
	p.configurations = make(map[string]*PolicyConfiguration)
	if p.settings["configurations"] != nil {
		// Load all PolicyConfigurations
		// and create a map based on id
		for _, v := range p.settings["configurations"].([]interface{}) {
			// Populate
			conf, err := p.NewPolicyConfiguration(v)
			if err != nil {
				return err
			}
			p.configurations[conf.Id] = conf
		}
	}
	p.policies = make(map[string]*Policy)
	if p.settings["policies"] != nil {
		// Load all Policies
		// and populate with references to the other structures
		// based on the configured id's
		for _, v := range p.settings["policies"].([]interface{}) {
			// Populate
			policy, err := p.NewPolicy(v)
			if err != nil {
				return err
			}
			p.policies[policy.Id] = policy
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

func (p *PolicyManager) loadPolicyConditions() error {
	p.logger.Err("loadPolicyConditions Not ready yet")
	return nil
}

// Returns the priority that policy can process packets from Nfqueue
func (p *PolicyManager) PacketProcessorPriority() uint {
	return policyPriority
}

// Startup function is called to allow plugin specific initialization.
func (p *PolicyManager) Startup() error {
	p.logger.Info("Plugin %v stating\n", p.Name())
	p.syncCallbackHandler()
	return nil
}

// Shutdown function called when the daemon is shutting down
func (p *PolicyManager) Shutdown() error {
	p.logger.Info("Plugin Shutdown(%s) has been called\n", pluginName)
	return nil
}

// Returns name of the plugin
func (p *PolicyManager) Name() string {
	return pluginName
}

// syncCallbackHandler is called when the settings are changed.
// Needs to load all policy settings.
func (p *PolicyManager) syncCallbackHandler() error {
	// Do nothing for now.
	return nil
}

// PluginEnabled function returns the status (if plugin is enabled (true) or disabled (false) currently)
func (p *PolicyManager) PluginEnabled() bool {
	p.policySettingsLock.RLock()
	defer p.policySettingsLock.RUnlock()

	return p.Enabled
}

// Handles a syscall
func (p *PolicyManager) Signal(message syscall.Signal) error {
	p.logger.Info("PluginSignal(%s) has been called \n", pluginName)
	switch message {
	case syscall.SIGHUP:
		return p.syncCallbackHandler()
	}

	return nil
}

// // HandleNfqueuePacket receives a PacketMessage which includes a Tuple and
// // a gopacket.Packet, along with the IP and TCP or UDP layer already extracted.
// // We do whatever we like with the data, and when finished, we return an
// // integer via the argumented channel with interface{} bits set that we want added to
// // the packet mark.
// func (p *PolicyManager) HandleNfqueuePacket(
// 	mess dispatch.PacketMessage,
// 	newSession bool) dispatch.PacketProcessingResult {

// 	srcAddr := mess.Session.GetClientSideTuple().ClientAddress

// 	// See if this sessions should have it policy enforced.
// 	if pol := p.findPolicy(srcAddr); pol != nil {
// 		if err := p.dict.AddSessionEntry(
// 			mess.Session.GetConntrackID(),
// 			"policy",
// 			pol); err != nil {
// 			// Log something appropriate.
// 		}
// 	}

// 	// More logic probably needed in case we can't determine policy based on first packet.
// 	return dispatch.PacketProcessingResult{
// 		SessionRelease: true,
// 	}

// }

// // findPolicy returns matching polices.
// func (p *PolicyManager) findPolicy(s net.IP) *string {
// 	for _, pol := range p.policySettings {
// 		if p.matchPolicy(s, pol) {
// 			return &pol.Name
// 		}
// 	}
// 	return nil
// }

// func (p *PolicyManager) matchPolicy(sourceAddr net.IP, pol *policySettingsType) bool {
// 	// Check if source matches.

// 	for _, pSrc := range pol.Source {
// 		if _, psource, err := net.ParseCIDR(pSrc); err == nil {
// 			if psource.Contains(sourceAddr) {
// 				p.logger.Debug("Policy match: %v", pol.Name)
// 				return true
// 			}
// 		} else {
// 			p.logger.Err("Error parsing policy source: %v\n", err)
// 		}
// 	}
// 	return false
// }

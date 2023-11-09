package policy

import (
	"encoding/json"
	"fmt"

	utilNet "github.com/untangle/golang-shared/util/net"
)

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

	// Rule types
	Conditions []string `json:"conditions,omitempty"`

	// Policy Object
	Rules []string `json:"rules,omitempty"`

	//Action overlaps a bit with Policy type
	Action *Action `json:"action,omitempty"`

	// Used for policy configuration objects
	Settings any `json:"settings,omitempty"`
}

// Policies are the root of our policy configurations. It includes pointers to substructure.
type Policy = Object

type PolicyConfiguration = Object

// Action struct is used for rule object types (Conditions + Action)
type Action struct {
	Key         string `json:"key"`
	UUID        string `json:"configuration_id,omitempty"`
	Type        string `json:"type"`
	DNATAddress string `json:"dnat_address,omitempty"`
	DNATPort    string `json:"dnat_port,omitempty"`
	SNATAddress string `json:"snat_address,omitempty"`
}

// ServiceEndpoint is a particular object type, a object may be
// identified by a list of these.
type ServiceEndpoint struct {
	Protocol []uint                        `json:"protocol"`
	Port     []utilNet.PortSpecifierString `json:"port"`
}

// ApplicationObject holds an array of Ports and an array of IPSpecifiers
// a match occurs if any of the ports are matched and any of the IPs are matched
type ApplicationObject struct {
	Port       []utilNet.PortSpecifierString `json:"port"`
	IPAddrList []utilNet.IPSpecifierString   `json:"ips"`
}

// setList is a utility function for setting a list in the Object.Items field. We
// use a trick where json.Unmarshal will look at an 'any' value and if
// it has a pointer to a specific type, unmarshall into that
// type. However, we don't want the pointer later on, we just want the
// slice. setting g.Items to []T{} where T is a type we want does not
// work.
func setList[T any](obj *Object) func() {
	list := []T{}
	obj.Items = &list
	return func() {
		obj.Items = list
	}
}

// ObjectTypeField is used to figure out what group type is being used within a group
type ObjectTypeField struct {
	Type ObjectType `json:"type"`
}

// UnmarshalJSON is a custom json unmarshaller for Objects.
func (obj *Object) UnmarshalJSON(data []byte) error {
	var typeField ObjectTypeField

	if err := json.Unmarshal(data, &typeField); err != nil {
		return fmt.Errorf("unable to unmarshal group: %w", err)
	}
	type aliasObject Object

	switch typeField.Type {
	// If type field is empty - then we need to use a different type of alias to marshal (just direct object alias?)
	case "":
		// Policies typically don't have a Type
		// drop down to the default return
	case PolicyType:
		// drop to default return

	case ApplicationControlRuleObject, CaptivePortalRuleObject, GeoipRuleObject, ThreatPreventionRuleObject,
		NATRuleObject, PortForwardRuleObject, SecurityRuleObject, ShapingRuleObject, WANPolicyRuleObject,
		WebFilterRuleObject, QuotaRuleObject:
		// drop down to the default return

	case GeoipConfigType, WebFilterConfigType, ThreatPreventionConfigType,
		WANPolicyConfigType, ApplicationControlConfigType,
		CaptivePortalConfigType, SecurityConfigType:
		// drop to default return

	case QuotaType:
		obj.Settings = &QuotaSettings{}
	case WANPolicyType:
		obj.Settings = &WANPolicySettings{}
	case IPObjectType:
		defer setList[utilNet.IPSpecifierString](obj)()
	case ApplicationGroupType, GeoIPObjectType, GeoIPObjectGroupType, IPAddressGroupType, ServiceEndpointGroupType:
		defer setList[string](obj)()
	case ServiceEndpointObjectType:
		defer setList[ServiceEndpoint](obj)()
	case ApplicationType:
		defer setList[ApplicationObject](obj)()
	case InterfaceType:
		defer setList[uint](obj)()
	case InterfaceObjectType:
		defer setList[string](obj)()
	case ConditionType:
		defer setList[*PolicyCondition](obj)()
	case ConditionGroupType:
		defer setList[string](obj)()
	default:
		return fmt.Errorf("error unmarshalling policy object: invalid object type: %s", typeField.Type)
	}

	// unmarshal PolicyConfiguration using struct tags
	if err := json.Unmarshal(data, (*aliasObject)(obj)); err != nil {
		return fmt.Errorf("error unmarshalling Object of type: %s: %w", typeField.Type, err)
	}
	return nil
}

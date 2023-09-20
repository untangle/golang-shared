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

	// DEPRECATED
	Configurations []string `json:"configurations,omitempty"`
	Flows          []string `json:"flows,omitempty"`
}

// Group is a deprecated concept, please use Object.
// Deprecated: Group is deprecated, use Object instead. See MFW-3517.
type Group = Object

// Policies are the root of our policy configurations. It includes pointers to substructure.
type Policy = Object

// Action struct is used for rule object types (Conditions + Action)
type Action struct {
	Key  string `json:"key"`
	UUID string `json:"configuration_id"`
	Type string `json:"type"`
}

// ServiceEndpoint is a particular group type, a group may be
// identified by a list of these.
type ServiceEndpoint struct {
	Protocol uint `json:"protocol"`
	Port     uint `json:"port"`
}

// setList is a utility function for setting a list in the Group.Items field. We
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

// UnmarshalJSON is a custom json unmarshaller for Objects.
func (obj *Object) UnmarshalJSON(data []byte) error {
	var typeField GroupTypeField

	if err := json.Unmarshal(data, &typeField); err != nil {
		return fmt.Errorf("unable to unmarshal group: %w", err)
	}
	type aliasObject Object

	switch typeField.Type {
	// If type field is empty - then we need to use a different type of alias to marshal (just direct object alias?)
	case "":
		// Policies typically don't have a Type
		// drop down to the defaul return

	case ApplicationControlRuleObject, CaptivePortalRuleObject, GeoipRuleObject,
		NATRuleObject, PortForwardRuleObject, SecurityRuleObject, ShapingRuleObject, WANPolicyRuleObject:
		// drop down to the defaul return

	case IPAddrListType, IPObjectType:
		defer setList[utilNet.IPSpecifierString](obj)()
	case GeoIPListType, GeoIPObjectType, GeoIPObjectGroupType:
		defer setList[string](obj)()
	case ServiceEndpointType, ServiceEndpointObjectType:
		defer setList[ServiceEndpoint](obj)()
	case InterfaceType, InterfaceObjectType:
		defer setList[uint](obj)()
	case ConditionType:
		defer setList[*PolicyCondition](obj)()
	case ConditionGroupType:
		defer setList[string](obj)()
	case ThreatPreventionType:
		defer setList[uint](obj)()
	case WebFilterCategoryType:
		defer setList[uint](obj)()
	default:
		return fmt.Errorf("error unmarshalling policy group: invalid group type: %s", typeField.Type)
	}

	// unmarshal PolicyConfiguration using struct tags
	return json.Unmarshal(data, (*aliasObject)(obj))
}

// ItemsStringList returns the Items of the group as a slice of
// strings if they can be interpreted this way, or an empty slice and
// false if not.
func (g *Group) ItemsStringList() ([]string, bool) {
	val, ok := g.Items.([]string)
	return val, ok
}

// ItemsIPSpecList returns the Items of a group as a slice of
// utilNet.IPSpecifierString and true if they can be interpreted this way,
// or an empty slice and false otherwise.
func (g *Group) ItemsIPSpecList() ([]utilNet.IPSpecifierString, bool) {
	val, ok := g.Items.([]utilNet.IPSpecifierString)
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

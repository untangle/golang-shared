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
	Name        string    `json:"name"`
	Type        GroupType `json:"type"`
	Description string    `json:"description"`
	ID          string    `json:"id"`
	Items       any       `json:"items"`
}

// Group is a deprecated concept, please use Object.
// Deprecated: Group is deprecated, use Object instead. See MFW-3517.
type Group = Object

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
		defer setList[utilNet.IPSpecifierString](g)()
	case GeoIPListType:
		defer setList[string](g)()
	case ServiceEndpointType:
		defer setList[ServiceEndpoint](g)()
	case InterfaceType:
		defer setList[uint](g)()
	case ThreatPreventionType:
		defer setList[uint](g)()
	case WebFilterCategoryType:
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

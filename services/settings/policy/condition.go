package policy

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// Empty interface for optoinal data in PolicyCondition
type ConditionExtraData interface{}

// PolicyCondition contains policy condition configuration.
type PolicyCondition struct {
	Op       string   `json:"op"`
	CType    string   `json:"type"`
	Value    []string `json:"value,omitempty"`
	GroupIDs []string `json:"object,omitempty"`

	// Extra Data for ApplicationObjectGroup or other things
	Extra ConditionExtraData `json:"omitempty"`

	// Deprecated
	GroupID string `json:"groupId,omitempty"`
}

// Application Object
type ApplicationObject struct {
	ConditionExtraData
	// Not sure we need protocol
	// If not present, treat as TCP
	Ports      []int    `json:"ports"`
	IPAddrList []string `json:"ipaddresslist"`
}

// Application Object Group
type ApplicationObjectGroup struct {
	Items []ApplicationObject `json:"items"`
}

var _ ConditionExtraData = ApplicationObjectGroup{}

// Unmarshal policy condition so that types of values can be checked
func (pCondition *PolicyCondition) UnmarshalJSON(data []byte) error {
	// unmarshal like normal first
	type aliasPolicyCondition PolicyCondition
	if err := json.Unmarshal(data, (*aliasPolicyCondition)(pCondition)); err != nil {
		return err
	}
	// Check to see if this is Conditions is using an ApplicationObject or ...Group
	if pCondition.CType == "APPLICATION_NAME" {
		// Using an ApplicationObject or ...Group
		// Try for an ApplicationObjecGroup with Items first
		appObjectGroup := ApplicationObjectGroup{}
		if err := json.Unmarshal(data, &appObjectGroup); err != nil {
			return err
		}
		if len(appObjectGroup.Items) == 0 {
			// Otherwise try for a single ApplicationObject
			// ApplicatoinObjectGroup with exactly 1 item
			appObjectGroup = ApplicationObjectGroup{
				Items: make([]ApplicationObject, 1),
			}
			if err := json.Unmarshal(data, &appObjectGroup.Items[0]); err != nil {
				return err
			}
		}
		pCondition.Extra = appObjectGroup
	} else if pCondition.Op != "match" && pCondition.Op != "in" {
		// check that pCondition.Value is formatted correctly for the CType
		for i, value := range pCondition.Value {
			switch pCondition.CType {
			case "CLIENT_ADDRESS", "SERVER_ADDRESS", "SOURCE_ADDRESS", "DESTINATION_ADDRESS":
				// Check that address is in CIDR format (w/ mask)
				if _, _, err := net.ParseCIDR(value); err != nil {
					// If address is a valid IP, but without a mask, just add the default
					if ip := net.ParseIP(value); ip != nil {
						if ip.To4() != nil {
							pCondition.Value[i] = fmt.Sprintf("%s%s", value, "/32")
						} else {
							pCondition.Value[i] = fmt.Sprintf("%s%s", value, "/128")
						}
					} else {
						return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
					}
				}
			case "CLIENT_PORT", "SERVER_PORT":
				if _, err := strconv.ParseUint(value, 10, 32); err != nil {
					return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
				}
			// just string type values on these, no need to validate
			case "DAY_OF_WEEK", "GEOIP_LOCATION", "INTERFACE", "SERVICE_ENDPOINT",
				"PROTOCOL_TYPE", "TIME_OF_DAY", "VLAN_ID", "THREATPREVENTION":

			default:
				return fmt.Errorf("error while unmarshalling policy condition: invalid type: %s", pCondition.CType)
			}
		}
	}
	return nil
}

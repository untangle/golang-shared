package policy

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// PolicyCondition contains policy condition configuration.
type PolicyCondition struct {
	Op       string   `json:"op"`
	CType    string   `json:"type"`
	Value    []string `json:"value,omitempty"`
	GroupIDs []string `json:"object,omitempty"`
}

// Unmarshal policy condition so that types of values can be checked
func (pCondition *PolicyCondition) UnmarshalJSON(data []byte) error {
	// unmarshal like normal first
	type aliasPolicyCondition PolicyCondition
	alias := &struct {
		ValueRaw json.RawMessage `json:"value"`
		*aliasPolicyCondition
	}{
		aliasPolicyCondition: (*aliasPolicyCondition)(pCondition),
	}
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}

	if alias.ValueRaw != nil {
		var valString []string
		if err := json.Unmarshal(alias.ValueRaw, &valString); err == nil {
			pCondition.Value = valString
		} else {
			var valInt []int
			if err := json.Unmarshal(alias.ValueRaw, &valInt); err != nil {
				return err
			}

			for _, v := range valInt {
				pCondition.Value = append(pCondition.Value, strconv.Itoa(v))
			}
		}
	}

	// Only use value if Group is not configured
	switch pCondition.Op {
	case "in", "match", "not_in", "not_match":
		// Special handling for objects
		// The Condition will contain one or more GUIDs in its GroupIDs array
	default:
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
			case "IP_PROTOCOL", "CLIENT_PORT", "SERVER_PORT",
				// CLIENT and SOURCE mean the same thing - support both
				// SERVER and DESTINATION mean the same thing - support both
				"CLIENT_INTERFACE_TYPE", "SERVER_INTERFACE_TYPE",
				"SOURCE_INTERFACE_TYPE", "DESTINATION_INTERFACE_TYPE",
				"APPLICATION_RISK", "APPLICATION_RISK_INFERRED",
				"APPLICATION_PRODUCTIVITY", "APPLICATION_PRODUCTIVITY_INFERRED":

				if _, err := strconv.ParseUint(value, 10, 32); err != nil {
					return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
				}
			// just string type values on these, no need to validate
			case "CERT_SUBJECT_CN", "CERT_SUBJECT_DNS", "CERT_SUBJECT_O",
				"DAY_OF_WEEK", "SERVER_GEOIP", "CLIENT_GEOIP", "INTERFACE", "SERVICE", "SERVER_SERVICE", "CLIENT_SERVICE",
				// CLIENT and SOURCE mean the same thing - support both
				// SERVER and DESTINATION mean the same thing - support both
				"CLIENT_INTERFACE_NAME", "SERVER_INTERFACE_NAME",
				"SOURCE_INTERFACE_NAME", "DESTINATION_INTERFACE_NAME",
				"CLIENT_INTERFACE_ZONE", "SERVER_INTERFACE_ZONE",
				"SOURCE_INTERFACE_ZONE", "DESTINATION_INTERFACE_ZONE",
				"SOURCE_INTERFACE", "DESTINATION_INTERFACE",

				"PROTOCOL_TYPE", "TIME_OF_DAY", "VLAN_TAG", "THREATPREVENTION",
				"APPLICATION", "SERVER_APPLICATION", "CLIENT_APPLICATION", "HOSTNAME", "SERVER_DNS_HINT", "CLIENT_DNS_HINT",
				"APPLICATION_NAME", "APPLICATION_NAME_INFERRED", "APPLICATION_CATEGORY", "APPLICATION_CATEGORY_INFERRED":

			default:
				// At the moment we allow undeclared fields.
			}
		}
	}
	return nil
}

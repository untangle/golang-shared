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
	if err := json.Unmarshal(data, (*aliasPolicyCondition)(pCondition)); err != nil {
		return err
	}

	// Only use value if Group is not configured
	if pCondition.Op != "match" && pCondition.Op != "in" {
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
				"SOURCE_INTERFACE_TYPE", "DESTINATION_INTERFACE_TYPE":
				if _, err := strconv.ParseUint(value, 10, 32); err != nil {
					return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
				}
			// just string type values on these, no need to validate
			case "CERT_SUBJECT_CN", "CERT_SUBJECT_DNS", "CERT_SUBJECT_O",
				"DAY_OF_WEEK", "SERVER_GEOIP", "CLIENT_GEOIP", "INTERFACE", "SERVICE", "SERVER_SERVICE", "CLIENT_SERVICE",
				"DESINTATION_INTERFACE_NAME", "DESTINATION_INTERFACE_TYPE", "DESTINATION_INTERFACE_ZONE",
				"SOURCE_INTERFACE_NAME", "SOURCE_INTERFACE_TYPE", "SOURCE_INTERFACE_ZONE",
				"PROTOCOL_TYPE", "APPLICATION_CATEGORY", "TIME_OF_DAY", "VLAN_TAG", "THREATPREVENTION",
				"APPLICATION", "SERVER_APPLICATION", "CLIENT_APPLICATION", "HOSTNAME":

			default:
				return fmt.Errorf("error while unmarshalling policy condition: invalid type: %s", pCondition.CType)
			}
		}
	}

	return nil
}

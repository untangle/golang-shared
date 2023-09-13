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

	// check that pCondition.Value is formatted correctly for the CType
	for i, value := range pCondition.Value {
		switch pCondition.CType {
		case "CLIENT_ADDRESS", "SERVER_ADDRESS":
			// Check that address is in CIDR format (w/ mask)
			if _, _, err := net.ParseCIDR(value); err != nil {
				// If address is a valid IP, but without a mask, just add the default
				if ip := net.ParseIP(value); ip != nil {
					if ip.To4() != nil {
						pCondition.Value[i] = fmt.Sprintf("%s%s", value, "/32")
					} else {
						pCondition.Value[i] = fmt.Sprintf("%s%s", value, "/64")
					}
				} else {
					return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
				}
			}
		case "CLIENT_PORT", "SERVER_PORT":
			if _, err := strconv.ParseUint(value, 10, 32); err != nil {
				return fmt.Errorf("error while unmarshalling policy condition: value does not match type (%s) due to error (%v)", pCondition.CType, err)
			}
		case "DAY_OF_WEEK", "DEST_ADDRESS", "GEOIP_LOCATION", "INTERFACE", "SERVICE_ENDPOINT", "SOURCE_ADDRESS", "PROTOCOL_TYPE", "TIME_OF_DAY", "VLAN_ID", "THREATPREVENTION":
			// These are not yet implemented and need to have a designated format
		default:
			return fmt.Errorf("error while unmarshalling policy condition: invalid type: %s", pCondition.CType)
		}
	}

	return nil
}

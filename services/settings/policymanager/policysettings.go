package policy

var policyConditionTypeMap = map[string]bool{
	"CLIENT_ADDRESS": true,
	"CLIENT_PORT":    true,
	"DAY_OF_WEEK":    true,
	"DEST_ADDRESS":   true,
	"INTERFACE":      true,
	"SERVER_ADDRESS": true,
	"SERVER_PORT":    true,
	"SOURCE_ADDRESS": true,
	"PROTOCOL_TYPE":  true,
	"TIME_OF_DAY":    true,
	"VLAN_ID":        true,
}

// Valid PolicyCondition Ops - there may be more at some point
// == implies an OR operation between the different entries in the value arrray
// != implies an AND operation between the different entries in the value array
// all other operations assume a single entry in the value array (or string)
var policyConditionOpsMap = map[string]bool{
	"==": true,
	"!=": true,
	"<":  true,
	">":  true,
	"<=": true,
	">=": true,
}

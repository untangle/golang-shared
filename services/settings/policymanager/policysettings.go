package policy

var policyConditionTypeMap = map[string]int{
	"CLIENT_ADDRESS": 1,
	"CLIENT_PORT":    1,
	"DAY_OF_WEEK":    1,
	"DEST_ADDRESS":   1,
	"INTERFACE":      1,
	"SERVER_ADDRESS": 1,
	"SERVER_PORT":    1,
	"SOURCE_ADDRESS": 1,
	"PROTOCOL_TYPE":  1,
	"TIME_OF_DAY":    1,
	"VLAN_ID":        1,
}

// Valid PolicyCondition Ops - there may be more at some point
// == implies an OR operation between the different entries in the value arrray
// != implies an AND operation between the different entries in the value array
// all other operations assume a single entry in the value array (or string)
var policyConditionOpsMap = map[string]int{"==": 1, "!=": 1, "<": 1, ">": 1, "<=": 1, ">=": 1}

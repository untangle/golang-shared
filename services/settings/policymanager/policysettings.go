package policy

var policyConditionTypeMap = map[string]true{
	"CLIENT_ADDRESS": bool,
	"CLIENT_PORT":    bool,
	"DAY_OF_WEEK":    bool,
	"DEST_ADDRESS":   bool,
	"INTERFACE":      bool,
	"SERVER_ADDRESS": bool,
	"SERVER_PORT":    bool,
	"SOURCE_ADDRESS": bool,
	"PROTOCOL_TYPE":  bool,
	"TIME_OF_DAY":    bool,
	"VLAN_ID":        bool,
}

// Valid PolicyCondition Ops - there may be more at some point
// == implies an OR operation between the different entries in the value arrray
// != implies an AND operation between the different entries in the value array
// all other operations assume a single entry in the value array (or string)
var policyConditionOpsMap = map[string]bool{
	"==": bool,
	"!=": bool,
	"<":  bool,
	">":  bool,
	"<=": bool,
	">=": bool,
}

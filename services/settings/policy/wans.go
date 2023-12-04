package policy

// WANCriteriaType is a type of WAN criteria
type WANCriteriaType struct {
	Type      string `json:"type"`
	Attribute string `json:"attribute"`
}

// WANInterfaceType is a type of WAN interface
type WANInterfaceType struct {
	ID uint `json:"interfaceId"`
}

// WANPolicy is an object with Type Object
type WANPolicy Object

// WANPolicySettings are settings for a WAN policy.
type WANPolicySettings struct {
	Criteria     []WANCriteriaType  `json:"criteria"`
	Interfaces   []WANInterfaceType `json:"interfaces"`
	Type         string             `json:"type"`
	BestOfMetric string             `json:"best_of_metric"`
	Description  string             `json:"description"`
	Name         string             `json:"name"`
}

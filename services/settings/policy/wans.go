package policy

type WANCriteriaType struct {
	Type string `json:"type"`
}

type WANInterfaceType struct {
	ID uint `json:"interfaceId"`
}

// WANPolicySettings are settings for a WAN policy.
type WANPolicySettings struct {
	Criteria     []WANCriteriaType  `json:"criteria"`
	Interfaces   []WANInterfaceType `json:"interfaces"`
	Type         string             `json:"type"`
	BestOfMetric string             `json:"best_of_metric"`
}

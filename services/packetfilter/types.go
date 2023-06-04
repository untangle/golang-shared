package packetfilter

import "github.com/untangle/golang-shared/booleval"

const (
	// Actions
	ActionAccept = iota
	ActionReject
	// Properties
	ApplicationCategoryInferred
	ApplicationIDInferred
	ApplicationNameInferred
	ApplicationProtocolInferred
	ApplicationRiskInferred
	CertSubjectCN
	CertSubjectDNS
	CertSubjectO
	ClientAddress
	ClientAddressV6
	ClientInterfaceType
	ClientInterfaceZone
	ClientPort
	ClientReverseDNS
	CTState
	IPProtocol
	ServerAddress
	ServerAddressV6
	ServerDNSHint
	ServerInterfaceType
	ServerInterfaceZone
	ServerPort
)

// PacketFilter is a set of conditions and an action to take if the conditions are true
type PacketFilter struct {
	Name              string                      // Name of the filter
	Enabled           bool                        // Whether the filter is enabled
	Conditions        []PacketCondition           // Conditions to evaluate
	AtomicExpressions []booleval.AtomicExpression // Expression translated from conditions
	Action            string                      // Action to take if the conditions is true
}

type PacketCondition struct {
	Property string
	Operator string
	Value    string
}

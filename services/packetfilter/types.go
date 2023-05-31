package packetfilter

import (
	"github.com/untangle/golang-shared/booleval"
)

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

type PacketFilter struct {
	Conditions []booleval.AtomicExpression // Conditions to evaluate
	Action     int                         // Action to take if the condition is true
}

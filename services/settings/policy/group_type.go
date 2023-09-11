package policy

// GroupType is the type of group that a Group is, used to demux the
// Items field.
type GroupType string

const (
	// GeoIPListType means that the Items of a Group are geoip countries.
	GeoIPListType GroupType = "GeoIPLocation"

	// IPAddrListType means that the Items of the Group are ip
	// specifications (ranges, CIDRs, or single IPs).
	IPAddrListType GroupType = "IPAddrList"

	// ServiceEndpointType means that the Items of a Group are
	// service endpoints.
	ServiceEndpointType GroupType = "ServiceEndpoint"

	// InterfaceType is a group type where all items are interface
	// IDs (integers)
	InterfaceType GroupType = "Interface"

	// WebFilterCategoryType means that the Items of the Group are web filter categories.
	WebFilterCategoryType GroupType = "WebFilterCategory"

	// ThreatPreventionType means that the Items of the Group are threat prevention score.
	ThreatPreventionType GroupType = "ThreatPrevention"
)

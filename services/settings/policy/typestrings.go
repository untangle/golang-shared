package policy

// ObjectType is a string used to demux the actual type of an object
// when loading from JSON.
type ObjectType string

// GroupType is the type of group that a Group is, used to demux the
// Items field.
type GroupType = ObjectType

// GroupTypeField is used to figure out what group type is being used within a group
type GroupTypeField struct {
	Type GroupType `json:"type"`
}

const (
	// Deprecated: GeoIPListType means that the Items of a Group are geoip countries.
	GeoIPListType ObjectType = "GeoIPLocation"

	// GeoIPObjectGroupType/GeoIPListType are new-style type names for
	// geoip objects and groups.
	GeoIPObjectType       ObjectType = "mfw-object-geoip"
	GeoIPObjectGroupType  ObjectType = "mfw-object-geoip-group"
	GeoipFilterRuleObject ObjectType = "GeoipFilterRuleObject"

	// Deprecated: IPAddrListType means that the Items of the Group are ip
	// specifications (ranges, CIDRs, or single IPs). (old)
	IPAddrListType ObjectType = "IPAddrList"

	// IPObjectType/IPAddressGroupType are the types for the
	// new-style matchable object/group that relate to IPs.
	IPObjectType       ObjectType = "mfw-object-ipaddress"
	IPAddressGroupType ObjectType = "mfw-object-ipaddress-group"

	// ServiceEndpointType means that the Items of a Group are
	// service endpoints.
	ServiceEndpointType ObjectType = "ServiceEndpoint"

	// ServiceEndpointObjecttype and ServiceEndpointGroup are types
	// for object/group in new schema, from cloud.
	ServiceEndpointObjectType ObjectType = "mfw-object-service"
	ServiceEndpointGroupType  ObjectType = "mfw-object-service-group"

	// InterfaceType is a group type where all items are interface
	// IDs (integers)
	InterfaceType            ObjectType = "Interface"
	InterfaceObjectType      ObjectType = "mfw-interfacezone-object"
	InterfaceObjectGroupType ObjectType = "mfw-interfacezone-group"

	// Deprecated: WebFilter* will be removed.
	// WebFilterCategoryType means that the Items of the Group are web filter categories.
	WebFilterCategoryType ObjectType = "WebFilterCategory"
	WebFilterRuleObject   ObjectType = "WebFilterRuleObject"

	// Deprecated: ThreatPreventionType will be removed.
	// ThreatPreventionType means that the Items of the Group are
	// threat prevention score.
	ThreatPreventionType ObjectType = "ThreatPrevention"

	// ConditionType,ConditionGrouType: type id strings, for the
	// object and the group.
	ConditionType      ObjectType = "mfw-object-condition"
	ConditionGroupType ObjectType = "mfw-object-condition-group"

	// Configuration Types used for marshalling configs out of configuration sections
	GeoipConfigType              ObjectType = "mfw-template-geoipfilter"
	WebFilterConfigType          ObjectType = "mfw-template-webfilter"
	ThreatPreventionConfigType   ObjectType = "mfw-template-threatprevention"
	WANPolicyConfigType          ObjectType = "mfw-template-wanpolicy"
	ApplicationControlConfigType ObjectType = "mfw-template-applicationcontrol"
	CaptivePortalConfigType      ObjectType = "mfw-template-captiveportal"
	SecurityConfigType           ObjectType = "mfw-template-security"

	// TODO: Impelment these object/group types
	HostType      ObjectType = "mfw-object-host"
	HostGroupType ObjectType = "mfw-object-host-group"

	DomainType      ObjectType = "mfw-object-domain"
	DomainGroupType ObjectType = "mfw-object-domain-group"

	UserType      ObjectType = "mfw-object-user"
	UserGroupType ObjectType = "mfw-object-user-group"

	VLANTagType      ObjectType = "mfw-object-vlantag"
	VLANTagGroupType ObjectType = "mfw-object-vlantag-group"

	ApplicationType      ObjectType = "mfw-object-application"
	ApplicationGroupType ObjectType = "mfw-object-application-group"
)

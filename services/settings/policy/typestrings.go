package policy

// ObjectType is a string used to demux the actual type of an object
// when loading from JSON.
type ObjectType string

// ObjectParentType is a string used for representing the parent grouping of
// a specific object
type ObjectParentType string

// GroupType is the type of group that a Group is, used to demux the
// Items field.
type GroupType = ObjectType

// GroupTypeField is used to figure out what group type is being used within a group
type GroupTypeField struct {
	Type GroupType `json:"type"`
}

// ObjectMetaLookup is a map for retrieving additional metadata about objects
var ObjectMetaLookup map[ObjectType]ObjectMetadata

// SettingsMetaLookup is a map for retrieving metadata from a specific settings
var SettingsMetaLookup map[string]ObjectMetadata

// ObjectMetadata keeps track of a set of metadata for each enumeration value
// for object enumerations.
type ObjectMetadata struct {
	// The object type
	Type ObjectType
	// The object parent type
	ParentType ObjectParentType
	// The original settings name
	SettingsName string
}

const (
	// Deprecated: GeoIPListType means that the Items of a Group are geoip countries.
	GeoIPListType ObjectType = "GeoIPLocation"

	// GeoIPObjectGroupType/GeoIPListType are new-style type names for
	// geoip objects and groups.
	GeoIPObjectType      ObjectType = "mfw-object-geoip"
	GeoIPObjectGroupType ObjectType = "mfw-object-geoip-group"

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

	// RuleTypes
	ApplicationControlRuleObject ObjectType = "mfw-rule-applicationcontrol"
	CaptivePortalRuleObject      ObjectType = "mfw-rule-captiveportal"
	GeoipRuleObject              ObjectType = "mfw-rule-geoip"
	NATRuleObject                ObjectType = "mfw-rule-nat"
	PortForwardRuleObject        ObjectType = "mfw-rule-portforward"
	SecurityRuleObject           ObjectType = "mfw-rule-security"
	ShapingRuleObject            ObjectType = "mfw-rule-shaping"
	WANPolicyRuleObject          ObjectType = "mfw-rule-wanpolicy"

	// Deprecated: WebFilter* will be removed.
	// WebFilterCategoryType means that the Items of the Group are web filter categories.
	WebFilterCategoryType ObjectType = "WebFilterCategory"
	WebFilterRuleObject   ObjectType = "WebFilterRuleObject"

	// Deprecated: ThreatPreventionType will be removed.
	// ThreatPreventionType means that the Items of the Group are
	// threat prevention score.
	ThreatPreventionType ObjectType = "ThreatPrevention"

	// ConditionType,ConditionGroupType: type id strings, for the
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

	PolicyParent         ObjectParentType = "policy"
	RuleParent           ObjectParentType = "rule"
	ConditionParent      ObjectParentType = "condition"
	ConditionGroupParent ObjectParentType = "conditiongroup"
	ConfigurationParent  ObjectParentType = "configuration"
	ObjectParent         ObjectParentType = "object"
	ObjectGroupParent    ObjectParentType = "objectgroup"

	GeoipSettingsKey      string = "geoip"
	WebfilterSettingsKey  string = "webfilter"
	TPSettingsKey         string = "threatprevention"
	AppControlSettingsKey string = "application_control"
	CaptiveSettingsKey    string = "captiveportal"
)

// init() builds the object metadata map
func init() {
	buildObjectMetadata()
}

// buildObjectMetadata constructs a metadata map used for retrieving details or properties of specific objects
func buildObjectMetadata() {

	// we can probably size these since we know all the values
	ObjectMetaLookup = make(map[ObjectType]ObjectMetadata)
	SettingsMetaLookup = make(map[string]ObjectMetadata)

	// Configs exist in both the SettingsMetaLookup and ObjectMetaLookup - so that we can easily translate settings config name -> template meta details
	var geoipMeta ObjectMetadata = ObjectMetadata{SettingsName: GeoipSettingsKey, Type: GeoipConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup[GeoipSettingsKey] = geoipMeta
	ObjectMetaLookup[GeoipConfigType] = geoipMeta

	var webfilterMeta ObjectMetadata = ObjectMetadata{SettingsName: WebfilterSettingsKey, Type: WebFilterConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup[WebfilterSettingsKey] = webfilterMeta
	ObjectMetaLookup[WebFilterConfigType] = webfilterMeta

	var tpMeta ObjectMetadata = ObjectMetadata{SettingsName: TPSettingsKey, Type: ThreatPreventionConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup[TPSettingsKey] = tpMeta
	ObjectMetaLookup[ThreatPreventionConfigType] = tpMeta

	var appMeta ObjectMetadata = ObjectMetadata{SettingsName: AppControlSettingsKey, Type: ApplicationControlConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup[AppControlSettingsKey] = appMeta
	ObjectMetaLookup[ApplicationControlConfigType] = appMeta

	var captiveMeta ObjectMetadata = ObjectMetadata{SettingsName: CaptiveSettingsKey, Type: CaptivePortalConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup[CaptiveSettingsKey] = captiveMeta
	ObjectMetaLookup[CaptivePortalConfigType] = captiveMeta

	// do we really need these? they don't have 'default configuration' per se
	var securityMeta ObjectMetadata = ObjectMetadata{SettingsName: "security", Type: SecurityConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup["security"] = securityMeta
	ObjectMetaLookup[SecurityConfigType] = securityMeta

	var wanMeta ObjectMetadata = ObjectMetadata{SettingsName: "wan_policy", Type: WANPolicyConfigType, ParentType: ConfigurationParent}
	SettingsMetaLookup["wan_policy"] = wanMeta
	ObjectMetaLookup[WANPolicyConfigType] = wanMeta
}

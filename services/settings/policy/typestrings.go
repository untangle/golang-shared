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
	// GeoIPObjectGroupType/GeoIPListType are new-style type names for
	// geoip objects and groups.
	GeoIPObjectType      ObjectType = "mfw-object-geoip"
	GeoIPObjectGroupType ObjectType = "mfw-object-geoip-group"

	// IPObjectType/IPAddressGroupType are the types for the
	// new-style matchable object/group that relate to IPs.
	IPObjectType       ObjectType = "mfw-object-ipaddress"
	IPAddressGroupType ObjectType = "mfw-object-ipaddress-group"

	// ServiceEndpointObjecttype and ServiceEndpointGroup are types
	// for object/group in new schema, from cloud.
	ServiceEndpointObjectType ObjectType = "mfw-object-service"
	ServiceEndpointGroupType  ObjectType = "mfw-object-service-group"

	// InterfacezoneObjectType is a group type where all items are interface
	// IDs (integers)
	InterfacezoneObjectType      ObjectType = "mfw-object-interfacezone"
	InterfacezoneObjectGroupType ObjectType = "mfw-object-interfacezone-group"

	// QuotaType string -- a quota type.
	QuotaType ObjectType = "mfw-quota"

	// RuleTypes
	ApplicationControlRuleObject ObjectType = "mfw-rule-applicationcontrol"
	CaptivePortalRuleObject      ObjectType = "mfw-rule-captiveportal"
	GeoipRuleObject              ObjectType = "mfw-rule-geoip"
	NATRuleObject                ObjectType = "mfw-rule-nat"
	PortForwardRuleObject        ObjectType = "mfw-rule-portforward"
	SecurityRuleObject           ObjectType = "mfw-rule-security"
	ShapingRuleObject            ObjectType = "mfw-rule-shaping"
	ThreatPreventionRuleObject   ObjectType = "mfw-rule-threatprevention"
	WANPolicyRuleObject          ObjectType = "mfw-rule-wanpolicy"
	WebFilterRuleObject          ObjectType = "mfw-rule-webfilter"
	QuotaRuleObject              ObjectType = "mfw-rule-quota"

	// ConditionType,ConditionGroupType: type id strings, for the
	// object and the group.
	ConditionType      ObjectType = "mfw-object-condition"
	ConditionGroupType ObjectType = "mfw-object-condition-group"

	// Configuration Types used for marshalling configs out of configuration sections
	GeoipConfigType              ObjectType = "mfw-config-geoipfilter"
	WebFilterConfigType          ObjectType = "mfw-config-webfilter"
	ThreatPreventionConfigType   ObjectType = "mfw-config-threatprevention"
	WANPolicyConfigType          ObjectType = "mfw-config-wanpolicy"
	ApplicationControlConfigType ObjectType = "mfw-config-applicationcontrol"
	CaptivePortalConfigType      ObjectType = "mfw-config-captiveportal"
	SecurityConfigType           ObjectType = "mfw-config-security"

	//Policy Type
	PolicyType ObjectType = "mfw-policy"

	HostType      ObjectType = "mfw-object-hostname"
	HostGroupType ObjectType = "mfw-object-hostname-group"

	DomainType      ObjectType = "mfw-object-domain"
	DomainGroupType ObjectType = "mfw-object-domain-group"

	VLANTagType      ObjectType = "mfw-object-vlantag"
	VLANTagGroupType ObjectType = "mfw-object-vlantag-group"

	ApplicationType      ObjectType = "mfw-object-application"
	ApplicationGroupType ObjectType = "mfw-object-application-group"

	// TODO: implemented fully when the UserType diverges from IPObjectType.
	// As of right now, they are functionally the same so the front end is using
	// IPObjectType
	UserType      ObjectType = "mfw-object-user"
	UserGroupType ObjectType = "mfw-object-user-group"

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

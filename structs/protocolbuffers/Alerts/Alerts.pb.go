// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: Alerts.proto

package Alerts

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AlertType int32

const (
	AlertType_UNKNOWN           AlertType = 0
	AlertType_USER              AlertType = 1
	AlertType_LINK              AlertType = 2
	AlertType_THREATPREVENTION  AlertType = 3
	AlertType_WEBFILTER         AlertType = 4
	AlertType_WEBCLASSIFICATION AlertType = 5
	AlertType_GEOIP             AlertType = 6
	AlertType_SETTINGS          AlertType = 7
	AlertType_DISCOVERY         AlertType = 8
	AlertType_DHCP              AlertType = 9
	AlertType_CRITICALERROR     AlertType = 10
	AlertType_VPN               AlertType = 11
	AlertType_CAPTIVEPORTAL     AlertType = 12
	AlertType_FIREWALLEVENT     AlertType = 13
	AlertType_DYNAMICLISTS      AlertType = 14
	AlertType_POLICYMANAGER     AlertType = 15
	AlertType_DATABASEMANAGER   AlertType = 16
	AlertType_DNSFILTERMANAGER  AlertType = 17
)

// Enum value maps for AlertType.
var (
	AlertType_name = map[int32]string{
		0:  "UNKNOWN",
		1:  "USER",
		2:  "LINK",
		3:  "THREATPREVENTION",
		4:  "WEBFILTER",
		5:  "WEBCLASSIFICATION",
		6:  "GEOIP",
		7:  "SETTINGS",
		8:  "DISCOVERY",
		9:  "DHCP",
		10: "CRITICALERROR",
		11: "VPN",
		12: "CAPTIVEPORTAL",
		13: "FIREWALLEVENT",
		14: "DYNAMICLISTS",
		15: "POLICYMANAGER",
		16: "DATABASEMANAGER",
		17: "DNSFILTERMANAGER",
	}
	AlertType_value = map[string]int32{
		"UNKNOWN":           0,
		"USER":              1,
		"LINK":              2,
		"THREATPREVENTION":  3,
		"WEBFILTER":         4,
		"WEBCLASSIFICATION": 5,
		"GEOIP":             6,
		"SETTINGS":          7,
		"DISCOVERY":         8,
		"DHCP":              9,
		"CRITICALERROR":     10,
		"VPN":               11,
		"CAPTIVEPORTAL":     12,
		"FIREWALLEVENT":     13,
		"DYNAMICLISTS":      14,
		"POLICYMANAGER":     15,
		"DATABASEMANAGER":   16,
		"DNSFILTERMANAGER":  17,
	}
)

func (x AlertType) Enum() *AlertType {
	p := new(AlertType)
	*p = x
	return p
}

func (x AlertType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AlertType) Descriptor() protoreflect.EnumDescriptor {
	return file_Alerts_proto_enumTypes[0].Descriptor()
}

func (AlertType) Type() protoreflect.EnumType {
	return &file_Alerts_proto_enumTypes[0]
}

func (x AlertType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AlertType.Descriptor instead.
func (AlertType) EnumDescriptor() ([]byte, []int) {
	return file_Alerts_proto_rawDescGZIP(), []int{0}
}

type AlertSeverity int32

const (
	AlertSeverity_INFO     AlertSeverity = 0
	AlertSeverity_WARN     AlertSeverity = 1
	AlertSeverity_ERROR    AlertSeverity = 2
	AlertSeverity_DEBUG    AlertSeverity = 3
	AlertSeverity_CRITICAL AlertSeverity = 4
)

// Enum value maps for AlertSeverity.
var (
	AlertSeverity_name = map[int32]string{
		0: "INFO",
		1: "WARN",
		2: "ERROR",
		3: "DEBUG",
		4: "CRITICAL",
	}
	AlertSeverity_value = map[string]int32{
		"INFO":     0,
		"WARN":     1,
		"ERROR":    2,
		"DEBUG":    3,
		"CRITICAL": 4,
	}
)

func (x AlertSeverity) Enum() *AlertSeverity {
	p := new(AlertSeverity)
	*p = x
	return p
}

func (x AlertSeverity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AlertSeverity) Descriptor() protoreflect.EnumDescriptor {
	return file_Alerts_proto_enumTypes[1].Descriptor()
}

func (AlertSeverity) Type() protoreflect.EnumType {
	return &file_Alerts_proto_enumTypes[1]
}

func (x AlertSeverity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AlertSeverity.Descriptor instead.
func (AlertSeverity) EnumDescriptor() ([]byte, []int) {
	return file_Alerts_proto_rawDescGZIP(), []int{1}
}

type Alert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type          AlertType         `protobuf:"varint,1,opt,name=type,proto3,enum=alerts.AlertType" json:"type,omitempty"`
	Severity      AlertSeverity     `protobuf:"varint,2,opt,name=severity,proto3,enum=alerts.AlertSeverity" json:"severity,omitempty"`
	Message       string            `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Params        map[string]string `protobuf:"bytes,4,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Datetime      int64             `protobuf:"varint,5,opt,name=datetime,proto3" json:"datetime,omitempty"`
	IsLoggerAlert bool              `protobuf:"varint,6,opt,name=isLoggerAlert,proto3" json:"isLoggerAlert,omitempty"`
}

func (x *Alert) Reset() {
	*x = Alert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Alerts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Alert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Alert) ProtoMessage() {}

func (x *Alert) ProtoReflect() protoreflect.Message {
	mi := &file_Alerts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Alert.ProtoReflect.Descriptor instead.
func (*Alert) Descriptor() ([]byte, []int) {
	return file_Alerts_proto_rawDescGZIP(), []int{0}
}

func (x *Alert) GetType() AlertType {
	if x != nil {
		return x.Type
	}
	return AlertType_UNKNOWN
}

func (x *Alert) GetSeverity() AlertSeverity {
	if x != nil {
		return x.Severity
	}
	return AlertSeverity_INFO
}

func (x *Alert) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Alert) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *Alert) GetDatetime() int64 {
	if x != nil {
		return x.Datetime
	}
	return 0
}

func (x *Alert) GetIsLoggerAlert() bool {
	if x != nil {
		return x.IsLoggerAlert
	}
	return false
}

var File_Alerts_proto protoreflect.FileDescriptor

var file_Alerts_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x22, 0xab, 0x02, 0x0a, 0x05, 0x41, 0x6c, 0x65, 0x72, 0x74,
	0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11,
	0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72,
	0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x61, 0x6c, 0x65, 0x72,
	0x74, 0x73, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79,
	0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x6c,
	0x65, 0x72, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x65, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x61, 0x74, 0x65, 0x74,
	0x69, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x69, 0x73, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x41,
	0x6c, 0x65, 0x72, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x73, 0x4c, 0x6f,
	0x67, 0x67, 0x65, 0x72, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x2a, 0xac, 0x02, 0x0a, 0x09, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x08, 0x0a, 0x04, 0x55, 0x53, 0x45, 0x52, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x49, 0x4e,
	0x4b, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x48, 0x52, 0x45, 0x41, 0x54, 0x50, 0x52, 0x45,
	0x56, 0x45, 0x4e, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x57, 0x45, 0x42,
	0x46, 0x49, 0x4c, 0x54, 0x45, 0x52, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x57, 0x45, 0x42, 0x43,
	0x4c, 0x41, 0x53, 0x53, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x05, 0x12,
	0x09, 0x0a, 0x05, 0x47, 0x45, 0x4f, 0x49, 0x50, 0x10, 0x06, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x45,
	0x54, 0x54, 0x49, 0x4e, 0x47, 0x53, 0x10, 0x07, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x49, 0x53, 0x43,
	0x4f, 0x56, 0x45, 0x52, 0x59, 0x10, 0x08, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x48, 0x43, 0x50, 0x10,
	0x09, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0x0a, 0x12, 0x07, 0x0a, 0x03, 0x56, 0x50, 0x4e, 0x10, 0x0b, 0x12, 0x11, 0x0a,
	0x0d, 0x43, 0x41, 0x50, 0x54, 0x49, 0x56, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x41, 0x4c, 0x10, 0x0c,
	0x12, 0x11, 0x0a, 0x0d, 0x46, 0x49, 0x52, 0x45, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x56, 0x45, 0x4e,
	0x54, 0x10, 0x0d, 0x12, 0x10, 0x0a, 0x0c, 0x44, 0x59, 0x4e, 0x41, 0x4d, 0x49, 0x43, 0x4c, 0x49,
	0x53, 0x54, 0x53, 0x10, 0x0e, 0x12, 0x11, 0x0a, 0x0d, 0x50, 0x4f, 0x4c, 0x49, 0x43, 0x59, 0x4d,
	0x41, 0x4e, 0x41, 0x47, 0x45, 0x52, 0x10, 0x0f, 0x12, 0x13, 0x0a, 0x0f, 0x44, 0x41, 0x54, 0x41,
	0x42, 0x41, 0x53, 0x45, 0x4d, 0x41, 0x4e, 0x41, 0x47, 0x45, 0x52, 0x10, 0x10, 0x12, 0x14, 0x0a,
	0x10, 0x44, 0x4e, 0x53, 0x46, 0x49, 0x4c, 0x54, 0x45, 0x52, 0x4d, 0x41, 0x4e, 0x41, 0x47, 0x45,
	0x52, 0x10, 0x11, 0x2a, 0x47, 0x0a, 0x0d, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x53, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x57, 0x41, 0x52, 0x4e, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x45, 0x42, 0x55, 0x47, 0x10, 0x03, 0x12, 0x0c,
	0x0a, 0x08, 0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x10, 0x04, 0x42, 0x42, 0x5a, 0x40,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e,
	0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Alerts_proto_rawDescOnce sync.Once
	file_Alerts_proto_rawDescData = file_Alerts_proto_rawDesc
)

func file_Alerts_proto_rawDescGZIP() []byte {
	file_Alerts_proto_rawDescOnce.Do(func() {
		file_Alerts_proto_rawDescData = protoimpl.X.CompressGZIP(file_Alerts_proto_rawDescData)
	})
	return file_Alerts_proto_rawDescData
}

var file_Alerts_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_Alerts_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_Alerts_proto_goTypes = []interface{}{
	(AlertType)(0),     // 0: alerts.AlertType
	(AlertSeverity)(0), // 1: alerts.AlertSeverity
	(*Alert)(nil),      // 2: alerts.Alert
	nil,                // 3: alerts.Alert.ParamsEntry
}
var file_Alerts_proto_depIdxs = []int32{
	0, // 0: alerts.Alert.type:type_name -> alerts.AlertType
	1, // 1: alerts.Alert.severity:type_name -> alerts.AlertSeverity
	3, // 2: alerts.Alert.params:type_name -> alerts.Alert.ParamsEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_Alerts_proto_init() }
func file_Alerts_proto_init() {
	if File_Alerts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Alerts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Alert); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_Alerts_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Alerts_proto_goTypes,
		DependencyIndexes: file_Alerts_proto_depIdxs,
		EnumInfos:         file_Alerts_proto_enumTypes,
		MessageInfos:      file_Alerts_proto_msgTypes,
	}.Build()
	File_Alerts_proto = out.File
	file_Alerts_proto_rawDesc = nil
	file_Alerts_proto_goTypes = nil
	file_Alerts_proto_depIdxs = nil
}

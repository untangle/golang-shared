// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.24.4
// source: DosEvent.proto

package DosEvent

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

type Protocol int32

const (
	Protocol_ALL  Protocol = 0 // Represents 'all'
	Protocol_TCP  Protocol = 1 // Represents 'tcp'
	Protocol_UDP  Protocol = 2 // Represents 'udp'
	Protocol_ICMP Protocol = 3 // Represents 'icmp'
)

// Enum value maps for Protocol.
var (
	Protocol_name = map[int32]string{
		0: "ALL",
		1: "TCP",
		2: "UDP",
		3: "ICMP",
	}
	Protocol_value = map[string]int32{
		"ALL":  0,
		"TCP":  1,
		"UDP":  2,
		"ICMP": 3,
	}
)

func (x Protocol) Enum() *Protocol {
	p := new(Protocol)
	*p = x
	return p
}

func (x Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_DosEvent_proto_enumTypes[0].Descriptor()
}

func (Protocol) Type() protoreflect.EnumType {
	return &file_DosEvent_proto_enumTypes[0]
}

func (x Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Protocol.Descriptor instead.
func (Protocol) EnumDescriptor() ([]byte, []int) {
	return file_DosEvent_proto_rawDescGZIP(), []int{0}
}

type DosEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromHost  string   `protobuf:"bytes,1,opt,name=fromHost,proto3" json:"fromHost,omitempty"`                         // Source IP address as string (e.g., 192.168.56.11)
	ToHost    string   `protobuf:"bytes,2,opt,name=toHost,proto3" json:"toHost,omitempty"`                             // Destination IP address as string (e.g., 192.168.56.11)
	Protocol  Protocol `protobuf:"varint,3,opt,name=protocol,proto3,enum=DosEvent.Protocol" json:"protocol,omitempty"` // Enum for protocol type (all, tcp, udp, icmp)
	RuleId    string   `protobuf:"bytes,4,opt,name=ruleId,proto3" json:"ruleId,omitempty"`                             // RuleID for the DOS rule
	TimeStamp int64    `protobuf:"varint,5,opt,name=timeStamp,proto3" json:"timeStamp,omitempty"`                      // Unix timestamp in milliseconds
}

func (x *DosEvent) Reset() {
	*x = DosEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_DosEvent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DosEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DosEvent) ProtoMessage() {}

func (x *DosEvent) ProtoReflect() protoreflect.Message {
	mi := &file_DosEvent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DosEvent.ProtoReflect.Descriptor instead.
func (*DosEvent) Descriptor() ([]byte, []int) {
	return file_DosEvent_proto_rawDescGZIP(), []int{0}
}

func (x *DosEvent) GetFromHost() string {
	if x != nil {
		return x.FromHost
	}
	return ""
}

func (x *DosEvent) GetToHost() string {
	if x != nil {
		return x.ToHost
	}
	return ""
}

func (x *DosEvent) GetProtocol() Protocol {
	if x != nil {
		return x.Protocol
	}
	return Protocol_ALL
}

func (x *DosEvent) GetRuleId() string {
	if x != nil {
		return x.RuleId
	}
	return ""
}

func (x *DosEvent) GetTimeStamp() int64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

var File_DosEvent_proto protoreflect.FileDescriptor

var file_DosEvent_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x44, 0x6f, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x44, 0x6f, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0xa4, 0x01, 0x0a, 0x08, 0x44,
	0x6f, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x48,
	0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x48,
	0x6f, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x48, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e,
	0x44, 0x6f, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x75, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x75, 0x6c,
	0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d,
	0x70, 0x2a, 0x2f, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x07, 0x0a,
	0x03, 0x41, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x43, 0x50, 0x10, 0x01, 0x12,
	0x07, 0x0a, 0x03, 0x55, 0x44, 0x50, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x43, 0x4d, 0x50,
	0x10, 0x03, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67,
	0x2d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f,
	0x44, 0x6f, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_DosEvent_proto_rawDescOnce sync.Once
	file_DosEvent_proto_rawDescData = file_DosEvent_proto_rawDesc
)

func file_DosEvent_proto_rawDescGZIP() []byte {
	file_DosEvent_proto_rawDescOnce.Do(func() {
		file_DosEvent_proto_rawDescData = protoimpl.X.CompressGZIP(file_DosEvent_proto_rawDescData)
	})
	return file_DosEvent_proto_rawDescData
}

var file_DosEvent_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_DosEvent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_DosEvent_proto_goTypes = []interface{}{
	(Protocol)(0),    // 0: DosEvent.Protocol
	(*DosEvent)(nil), // 1: DosEvent.DosEvent
}
var file_DosEvent_proto_depIdxs = []int32{
	0, // 0: DosEvent.DosEvent.protocol:type_name -> DosEvent.Protocol
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_DosEvent_proto_init() }
func file_DosEvent_proto_init() {
	if File_DosEvent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_DosEvent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DosEvent); i {
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
			RawDescriptor: file_DosEvent_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_DosEvent_proto_goTypes,
		DependencyIndexes: file_DosEvent_proto_depIdxs,
		EnumInfos:         file_DosEvent_proto_enumTypes,
		MessageInfos:      file_DosEvent_proto_msgTypes,
	}.Build()
	File_DosEvent_proto = out.File
	file_DosEvent_proto_rawDesc = nil
	file_DosEvent_proto_goTypes = nil
	file_DosEvent_proto_depIdxs = nil
}

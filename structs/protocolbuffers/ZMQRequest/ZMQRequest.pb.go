// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.21.12
// source: ZMQRequest.proto

package ZMQRequest

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

type ZMQRequest_Service int32

const (
	ZMQRequest_PACKETD   ZMQRequest_Service = 0
	ZMQRequest_REPORTD   ZMQRequest_Service = 1
	ZMQRequest_DISCOVERD ZMQRequest_Service = 2
)

// Enum value maps for ZMQRequest_Service.
var (
	ZMQRequest_Service_name = map[int32]string{
		0: "PACKETD",
		1: "REPORTD",
		2: "DISCOVERD",
	}
	ZMQRequest_Service_value = map[string]int32{
		"PACKETD":   0,
		"REPORTD":   1,
		"DISCOVERD": 2,
	}
)

func (x ZMQRequest_Service) Enum() *ZMQRequest_Service {
	p := new(ZMQRequest_Service)
	*p = x
	return p
}

func (x ZMQRequest_Service) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ZMQRequest_Service) Descriptor() protoreflect.EnumDescriptor {
	return file_ZMQRequest_proto_enumTypes[0].Descriptor()
}

func (ZMQRequest_Service) Type() protoreflect.EnumType {
	return &file_ZMQRequest_proto_enumTypes[0]
}

func (x ZMQRequest_Service) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ZMQRequest_Service.Descriptor instead.
func (ZMQRequest_Service) EnumDescriptor() ([]byte, []int) {
	return file_ZMQRequest_proto_rawDescGZIP(), []int{0, 0}
}

type ZMQRequest_Function int32

const (
	ZMQRequest_TEST_INFO    ZMQRequest_Function = 0
	ZMQRequest_GET_SESSIONS ZMQRequest_Function = 1
	ZMQRequest_GET_DEVICES  ZMQRequest_Function = 5 // Request all known devices from discoverd.
)

// Enum value maps for ZMQRequest_Function.
var (
	ZMQRequest_Function_name = map[int32]string{
		0: "TEST_INFO",
		1: "GET_SESSIONS",
		5: "GET_DEVICES",
	}
	ZMQRequest_Function_value = map[string]int32{
		"TEST_INFO":    0,
		"GET_SESSIONS": 1,
		"GET_DEVICES":  5,
	}
)

func (x ZMQRequest_Function) Enum() *ZMQRequest_Function {
	p := new(ZMQRequest_Function)
	*p = x
	return p
}

func (x ZMQRequest_Function) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ZMQRequest_Function) Descriptor() protoreflect.EnumDescriptor {
	return file_ZMQRequest_proto_enumTypes[1].Descriptor()
}

func (ZMQRequest_Function) Type() protoreflect.EnumType {
	return &file_ZMQRequest_proto_enumTypes[1]
}

func (x ZMQRequest_Function) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ZMQRequest_Function.Descriptor instead.
func (ZMQRequest_Function) EnumDescriptor() ([]byte, []int) {
	return file_ZMQRequest_proto_rawDescGZIP(), []int{0, 1}
}

type ZMQRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service  ZMQRequest_Service  `protobuf:"varint,1,opt,name=service,proto3,enum=reports.ZMQRequest_Service" json:"service,omitempty"`
	Function ZMQRequest_Function `protobuf:"varint,2,opt,name=function,proto3,enum=reports.ZMQRequest_Function" json:"function,omitempty"`
	Data     string              `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ZMQRequest) Reset() {
	*x = ZMQRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ZMQRequest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZMQRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZMQRequest) ProtoMessage() {}

func (x *ZMQRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ZMQRequest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZMQRequest.ProtoReflect.Descriptor instead.
func (*ZMQRequest) Descriptor() ([]byte, []int) {
	return file_ZMQRequest_proto_rawDescGZIP(), []int{0}
}

func (x *ZMQRequest) GetService() ZMQRequest_Service {
	if x != nil {
		return x.Service
	}
	return ZMQRequest_PACKETD
}

func (x *ZMQRequest) GetFunction() ZMQRequest_Function {
	if x != nil {
		return x.Function
	}
	return ZMQRequest_TEST_INFO
}

func (x *ZMQRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_ZMQRequest_proto protoreflect.FileDescriptor

var file_ZMQRequest_proto_rawDesc = []byte{
	0x0a, 0x10, 0x5a, 0x4d, 0x51, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x22, 0x83, 0x02, 0x0a, 0x0a,
	0x5a, 0x4d, 0x51, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x07, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x73, 0x2e, 0x5a, 0x4d, 0x51, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x38, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2e, 0x5a, 0x4d,
	0x51, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x32, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x41,
	0x43, 0x4b, 0x45, 0x54, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x50, 0x4f, 0x52,
	0x54, 0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x56, 0x45, 0x52,
	0x44, 0x10, 0x02, 0x22, 0x3c, 0x0a, 0x08, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0d, 0x0a, 0x09, 0x54, 0x45, 0x53, 0x54, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x53, 0x10, 0x01,
	0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x45, 0x54, 0x5f, 0x44, 0x45, 0x56, 0x49, 0x43, 0x45, 0x53, 0x10,
	0x05, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x5a,
	0x4d, 0x51, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_ZMQRequest_proto_rawDescOnce sync.Once
	file_ZMQRequest_proto_rawDescData = file_ZMQRequest_proto_rawDesc
)

func file_ZMQRequest_proto_rawDescGZIP() []byte {
	file_ZMQRequest_proto_rawDescOnce.Do(func() {
		file_ZMQRequest_proto_rawDescData = protoimpl.X.CompressGZIP(file_ZMQRequest_proto_rawDescData)
	})
	return file_ZMQRequest_proto_rawDescData
}

var file_ZMQRequest_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_ZMQRequest_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ZMQRequest_proto_goTypes = []interface{}{
	(ZMQRequest_Service)(0),  // 0: reports.ZMQRequest.Service
	(ZMQRequest_Function)(0), // 1: reports.ZMQRequest.Function
	(*ZMQRequest)(nil),       // 2: reports.ZMQRequest
}
var file_ZMQRequest_proto_depIdxs = []int32{
	0, // 0: reports.ZMQRequest.service:type_name -> reports.ZMQRequest.Service
	1, // 1: reports.ZMQRequest.function:type_name -> reports.ZMQRequest.Function
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_ZMQRequest_proto_init() }
func file_ZMQRequest_proto_init() {
	if File_ZMQRequest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ZMQRequest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZMQRequest); i {
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
			RawDescriptor: file_ZMQRequest_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ZMQRequest_proto_goTypes,
		DependencyIndexes: file_ZMQRequest_proto_depIdxs,
		EnumInfos:         file_ZMQRequest_proto_enumTypes,
		MessageInfos:      file_ZMQRequest_proto_msgTypes,
	}.Build()
	File_ZMQRequest_proto = out.File
	file_ZMQRequest_proto_rawDesc = nil
	file_ZMQRequest_proto_goTypes = nil
	file_ZMQRequest_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0-devel
// 	protoc        v3.6.1
// source: SessionEvent.proto

package SessionEvent

import (
	_struct "github.com/golang/protobuf/ptypes/struct"
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

type SessionEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            string          `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Table           string          `protobuf:"bytes,2,opt,name=Table,proto3" json:"Table,omitempty"`
	SQLOp           int32           `protobuf:"varint,3,opt,name=SQLOp,proto3" json:"SQLOp,omitempty"`
	Columns         *_struct.Struct `protobuf:"bytes,4,opt,name=Columns,proto3" json:"Columns,omitempty"`
	ModifiedColumns *_struct.Struct `protobuf:"bytes,5,opt,name=ModifiedColumns,proto3" json:"ModifiedColumns,omitempty"`
}

func (x *SessionEvent) Reset() {
	*x = SessionEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SessionEvent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionEvent) ProtoMessage() {}

func (x *SessionEvent) ProtoReflect() protoreflect.Message {
	mi := &file_SessionEvent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionEvent.ProtoReflect.Descriptor instead.
func (*SessionEvent) Descriptor() ([]byte, []int) {
	return file_SessionEvent_proto_rawDescGZIP(), []int{0}
}

func (x *SessionEvent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SessionEvent) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *SessionEvent) GetSQLOp() int32 {
	if x != nil {
		return x.SQLOp
	}
	return 0
}

func (x *SessionEvent) GetColumns() *_struct.Struct {
	if x != nil {
		return x.Columns
	}
	return nil
}

func (x *SessionEvent) GetModifiedColumns() *_struct.Struct {
	if x != nil {
		return x.ModifiedColumns
	}
	return nil
}

var File_SessionEvent_proto protoreflect.FileDescriptor

var file_SessionEvent_proto_rawDesc = []byte{
	0x0a, 0x12, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc4, 0x01, 0x0a, 0x0c,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x51, 0x4c, 0x4f, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x53, 0x51, 0x4c, 0x4f, 0x70, 0x12, 0x31, 0x0a, 0x07,
	0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x12,
	0x41, 0x0a, 0x0f, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x43, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x52, 0x0f, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x43, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x73, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67,
	0x2d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_SessionEvent_proto_rawDescOnce sync.Once
	file_SessionEvent_proto_rawDescData = file_SessionEvent_proto_rawDesc
)

func file_SessionEvent_proto_rawDescGZIP() []byte {
	file_SessionEvent_proto_rawDescOnce.Do(func() {
		file_SessionEvent_proto_rawDescData = protoimpl.X.CompressGZIP(file_SessionEvent_proto_rawDescData)
	})
	return file_SessionEvent_proto_rawDescData
}

var file_SessionEvent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_SessionEvent_proto_goTypes = []interface{}{
	(*SessionEvent)(nil),   // 0: reports.SessionEvent
	(*_struct.Struct)(nil), // 1: google.protobuf.Struct
}
var file_SessionEvent_proto_depIdxs = []int32{
	1, // 0: reports.SessionEvent.Columns:type_name -> google.protobuf.Struct
	1, // 1: reports.SessionEvent.ModifiedColumns:type_name -> google.protobuf.Struct
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_SessionEvent_proto_init() }
func file_SessionEvent_proto_init() {
	if File_SessionEvent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_SessionEvent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionEvent); i {
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
			RawDescriptor: file_SessionEvent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_SessionEvent_proto_goTypes,
		DependencyIndexes: file_SessionEvent_proto_depIdxs,
		MessageInfos:      file_SessionEvent_proto_msgTypes,
	}.Build()
	File_SessionEvent_proto = out.File
	file_SessionEvent_proto_rawDesc = nil
	file_SessionEvent_proto_goTypes = nil
	file_SessionEvent_proto_depIdxs = nil
}

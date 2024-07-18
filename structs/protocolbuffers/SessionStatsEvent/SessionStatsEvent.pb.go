// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: SessionStatsEvent.proto

package SessionStatsEvent

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

type SessionStatsEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionID        int64  `protobuf:"varint,1,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	TimeStamp        int64  `protobuf:"varint,2,opt,name=timeStamp,proto3" json:"timeStamp,omitempty"`
	Bytes            uint64 `protobuf:"varint,3,opt,name=bytes,proto3" json:"bytes,omitempty"`
	ClientBytes      uint64 `protobuf:"varint,4,opt,name=clientBytes,proto3" json:"clientBytes,omitempty"`
	ServerBytes      uint64 `protobuf:"varint,5,opt,name=serverBytes,proto3" json:"serverBytes,omitempty"`
	ByteRate         uint32 `protobuf:"varint,6,opt,name=byteRate,proto3" json:"byteRate,omitempty"`
	ClientByteRate   uint32 `protobuf:"varint,7,opt,name=clientByteRate,proto3" json:"clientByteRate,omitempty"`
	ServerByteRate   uint32 `protobuf:"varint,8,opt,name=serverByteRate,proto3" json:"serverByteRate,omitempty"`
	Packets          uint64 `protobuf:"varint,9,opt,name=packets,proto3" json:"packets,omitempty"`
	ClientPackets    uint64 `protobuf:"varint,10,opt,name=clientPackets,proto3" json:"clientPackets,omitempty"`
	ServerPackets    uint64 `protobuf:"varint,11,opt,name=serverPackets,proto3" json:"serverPackets,omitempty"`
	PacketRate       uint32 `protobuf:"varint,12,opt,name=packetRate,proto3" json:"packetRate,omitempty"`
	ClientPacketRate uint32 `protobuf:"varint,13,opt,name=clientPacketRate,proto3" json:"clientPacketRate,omitempty"`
	ServerPacketRate uint32 `protobuf:"varint,14,opt,name=serverPacketRate,proto3" json:"serverPacketRate,omitempty"`
	// Network address of the client
	ClientNetworkAddress string `protobuf:"bytes,15,opt,name=clientNetworkAddress,proto3" json:"clientNetworkAddress,omitempty"`
}

func (x *SessionStatsEvent) Reset() {
	*x = SessionStatsEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SessionStatsEvent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionStatsEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionStatsEvent) ProtoMessage() {}

func (x *SessionStatsEvent) ProtoReflect() protoreflect.Message {
	mi := &file_SessionStatsEvent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionStatsEvent.ProtoReflect.Descriptor instead.
func (*SessionStatsEvent) Descriptor() ([]byte, []int) {
	return file_SessionStatsEvent_proto_rawDescGZIP(), []int{0}
}

func (x *SessionStatsEvent) GetSessionID() int64 {
	if x != nil {
		return x.SessionID
	}
	return 0
}

func (x *SessionStatsEvent) GetTimeStamp() int64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *SessionStatsEvent) GetBytes() uint64 {
	if x != nil {
		return x.Bytes
	}
	return 0
}

func (x *SessionStatsEvent) GetClientBytes() uint64 {
	if x != nil {
		return x.ClientBytes
	}
	return 0
}

func (x *SessionStatsEvent) GetServerBytes() uint64 {
	if x != nil {
		return x.ServerBytes
	}
	return 0
}

func (x *SessionStatsEvent) GetByteRate() uint32 {
	if x != nil {
		return x.ByteRate
	}
	return 0
}

func (x *SessionStatsEvent) GetClientByteRate() uint32 {
	if x != nil {
		return x.ClientByteRate
	}
	return 0
}

func (x *SessionStatsEvent) GetServerByteRate() uint32 {
	if x != nil {
		return x.ServerByteRate
	}
	return 0
}

func (x *SessionStatsEvent) GetPackets() uint64 {
	if x != nil {
		return x.Packets
	}
	return 0
}

func (x *SessionStatsEvent) GetClientPackets() uint64 {
	if x != nil {
		return x.ClientPackets
	}
	return 0
}

func (x *SessionStatsEvent) GetServerPackets() uint64 {
	if x != nil {
		return x.ServerPackets
	}
	return 0
}

func (x *SessionStatsEvent) GetPacketRate() uint32 {
	if x != nil {
		return x.PacketRate
	}
	return 0
}

func (x *SessionStatsEvent) GetClientPacketRate() uint32 {
	if x != nil {
		return x.ClientPacketRate
	}
	return 0
}

func (x *SessionStatsEvent) GetServerPacketRate() uint32 {
	if x != nil {
		return x.ServerPacketRate
	}
	return 0
}

func (x *SessionStatsEvent) GetClientNetworkAddress() string {
	if x != nil {
		return x.ClientNetworkAddress
	}
	return ""
}

var File_SessionStatsEvent_proto protoreflect.FileDescriptor

var file_SessionStatsEvent_proto_rawDesc = []byte{
	0x0a, 0x17, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x73, 0x22, 0xa7, 0x04, 0x0a, 0x11, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x62, 0x79, 0x74, 0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x08, 0x62, 0x79, 0x74, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x74, 0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x74, 0x65, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x42, 0x79, 0x74, 0x65,
	0x52, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x42, 0x79, 0x74, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x70, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x2a, 0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x61, 0x74, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65, 0x12, 0x2a, 0x0a, 0x10,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x14, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x4d, 0x5a, 0x4b,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e,
	0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_SessionStatsEvent_proto_rawDescOnce sync.Once
	file_SessionStatsEvent_proto_rawDescData = file_SessionStatsEvent_proto_rawDesc
)

func file_SessionStatsEvent_proto_rawDescGZIP() []byte {
	file_SessionStatsEvent_proto_rawDescOnce.Do(func() {
		file_SessionStatsEvent_proto_rawDescData = protoimpl.X.CompressGZIP(file_SessionStatsEvent_proto_rawDescData)
	})
	return file_SessionStatsEvent_proto_rawDescData
}

var file_SessionStatsEvent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_SessionStatsEvent_proto_goTypes = []any{
	(*SessionStatsEvent)(nil), // 0: reports.SessionStatsEvent
}
var file_SessionStatsEvent_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_SessionStatsEvent_proto_init() }
func file_SessionStatsEvent_proto_init() {
	if File_SessionStatsEvent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_SessionStatsEvent_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SessionStatsEvent); i {
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
			RawDescriptor: file_SessionStatsEvent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_SessionStatsEvent_proto_goTypes,
		DependencyIndexes: file_SessionStatsEvent_proto_depIdxs,
		MessageInfos:      file_SessionStatsEvent_proto_msgTypes,
	}.Build()
	File_SessionStatsEvent_proto = out.File
	file_SessionStatsEvent_proto_rawDesc = nil
	file_SessionStatsEvent_proto_goTypes = nil
	file_SessionStatsEvent_proto_depIdxs = nil
}

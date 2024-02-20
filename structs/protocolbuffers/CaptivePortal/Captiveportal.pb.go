// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: Captiveportal.proto

package CaptivePortal

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

type UserGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientIp string `protobuf:"bytes,1,opt,name=ClientIp,proto3" json:"ClientIp,omitempty"`
}

func (x *UserGetRequest) Reset() {
	*x = UserGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGetRequest) ProtoMessage() {}

func (x *UserGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Captiveportal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserGetRequest.ProtoReflect.Descriptor instead.
func (*UserGetRequest) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{0}
}

func (x *UserGetRequest) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

// Get user info to redirect to respective captive portal.
type UserGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientIp string `protobuf:"bytes,1,opt,name=ClientIp,proto3" json:"ClientIp,omitempty"`
	ConfigId string `protobuf:"bytes,2,opt,name=ConfigId,proto3" json:"ConfigId,omitempty"`
}

func (x *UserGetResponse) Reset() {
	*x = UserGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGetResponse) ProtoMessage() {}

func (x *UserGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Captiveportal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserGetResponse.ProtoReflect.Descriptor instead.
func (*UserGetResponse) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{1}
}

func (x *UserGetResponse) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

func (x *UserGetResponse) GetConfigId() string {
	if x != nil {
		return x.ConfigId
	}
	return ""
}

type UserSetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientIp string `protobuf:"bytes,1,opt,name=ClientIp,proto3" json:"ClientIp,omitempty"`
}

func (x *UserSetRequest) Reset() {
	*x = UserSetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetRequest) ProtoMessage() {}

func (x *UserSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Captiveportal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSetRequest.ProtoReflect.Descriptor instead.
func (*UserSetRequest) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{2}
}

func (x *UserSetRequest) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

// Update captive portal t&c accepted status.
type UserSetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Done bool `protobuf:"varint,1,opt,name=Done,proto3" json:"Done,omitempty"`
}

func (x *UserSetResponse) Reset() {
	*x = UserSetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetResponse) ProtoMessage() {}

func (x *UserSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Captiveportal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSetResponse.ProtoReflect.Descriptor instead.
func (*UserSetResponse) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{3}
}

func (x *UserSetResponse) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

type CpUserEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigId              string `protobuf:"bytes,1,opt,name=ConfigId,proto3" json:"ConfigId,omitempty"`
	TimeoutDuration       int64  `protobuf:"varint,2,opt,name=TimeoutDuration,proto3" json:"TimeoutDuration,omitempty"`
	LastAcceptedTimeStamp int64  `protobuf:"varint,3,opt,name=LastAcceptedTimeStamp,proto3" json:"LastAcceptedTimeStamp,omitempty"`
	LastSeenTimeStamp     int64  `protobuf:"varint,4,opt,name=LastSeenTimeStamp,proto3" json:"LastSeenTimeStamp,omitempty"`
	Description           string `protobuf:"bytes,5,opt,name=Description,proto3" json:"Description,omitempty"`
	Host                  string `protobuf:"bytes,6,opt,name=Host,proto3" json:"Host,omitempty"`
}

func (x *CpUserEntry) Reset() {
	*x = CpUserEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpUserEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpUserEntry) ProtoMessage() {}

func (x *CpUserEntry) ProtoReflect() protoreflect.Message {
	mi := &file_Captiveportal_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CpUserEntry.ProtoReflect.Descriptor instead.
func (*CpUserEntry) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{4}
}

func (x *CpUserEntry) GetConfigId() string {
	if x != nil {
		return x.ConfigId
	}
	return ""
}

func (x *CpUserEntry) GetTimeoutDuration() int64 {
	if x != nil {
		return x.TimeoutDuration
	}
	return 0
}

func (x *CpUserEntry) GetLastAcceptedTimeStamp() int64 {
	if x != nil {
		return x.LastAcceptedTimeStamp
	}
	return 0
}

func (x *CpUserEntry) GetLastSeenTimeStamp() int64 {
	if x != nil {
		return x.LastSeenTimeStamp
	}
	return 0
}

func (x *CpUserEntry) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CpUserEntry) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

var File_Captiveportal_proto protoreflect.FileDescriptor

var file_Captiveportal_proto_rawDesc = []byte{
	0x0a, 0x13, 0x43, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x22, 0x2c, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x70, 0x22, 0x49, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x70, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x22, 0x2c, 0x0a,
	0x0e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x22, 0x25, 0x0a, 0x0f, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x44, 0x6f, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x44, 0x6f,
	0x6e, 0x65, 0x22, 0xed, 0x01, 0x0a, 0x0b, 0x43, 0x70, 0x55, 0x73, 0x65, 0x72, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x12, 0x28,
	0x0a, 0x0f, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x15, 0x4c, 0x61, 0x73, 0x74,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x15, 0x4c, 0x61, 0x73, 0x74, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x2c,
	0x0a, 0x11, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x4c, 0x61, 0x73, 0x74, 0x53,
	0x65, 0x65, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x20, 0x0a, 0x0b,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x6f,
	0x73, 0x74, 0x32, 0xcc, 0x01, 0x0a, 0x18, 0x43, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x47, 0x72, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x57, 0x0a, 0x14, 0x67, 0x65, 0x74, 0x43, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x14, 0x73, 0x65, 0x74, 0x43,
	0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x1d, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1e, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x43,
	0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Captiveportal_proto_rawDescOnce sync.Once
	file_Captiveportal_proto_rawDescData = file_Captiveportal_proto_rawDesc
)

func file_Captiveportal_proto_rawDescGZIP() []byte {
	file_Captiveportal_proto_rawDescOnce.Do(func() {
		file_Captiveportal_proto_rawDescData = protoimpl.X.CompressGZIP(file_Captiveportal_proto_rawDescData)
	})
	return file_Captiveportal_proto_rawDescData
}

var file_Captiveportal_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_Captiveportal_proto_goTypes = []interface{}{
	(*UserGetRequest)(nil),  // 0: captiveportal.UserGetRequest
	(*UserGetResponse)(nil), // 1: captiveportal.UserGetResponse
	(*UserSetRequest)(nil),  // 2: captiveportal.UserSetRequest
	(*UserSetResponse)(nil), // 3: captiveportal.UserSetResponse
	(*CpUserEntry)(nil),     // 4: captiveportal.CpUserEntry
}
var file_Captiveportal_proto_depIdxs = []int32{
	0, // 0: captiveportal.CaptivePortalGrpcService.getCaptivePortalUser:input_type -> captiveportal.UserGetRequest
	2, // 1: captiveportal.CaptivePortalGrpcService.setCaptivePortalUser:input_type -> captiveportal.UserSetRequest
	1, // 2: captiveportal.CaptivePortalGrpcService.getCaptivePortalUser:output_type -> captiveportal.UserGetResponse
	3, // 3: captiveportal.CaptivePortalGrpcService.setCaptivePortalUser:output_type -> captiveportal.UserSetResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_Captiveportal_proto_init() }
func file_Captiveportal_proto_init() {
	if File_Captiveportal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Captiveportal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserGetRequest); i {
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
		file_Captiveportal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserGetResponse); i {
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
		file_Captiveportal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserSetRequest); i {
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
		file_Captiveportal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserSetResponse); i {
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
		file_Captiveportal_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CpUserEntry); i {
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
			RawDescriptor: file_Captiveportal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_Captiveportal_proto_goTypes,
		DependencyIndexes: file_Captiveportal_proto_depIdxs,
		MessageInfos:      file_Captiveportal_proto_msgTypes,
	}.Build()
	File_Captiveportal_proto = out.File
	file_Captiveportal_proto_rawDesc = nil
	file_Captiveportal_proto_goTypes = nil
	file_Captiveportal_proto_depIdxs = nil
}

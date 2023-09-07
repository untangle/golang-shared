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

type ImageDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageData string `protobuf:"bytes,1,opt,name=imageData,proto3" json:"imageData,omitempty"`
	ImageName string `protobuf:"bytes,2,opt,name=imageName,proto3" json:"imageName,omitempty"`
}

func (x *ImageDetails) Reset() {
	*x = ImageDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageDetails) ProtoMessage() {}

func (x *ImageDetails) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ImageDetails.ProtoReflect.Descriptor instead.
func (*ImageDetails) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{0}
}

func (x *ImageDetails) GetImageData() string {
	if x != nil {
		return x.ImageData
	}
	return ""
}

func (x *ImageDetails) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

type CpSettingType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enabled          bool          `protobuf:"varint,1,opt,name=Enabled,json=enabled,proto3" json:"Enabled,omitempty"`
	TimeoutValue     int32         `protobuf:"varint,2,opt,name=TimeoutValue,json=timeoutValue,proto3" json:"TimeoutValue,omitempty"`
	TimeoutPeriod    string        `protobuf:"bytes,3,opt,name=TimeoutPeriod,json=timeoutPeriod,proto3" json:"TimeoutPeriod,omitempty"`
	AcceptText       string        `protobuf:"bytes,4,opt,name=AcceptText,json=acceptText,proto3" json:"AcceptText,omitempty"`
	AcceptButtonText string        `protobuf:"bytes,5,opt,name=AcceptButtonText,json=acceptButtonText,proto3" json:"AcceptButtonText,omitempty"`
	MessageHeading   string        `protobuf:"bytes,6,opt,name=MessageHeading,json=messageHeading,proto3" json:"MessageHeading,omitempty"`
	MessageText      string        `protobuf:"bytes,7,opt,name=MessageText,json=messageText,proto3" json:"MessageText,omitempty"`
	WelcomeText      string        `protobuf:"bytes,8,opt,name=WelcomeText,json=welcomeText,proto3" json:"WelcomeText,omitempty"`
	PageTitle        string        `protobuf:"bytes,9,opt,name=PageTitle,json=pageTitle,proto3" json:"PageTitle,omitempty"`
	Logo             *ImageDetails `protobuf:"bytes,10,opt,name=logo,proto3" json:"logo,omitempty"`
}

func (x *CpSettingType) Reset() {
	*x = CpSettingType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpSettingType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpSettingType) ProtoMessage() {}

func (x *CpSettingType) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use CpSettingType.ProtoReflect.Descriptor instead.
func (*CpSettingType) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{1}
}

func (x *CpSettingType) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *CpSettingType) GetTimeoutValue() int32 {
	if x != nil {
		return x.TimeoutValue
	}
	return 0
}

func (x *CpSettingType) GetTimeoutPeriod() string {
	if x != nil {
		return x.TimeoutPeriod
	}
	return ""
}

func (x *CpSettingType) GetAcceptText() string {
	if x != nil {
		return x.AcceptText
	}
	return ""
}

func (x *CpSettingType) GetAcceptButtonText() string {
	if x != nil {
		return x.AcceptButtonText
	}
	return ""
}

func (x *CpSettingType) GetMessageHeading() string {
	if x != nil {
		return x.MessageHeading
	}
	return ""
}

func (x *CpSettingType) GetMessageText() string {
	if x != nil {
		return x.MessageText
	}
	return ""
}

func (x *CpSettingType) GetWelcomeText() string {
	if x != nil {
		return x.WelcomeText
	}
	return ""
}

func (x *CpSettingType) GetPageTitle() string {
	if x != nil {
		return x.PageTitle
	}
	return ""
}

func (x *CpSettingType) GetLogo() *ImageDetails {
	if x != nil {
		return x.Logo
	}
	return nil
}

type UserGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientIp string `protobuf:"bytes,1,opt,name=ClientIp,proto3" json:"ClientIp,omitempty"`
}

func (x *UserGetRequest) Reset() {
	*x = UserGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGetRequest) ProtoMessage() {}

func (x *UserGetRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use UserGetRequest.ProtoReflect.Descriptor instead.
func (*UserGetRequest) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{2}
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

	ClientIp     string         `protobuf:"bytes,1,opt,name=ClientIp,proto3" json:"ClientIp,omitempty"`
	PolicyId     string         `protobuf:"bytes,2,opt,name=PolicyId,proto3" json:"PolicyId,omitempty"`
	PolicyConfig *CpSettingType `protobuf:"bytes,3,opt,name=PolicyConfig,proto3" json:"PolicyConfig,omitempty"`
}

func (x *UserGetResponse) Reset() {
	*x = UserGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Captiveportal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGetResponse) ProtoMessage() {}

func (x *UserGetResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use UserGetResponse.ProtoReflect.Descriptor instead.
func (*UserGetResponse) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{3}
}

func (x *UserGetResponse) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

func (x *UserGetResponse) GetPolicyId() string {
	if x != nil {
		return x.PolicyId
	}
	return ""
}

func (x *UserGetResponse) GetPolicyConfig() *CpSettingType {
	if x != nil {
		return x.PolicyConfig
	}
	return nil
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
		mi := &file_Captiveportal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetRequest) ProtoMessage() {}

func (x *UserSetRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use UserSetRequest.ProtoReflect.Descriptor instead.
func (*UserSetRequest) Descriptor() ([]byte, []int) {
	return file_Captiveportal_proto_rawDescGZIP(), []int{4}
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
		mi := &file_Captiveportal_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetResponse) ProtoMessage() {}

func (x *UserSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Captiveportal_proto_msgTypes[5]
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
	return file_Captiveportal_proto_rawDescGZIP(), []int{5}
}

func (x *UserSetResponse) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

var File_Captiveportal_proto protoreflect.FileDescriptor

var file_Captiveportal_proto_rawDesc = []byte{
	0x0a, 0x13, 0x43, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x6c, 0x22, 0x4a, 0x0a, 0x0c, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0xfa, 0x02, 0x0a, 0x0d, 0x43, 0x70, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x22, 0x0a, 0x0c,
	0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x24, 0x0a, 0x0d, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x50, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x54, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x65,
	0x70, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x54, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x10, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x54, 0x65,
	0x78, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x65, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x48, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x65, 0x78, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x54, 0x65, 0x78, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x54, 0x65, 0x78, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x50, 0x61, 0x67, 0x65, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x2f, 0x0a, 0x04,
	0x6c, 0x6f, 0x67, 0x6f, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x61, 0x70,
	0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x22, 0x2c, 0x0a,
	0x0e, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x22, 0x8b, 0x01, 0x0a, 0x0f,
	0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x49, 0x64, 0x12, 0x40, 0x0a, 0x0c, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x43, 0x70,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x2c, 0x0a, 0x0e, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x22, 0x25, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x6f,
	0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x44, 0x6f, 0x6e, 0x65, 0x32, 0xcc,
	0x01, 0x0a, 0x18, 0x43, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c,
	0x47, 0x72, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x14, 0x67,
	0x65, 0x74, 0x43, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x61, 0x6c, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x14, 0x73, 0x65, 0x74, 0x43, 0x61, 0x70, 0x74, 0x69,
	0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x63,
	0x61, 0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x61,
	0x70, 0x74, 0x69, 0x76, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x49, 0x5a,
	0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x74, 0x61,
	0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x43, 0x61, 0x70, 0x74, 0x69,
	0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_Captiveportal_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_Captiveportal_proto_goTypes = []interface{}{
	(*ImageDetails)(nil),    // 0: captiveportal.ImageDetails
	(*CpSettingType)(nil),   // 1: captiveportal.CpSettingType
	(*UserGetRequest)(nil),  // 2: captiveportal.UserGetRequest
	(*UserGetResponse)(nil), // 3: captiveportal.UserGetResponse
	(*UserSetRequest)(nil),  // 4: captiveportal.UserSetRequest
	(*UserSetResponse)(nil), // 5: captiveportal.UserSetResponse
}
var file_Captiveportal_proto_depIdxs = []int32{
	0, // 0: captiveportal.CpSettingType.logo:type_name -> captiveportal.ImageDetails
	1, // 1: captiveportal.UserGetResponse.PolicyConfig:type_name -> captiveportal.CpSettingType
	2, // 2: captiveportal.CaptivePortalGrpcService.getCaptivePortalUser:input_type -> captiveportal.UserGetRequest
	4, // 3: captiveportal.CaptivePortalGrpcService.setCaptivePortalUser:input_type -> captiveportal.UserSetRequest
	3, // 4: captiveportal.CaptivePortalGrpcService.getCaptivePortalUser:output_type -> captiveportal.UserGetResponse
	5, // 5: captiveportal.CaptivePortalGrpcService.setCaptivePortalUser:output_type -> captiveportal.UserSetResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_Captiveportal_proto_init() }
func file_Captiveportal_proto_init() {
	if File_Captiveportal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Captiveportal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageDetails); i {
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
			switch v := v.(*CpSettingType); i {
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
		file_Captiveportal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_Captiveportal_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_Captiveportal_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_Captiveportal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
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

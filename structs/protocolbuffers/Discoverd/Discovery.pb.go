// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.6.1
// source: Discovery.proto

package Discoverd

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

type ResponseCode int32

const (
	ResponseCode_UNKNOWN         ResponseCode = 0
	ResponseCode_OK              ResponseCode = 1
	ResponseCode_ERROR           ResponseCode = 2
	ResponseCode_INVALID_REQUEST ResponseCode = 3
)

// Enum value maps for ResponseCode.
var (
	ResponseCode_name = map[int32]string{
		0: "UNKNOWN",
		1: "OK",
		2: "ERROR",
		3: "INVALID_REQUEST",
	}
	ResponseCode_value = map[string]int32{
		"UNKNOWN":         0,
		"OK":              1,
		"ERROR":           2,
		"INVALID_REQUEST": 3,
	}
)

func (x ResponseCode) Enum() *ResponseCode {
	p := new(ResponseCode)
	*p = x
	return p
}

func (x ResponseCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseCode) Descriptor() protoreflect.EnumDescriptor {
	return file_Discovery_proto_enumTypes[0].Descriptor()
}

func (ResponseCode) Type() protoreflect.EnumType {
	return &file_Discovery_proto_enumTypes[0]
}

func (x ResponseCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseCode.Descriptor instead.
func (ResponseCode) EnumDescriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{0}
}

type EmptyParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyParam) Reset() {
	*x = EmptyParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyParam) ProtoMessage() {}

func (x *EmptyParam) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyParam.ProtoReflect.Descriptor instead.
func (*EmptyParam) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{0}
}

type RequestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result ResponseCode `protobuf:"varint,1,opt,name=result,proto3,enum=discoverd.ResponseCode" json:"result,omitempty"`
}

func (x *RequestResponse) Reset() {
	*x = RequestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestResponse) ProtoMessage() {}

func (x *RequestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestResponse.ProtoReflect.Descriptor instead.
func (*RequestResponse) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{1}
}

func (x *RequestResponse) GetResult() ResponseCode {
	if x != nil {
		return x.Result
	}
	return ResponseCode_UNKNOWN
}

// param: net is a string in CIDR notation
type ScanNetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Net []string `protobuf:"bytes,1,rep,name=net,proto3" json:"net,omitempty"`
}

func (x *ScanNetRequest) Reset() {
	*x = ScanNetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanNetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanNetRequest) ProtoMessage() {}

func (x *ScanNetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanNetRequest.ProtoReflect.Descriptor instead.
func (*ScanNetRequest) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{2}
}

func (x *ScanNetRequest) GetNet() []string {
	if x != nil {
		return x.Net
	}
	return nil
}

// param: host is a string in host/ip notation
type ScanHostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host []string `protobuf:"bytes,1,rep,name=host,proto3" json:"host,omitempty"`
}

func (x *ScanHostRequest) Reset() {
	*x = ScanHostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanHostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanHostRequest) ProtoMessage() {}

func (x *ScanHostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanHostRequest.ProtoReflect.Descriptor instead.
func (*ScanHostRequest) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{3}
}

func (x *ScanHostRequest) GetHost() []string {
	if x != nil {
		return x.Host
	}
	return nil
}

type DiscoveryEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MacAddress  string `protobuf:"bytes,1,opt,name=macAddress,proto3" json:"macAddress,omitempty"`
	IPv4Address string `protobuf:"bytes,2,opt,name=IPv4Address,proto3" json:"IPv4Address,omitempty"`
	LastUpdate  int64  `protobuf:"varint,3,opt,name=LastUpdate,proto3" json:"LastUpdate,omitempty"`
	Lldp        *LLDP  `protobuf:"bytes,10,opt,name=lldp,proto3" json:"lldp,omitempty"`
	Arp         *ARP   `protobuf:"bytes,11,opt,name=arp,proto3" json:"arp,omitempty"`
	Nmap        *NMAP  `protobuf:"bytes,12,opt,name=nmap,proto3" json:"nmap,omitempty"`
}

func (x *DiscoveryEntry) Reset() {
	*x = DiscoveryEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryEntry) ProtoMessage() {}

func (x *DiscoveryEntry) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryEntry.ProtoReflect.Descriptor instead.
func (*DiscoveryEntry) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{4}
}

func (x *DiscoveryEntry) GetMacAddress() string {
	if x != nil {
		return x.MacAddress
	}
	return ""
}

func (x *DiscoveryEntry) GetIPv4Address() string {
	if x != nil {
		return x.IPv4Address
	}
	return ""
}

func (x *DiscoveryEntry) GetLastUpdate() int64 {
	if x != nil {
		return x.LastUpdate
	}
	return 0
}

func (x *DiscoveryEntry) GetLldp() *LLDP {
	if x != nil {
		return x.Lldp
	}
	return nil
}

func (x *DiscoveryEntry) GetArp() *ARP {
	if x != nil {
		return x.Arp
	}
	return nil
}

func (x *DiscoveryEntry) GetNmap() *NMAP {
	if x != nil {
		return x.Nmap
	}
	return nil
}

type LLDP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Chassis
	SysName             string              `protobuf:"bytes,1,opt,name=sysName,proto3" json:"sysName,omitempty"`
	SysDesc             string              `protobuf:"bytes,2,opt,name=sysDesc,proto3" json:"sysDesc,omitempty"`
	ChassisCapabilities []*LLDPCapabilities `protobuf:"bytes,3,rep,name=chassisCapabilities,proto3" json:"chassisCapabilities,omitempty"`
	// LLDP-MED
	DeviceType      string              `protobuf:"bytes,4,opt,name=deviceType,proto3" json:"deviceType,omitempty"`
	MedCapabilities []*LLDPCapabilities `protobuf:"bytes,5,rep,name=medCapabilities,proto3" json:"medCapabilities,omitempty"`
	// LLDP-MED-DEVICE
	InventoryHWRev    string `protobuf:"bytes,6,opt,name=inventoryHWRev,proto3" json:"inventoryHWRev,omitempty"`
	InventorySoftRev  string `protobuf:"bytes,7,opt,name=inventorySoftRev,proto3" json:"inventorySoftRev,omitempty"`
	InventorySerial   string `protobuf:"bytes,8,opt,name=inventorySerial,proto3" json:"inventorySerial,omitempty"`
	InventoryAssetTag string `protobuf:"bytes,9,opt,name=inventoryAssetTag,proto3" json:"inventoryAssetTag,omitempty"`
	InventoryModel    string `protobuf:"bytes,10,opt,name=inventoryModel,proto3" json:"inventoryModel,omitempty"`
	InventoryVendor   string `protobuf:"bytes,11,opt,name=inventoryVendor,proto3" json:"inventoryVendor,omitempty"`
}

func (x *LLDP) Reset() {
	*x = LLDP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LLDP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LLDP) ProtoMessage() {}

func (x *LLDP) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LLDP.ProtoReflect.Descriptor instead.
func (*LLDP) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{5}
}

func (x *LLDP) GetSysName() string {
	if x != nil {
		return x.SysName
	}
	return ""
}

func (x *LLDP) GetSysDesc() string {
	if x != nil {
		return x.SysDesc
	}
	return ""
}

func (x *LLDP) GetChassisCapabilities() []*LLDPCapabilities {
	if x != nil {
		return x.ChassisCapabilities
	}
	return nil
}

func (x *LLDP) GetDeviceType() string {
	if x != nil {
		return x.DeviceType
	}
	return ""
}

func (x *LLDP) GetMedCapabilities() []*LLDPCapabilities {
	if x != nil {
		return x.MedCapabilities
	}
	return nil
}

func (x *LLDP) GetInventoryHWRev() string {
	if x != nil {
		return x.InventoryHWRev
	}
	return ""
}

func (x *LLDP) GetInventorySoftRev() string {
	if x != nil {
		return x.InventorySoftRev
	}
	return ""
}

func (x *LLDP) GetInventorySerial() string {
	if x != nil {
		return x.InventorySerial
	}
	return ""
}

func (x *LLDP) GetInventoryAssetTag() string {
	if x != nil {
		return x.InventoryAssetTag
	}
	return ""
}

func (x *LLDP) GetInventoryModel() string {
	if x != nil {
		return x.InventoryModel
	}
	return ""
}

func (x *LLDP) GetInventoryVendor() string {
	if x != nil {
		return x.InventoryVendor
	}
	return ""
}

type LLDPCapabilities struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Capability string `protobuf:"bytes,1,opt,name=capability,proto3" json:"capability,omitempty"`
	Enabled    bool   `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
}

func (x *LLDPCapabilities) Reset() {
	*x = LLDPCapabilities{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LLDPCapabilities) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LLDPCapabilities) ProtoMessage() {}

func (x *LLDPCapabilities) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LLDPCapabilities.ProtoReflect.Descriptor instead.
func (*LLDPCapabilities) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{6}
}

func (x *LLDPCapabilities) GetCapability() string {
	if x != nil {
		return x.Capability
	}
	return ""
}

func (x *LLDPCapabilities) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

type ARP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip  string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Mac string `protobuf:"bytes,2,opt,name=mac,proto3" json:"mac,omitempty"`
}

func (x *ARP) Reset() {
	*x = ARP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ARP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ARP) ProtoMessage() {}

func (x *ARP) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ARP.ProtoReflect.Descriptor instead.
func (*ARP) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{7}
}

func (x *ARP) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *ARP) GetMac() string {
	if x != nil {
		return x.Mac
	}
	return ""
}

type NMAP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname  string       `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	MacVendor string       `protobuf:"bytes,2,opt,name=macVendor,proto3" json:"macVendor,omitempty"`
	Uptime    string       `protobuf:"bytes,3,opt,name=uptime,proto3" json:"uptime,omitempty"`
	LastBoot  string       `protobuf:"bytes,4,opt,name=lastBoot,proto3" json:"lastBoot,omitempty"`
	Os        string       `protobuf:"bytes,5,opt,name=os,proto3" json:"os,omitempty"`
	OpenPorts []*NMAPPorts `protobuf:"bytes,6,rep,name=openPorts,proto3" json:"openPorts,omitempty"`
}

func (x *NMAP) Reset() {
	*x = NMAP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NMAP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NMAP) ProtoMessage() {}

func (x *NMAP) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NMAP.ProtoReflect.Descriptor instead.
func (*NMAP) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{8}
}

func (x *NMAP) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *NMAP) GetMacVendor() string {
	if x != nil {
		return x.MacVendor
	}
	return ""
}

func (x *NMAP) GetUptime() string {
	if x != nil {
		return x.Uptime
	}
	return ""
}

func (x *NMAP) GetLastBoot() string {
	if x != nil {
		return x.LastBoot
	}
	return ""
}

func (x *NMAP) GetOs() string {
	if x != nil {
		return x.Os
	}
	return ""
}

func (x *NMAP) GetOpenPorts() []*NMAPPorts {
	if x != nil {
		return x.OpenPorts
	}
	return nil
}

type NMAPPorts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port     int32  `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Protocol string `protobuf:"bytes,2,opt,name=protocol,proto3" json:"protocol,omitempty"`
	State    string `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *NMAPPorts) Reset() {
	*x = NMAPPorts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Discovery_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NMAPPorts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NMAPPorts) ProtoMessage() {}

func (x *NMAPPorts) ProtoReflect() protoreflect.Message {
	mi := &file_Discovery_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NMAPPorts.ProtoReflect.Descriptor instead.
func (*NMAPPorts) Descriptor() ([]byte, []int) {
	return file_Discovery_proto_rawDescGZIP(), []int{9}
}

func (x *NMAPPorts) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *NMAPPorts) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *NMAPPorts) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

var File_Discovery_proto protoreflect.FileDescriptor

var file_Discovery_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x22, 0x0c, 0x0a, 0x0a,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x22, 0x42, 0x0a, 0x0f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x22,
	0x0a, 0x0e, 0x53, 0x63, 0x61, 0x6e, 0x4e, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6e, 0x65, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x6e,
	0x65, 0x74, 0x22, 0x25, 0x0a, 0x0f, 0x53, 0x63, 0x61, 0x6e, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0xde, 0x01, 0x0a, 0x0e, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1e, 0x0a, 0x0a,
	0x6d, 0x61, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x6d, 0x61, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x20, 0x0a, 0x0b,
	0x49, 0x50, 0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x49, 0x50, 0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1e,
	0x0a, 0x0a, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x23,
	0x0a, 0x04, 0x6c, 0x6c, 0x64, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x4c, 0x4c, 0x44, 0x50, 0x52, 0x04, 0x6c,
	0x6c, 0x64, 0x70, 0x12, 0x20, 0x0a, 0x03, 0x61, 0x72, 0x70, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x41, 0x52, 0x50,
	0x52, 0x03, 0x61, 0x72, 0x70, 0x12, 0x23, 0x0a, 0x04, 0x6e, 0x6d, 0x61, 0x70, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e,
	0x4e, 0x4d, 0x41, 0x50, 0x52, 0x04, 0x6e, 0x6d, 0x61, 0x70, 0x22, 0xee, 0x03, 0x0a, 0x04, 0x4c,
	0x4c, 0x44, 0x50, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x79, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x79, 0x73, 0x44, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x79, 0x73, 0x44, 0x65, 0x73, 0x63, 0x12, 0x4d, 0x0a, 0x13, 0x63, 0x68, 0x61, 0x73, 0x73,
	0x69, 0x73, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64,
	0x2e, 0x4c, 0x4c, 0x44, 0x50, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x52, 0x13, 0x63, 0x68, 0x61, 0x73, 0x73, 0x69, 0x73, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x45, 0x0a, 0x0f, 0x6d, 0x65, 0x64, 0x43, 0x61, 0x70,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x4c, 0x4c, 0x44, 0x50,
	0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x0f, 0x6d, 0x65,
	0x64, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x26, 0x0a,
	0x0e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x48, 0x57, 0x52, 0x65, 0x76, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79,
	0x48, 0x57, 0x52, 0x65, 0x76, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x53, 0x6f, 0x66, 0x74, 0x52, 0x65, 0x76, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x10, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x6f, 0x66, 0x74, 0x52, 0x65,
	0x76, 0x12, 0x28, 0x0a, 0x0f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x69, 0x6e, 0x76, 0x65,
	0x6e, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x2c, 0x0a, 0x11, 0x69,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x61, 0x67,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x61, 0x67, 0x12, 0x26, 0x0a, 0x0e, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x4d, 0x6f, 0x64, 0x65,
	0x6c, 0x12, 0x28, 0x0a, 0x0f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x56, 0x65,
	0x6e, 0x64, 0x6f, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x69, 0x6e, 0x76, 0x65,
	0x6e, 0x74, 0x6f, 0x72, 0x79, 0x56, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x22, 0x4c, 0x0a, 0x10, 0x4c,
	0x4c, 0x44, 0x50, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12,
	0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x22, 0x27, 0x0a, 0x03, 0x41, 0x52, 0x50,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d,
	0x61, 0x63, 0x22, 0xb8, 0x01, 0x0a, 0x04, 0x4e, 0x4d, 0x41, 0x50, 0x12, 0x1a, 0x0a, 0x08, 0x68,
	0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68,
	0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x63, 0x56, 0x65,
	0x6e, 0x64, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x61, 0x63, 0x56,
	0x65, 0x6e, 0x64, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x61, 0x73, 0x74, 0x42, 0x6f, 0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x61, 0x73, 0x74, 0x42, 0x6f, 0x6f, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x32, 0x0a, 0x09, 0x6f, 0x70, 0x65,
	0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x4e, 0x4d, 0x41, 0x50, 0x50, 0x6f, 0x72,
	0x74, 0x73, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x22, 0x51, 0x0a,
	0x09, 0x4e, 0x4d, 0x41, 0x50, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x2a, 0x43, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x06, 0x0a,
	0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x02,
	0x12, 0x13, 0x0a, 0x0f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55,
	0x45, 0x53, 0x54, 0x10, 0x03, 0x32, 0xdf, 0x01, 0x0a, 0x09, 0x44, 0x69, 0x73, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x64, 0x12, 0x48, 0x0a, 0x11, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x6c,
	0x6c, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x15, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a,
	0x1a, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x42, 0x0a,
	0x07, 0x53, 0x63, 0x61, 0x6e, 0x4e, 0x65, 0x74, 0x12, 0x19, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x64, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x4e, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x44, 0x0a, 0x08, 0x53, 0x63, 0x61, 0x6e, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x1a, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x48, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67,
	0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66,
	0x66, 0x65, 0x72, 0x73, 0x2f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x64, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Discovery_proto_rawDescOnce sync.Once
	file_Discovery_proto_rawDescData = file_Discovery_proto_rawDesc
)

func file_Discovery_proto_rawDescGZIP() []byte {
	file_Discovery_proto_rawDescOnce.Do(func() {
		file_Discovery_proto_rawDescData = protoimpl.X.CompressGZIP(file_Discovery_proto_rawDescData)
	})
	return file_Discovery_proto_rawDescData
}

var file_Discovery_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_Discovery_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_Discovery_proto_goTypes = []interface{}{
	(ResponseCode)(0),        // 0: discoverd.ResponseCode
	(*EmptyParam)(nil),       // 1: discoverd.EmptyParam
	(*RequestResponse)(nil),  // 2: discoverd.RequestResponse
	(*ScanNetRequest)(nil),   // 3: discoverd.ScanNetRequest
	(*ScanHostRequest)(nil),  // 4: discoverd.ScanHostRequest
	(*DiscoveryEntry)(nil),   // 5: discoverd.DiscoveryEntry
	(*LLDP)(nil),             // 6: discoverd.LLDP
	(*LLDPCapabilities)(nil), // 7: discoverd.LLDPCapabilities
	(*ARP)(nil),              // 8: discoverd.ARP
	(*NMAP)(nil),             // 9: discoverd.NMAP
	(*NMAPPorts)(nil),        // 10: discoverd.NMAPPorts
}
var file_Discovery_proto_depIdxs = []int32{
	0,  // 0: discoverd.RequestResponse.result:type_name -> discoverd.ResponseCode
	6,  // 1: discoverd.DiscoveryEntry.lldp:type_name -> discoverd.LLDP
	8,  // 2: discoverd.DiscoveryEntry.arp:type_name -> discoverd.ARP
	9,  // 3: discoverd.DiscoveryEntry.nmap:type_name -> discoverd.NMAP
	7,  // 4: discoverd.LLDP.chassisCapabilities:type_name -> discoverd.LLDPCapabilities
	7,  // 5: discoverd.LLDP.medCapabilities:type_name -> discoverd.LLDPCapabilities
	10, // 6: discoverd.NMAP.openPorts:type_name -> discoverd.NMAPPorts
	1,  // 7: discoverd.Disoveryd.RequestAllEntries:input_type -> discoverd.EmptyParam
	3,  // 8: discoverd.Disoveryd.ScanNet:input_type -> discoverd.ScanNetRequest
	4,  // 9: discoverd.Disoveryd.ScanHost:input_type -> discoverd.ScanHostRequest
	2,  // 10: discoverd.Disoveryd.RequestAllEntries:output_type -> discoverd.RequestResponse
	2,  // 11: discoverd.Disoveryd.ScanNet:output_type -> discoverd.RequestResponse
	2,  // 12: discoverd.Disoveryd.ScanHost:output_type -> discoverd.RequestResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_Discovery_proto_init() }
func file_Discovery_proto_init() {
	if File_Discovery_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Discovery_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyParam); i {
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
		file_Discovery_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestResponse); i {
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
		file_Discovery_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanNetRequest); i {
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
		file_Discovery_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanHostRequest); i {
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
		file_Discovery_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryEntry); i {
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
		file_Discovery_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LLDP); i {
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
		file_Discovery_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LLDPCapabilities); i {
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
		file_Discovery_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ARP); i {
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
		file_Discovery_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NMAP); i {
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
		file_Discovery_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NMAPPorts); i {
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
			RawDescriptor: file_Discovery_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_Discovery_proto_goTypes,
		DependencyIndexes: file_Discovery_proto_depIdxs,
		EnumInfos:         file_Discovery_proto_enumTypes,
		MessageInfos:      file_Discovery_proto_msgTypes,
	}.Build()
	File_Discovery_proto = out.File
	file_Discovery_proto_rawDesc = nil
	file_Discovery_proto_goTypes = nil
	file_Discovery_proto_depIdxs = nil
}

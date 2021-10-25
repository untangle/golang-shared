// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.6.1
// source: ActiveSessions.proto

package ActiveSessions

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

type ActiveSessions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionsList []*Session `protobuf:"bytes,1,rep,name=sessionsList,proto3" json:"sessionsList,omitempty"`
}

func (x *ActiveSessions) Reset() {
	*x = ActiveSessions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ActiveSessions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActiveSessions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActiveSessions) ProtoMessage() {}

func (x *ActiveSessions) ProtoReflect() protoreflect.Message {
	mi := &file_ActiveSessions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActiveSessions.ProtoReflect.Descriptor instead.
func (*ActiveSessions) Descriptor() ([]byte, []int) {
	return file_ActiveSessions_proto_rawDescGZIP(), []int{0}
}

func (x *ActiveSessions) GetSessionsList() []*Session {
	if x != nil {
		return x.SessionsList
	}
	return nil
}

type Session struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AgeMilliseconds                 uint64 `protobuf:"varint,1,opt,name=age_milliseconds,json=ageMilliseconds,proto3" json:"age_milliseconds,omitempty"`
	Bytes                           uint64 `protobuf:"varint,3,opt,name=bytes,proto3" json:"bytes,omitempty"`
	ClientBytes                     uint64 `protobuf:"varint,4,opt,name=client_bytes,json=clientBytes,proto3" json:"client_bytes,omitempty"`
	ServerBytes                     uint64 `protobuf:"varint,5,opt,name=server_bytes,json=serverBytes,proto3" json:"server_bytes,omitempty"`
	ByteRate                        int64  `protobuf:"varint,6,opt,name=byte_rate,json=byteRate,proto3" json:"byte_rate,omitempty"`
	ClientByteRate                  int64  `protobuf:"varint,7,opt,name=client_byte_rate,json=clientByteRate,proto3" json:"client_byte_rate,omitempty"`
	ServerByteRate                  int64  `protobuf:"varint,8,opt,name=server_byte_rate,json=serverByteRate,proto3" json:"server_byte_rate,omitempty"`
	Packets                         uint64 `protobuf:"varint,9,opt,name=packets,proto3" json:"packets,omitempty"`
	ClientPackets                   uint64 `protobuf:"varint,10,opt,name=client_packets,json=clientPackets,proto3" json:"client_packets,omitempty"`
	ServerPackets                   uint64 `protobuf:"varint,11,opt,name=server_packets,json=serverPackets,proto3" json:"server_packets,omitempty"`
	PacketRate                      int64  `protobuf:"varint,12,opt,name=packet_rate,json=packetRate,proto3" json:"packet_rate,omitempty"`
	ClientPacketRate                int64  `protobuf:"varint,13,opt,name=client_packet_rate,json=clientPacketRate,proto3" json:"client_packet_rate,omitempty"`
	ServerPacketRate                int64  `protobuf:"varint,14,opt,name=server_packet_rate,json=serverPacketRate,proto3" json:"server_packet_rate,omitempty"`
	ClientAddress                   string `protobuf:"bytes,15,opt,name=client_address,json=clientAddress,proto3" json:"client_address,omitempty"`
	ClientAddressNew                string `protobuf:"bytes,16,opt,name=client_address_new,json=clientAddressNew,proto3" json:"client_address_new,omitempty"`
	ClientInterfaceId               uint32 `protobuf:"varint,17,opt,name=client_interface_id,json=clientInterfaceId,proto3" json:"client_interface_id,omitempty"`
	ClientInterfaceType             uint32 `protobuf:"varint,18,opt,name=client_interface_type,json=clientInterfaceType,proto3" json:"client_interface_type,omitempty"`
	ConntrackId                     uint32 `protobuf:"varint,20,opt,name=conntrack_id,json=conntrackId,proto3" json:"conntrack_id,omitempty"`
	Family                          uint32 `protobuf:"varint,21,opt,name=family,proto3" json:"family,omitempty"`
	IpProtocol                      uint32 `protobuf:"varint,22,opt,name=ip_protocol,json=ipProtocol,proto3" json:"ip_protocol,omitempty"`
	Mark                            uint32 `protobuf:"varint,23,opt,name=mark,proto3" json:"mark,omitempty"`
	Priority                        uint32 `protobuf:"varint,24,opt,name=priority,proto3" json:"priority,omitempty"`
	ServerAddress                   string `protobuf:"bytes,25,opt,name=server_address,json=serverAddress,proto3" json:"server_address,omitempty"`
	ServerAddressNew                string `protobuf:"bytes,26,opt,name=server_address_new,json=serverAddressNew,proto3" json:"server_address_new,omitempty"`
	ServerInterfaceId               uint32 `protobuf:"varint,27,opt,name=server_interface_id,json=serverInterfaceId,proto3" json:"server_interface_id,omitempty"`
	ServerInterfaceType             uint32 `protobuf:"varint,28,opt,name=server_interface_type,json=serverInterfaceType,proto3" json:"server_interface_type,omitempty"`
	ServerPort                      uint32 `protobuf:"varint,29,opt,name=server_port,json=serverPort,proto3" json:"server_port,omitempty"`
	ServerPortNew                   uint32 `protobuf:"varint,30,opt,name=server_port_new,json=serverPortNew,proto3" json:"server_port_new,omitempty"`
	SessionId                       int64  `protobuf:"varint,31,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	TcpState                        uint32 `protobuf:"varint,32,opt,name=tcp_state,json=tcpState,proto3" json:"tcp_state,omitempty"`
	TimestampStart                  uint64 `protobuf:"varint,33,opt,name=timestamp_start,json=timestampStart,proto3" json:"timestamp_start,omitempty"`
	WanPolicy                       string `protobuf:"bytes,34,opt,name=wan_policy,json=wanPolicy,proto3" json:"wan_policy,omitempty"`
	ApplicationCategory             string `protobuf:"bytes,35,opt,name=application_category,json=applicationCategory,proto3" json:"application_category,omitempty"`
	ApplicationCategoryInferred     string `protobuf:"bytes,36,opt,name=application_category_inferred,json=applicationCategoryInferred,proto3" json:"application_category_inferred,omitempty"`
	ApplicationConfidence           int32  `protobuf:"varint,37,opt,name=application_confidence,json=applicationConfidence,proto3" json:"application_confidence,omitempty"`
	ApplicationConfidenceInferred   int32  `protobuf:"varint,38,opt,name=application_confidence_inferred,json=applicationConfidenceInferred,proto3" json:"application_confidence_inferred,omitempty"`
	ApplicationId                   string `protobuf:"bytes,39,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	ApplicationIdInferred           string `protobuf:"bytes,40,opt,name=application_id_inferred,json=applicationIdInferred,proto3" json:"application_id_inferred,omitempty"`
	ApplicationProductivity         int32  `protobuf:"varint,41,opt,name=application_productivity,json=applicationProductivity,proto3" json:"application_productivity,omitempty"`
	ApplicationProductivityInferred int32  `protobuf:"varint,42,opt,name=application_productivity_inferred,json=applicationProductivityInferred,proto3" json:"application_productivity_inferred,omitempty"`
	ApplicationProtochain           string `protobuf:"bytes,43,opt,name=application_protochain,json=applicationProtochain,proto3" json:"application_protochain,omitempty"`
	ApplicationProtochainInferred   string `protobuf:"bytes,44,opt,name=application_protochain_inferred,json=applicationProtochainInferred,proto3" json:"application_protochain_inferred,omitempty"`
	ApplicationRisk                 int32  `protobuf:"varint,45,opt,name=application_risk,json=applicationRisk,proto3" json:"application_risk,omitempty"`
	ApplicationRiskInferred         int32  `protobuf:"varint,46,opt,name=application_risk_inferred,json=applicationRiskInferred,proto3" json:"application_risk_inferred,omitempty"`
	CertDnsNames                    string `protobuf:"bytes,47,opt,name=cert_dns_names,json=certDnsNames,proto3" json:"cert_dns_names,omitempty"`
	CertificateIssuerC              string `protobuf:"bytes,48,opt,name=certificate_issuer_c,json=certificateIssuerC,proto3" json:"certificate_issuer_c,omitempty"`
	CertificateIssuerCn             string `protobuf:"bytes,49,opt,name=certificate_issuer_cn,json=certificateIssuerCn,proto3" json:"certificate_issuer_cn,omitempty"`
	CertificateIssuerL              string `protobuf:"bytes,50,opt,name=certificate_issuer_l,json=certificateIssuerL,proto3" json:"certificate_issuer_l,omitempty"`
	CertificateIssuerO              string `protobuf:"bytes,51,opt,name=certificate_issuer_o,json=certificateIssuerO,proto3" json:"certificate_issuer_o,omitempty"`
	CertificateIssuerOu             string `protobuf:"bytes,52,opt,name=certificate_issuer_ou,json=certificateIssuerOu,proto3" json:"certificate_issuer_ou,omitempty"`
	CertificateIssuerP              string `protobuf:"bytes,53,opt,name=certificate_issuer_p,json=certificateIssuerP,proto3" json:"certificate_issuer_p,omitempty"`
	CertificateSubjectCn            string `protobuf:"bytes,54,opt,name=certificate_subject_cn,json=certificateSubjectCn,proto3" json:"certificate_subject_cn,omitempty"`
	CertificateSubjectO             string `protobuf:"bytes,55,opt,name=certificate_subject_o,json=certificateSubjectO,proto3" json:"certificate_subject_o,omitempty"`
	CertificateSubjectSan           string `protobuf:"bytes,56,opt,name=certificate_subject_san,json=certificateSubjectSan,proto3" json:"certificate_subject_san,omitempty"`
	ServerReverseDns                string `protobuf:"bytes,57,opt,name=server_reverse_dns,json=serverReverseDns,proto3" json:"server_reverse_dns,omitempty"`
}

func (x *Session) Reset() {
	*x = Session{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ActiveSessions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Session) ProtoMessage() {}

func (x *Session) ProtoReflect() protoreflect.Message {
	mi := &file_ActiveSessions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Session.ProtoReflect.Descriptor instead.
func (*Session) Descriptor() ([]byte, []int) {
	return file_ActiveSessions_proto_rawDescGZIP(), []int{1}
}

func (x *Session) GetAgeMilliseconds() uint64 {
	if x != nil {
		return x.AgeMilliseconds
	}
	return 0
}

func (x *Session) GetBytes() uint64 {
	if x != nil {
		return x.Bytes
	}
	return 0
}

func (x *Session) GetClientBytes() uint64 {
	if x != nil {
		return x.ClientBytes
	}
	return 0
}

func (x *Session) GetServerBytes() uint64 {
	if x != nil {
		return x.ServerBytes
	}
	return 0
}

func (x *Session) GetByteRate() int64 {
	if x != nil {
		return x.ByteRate
	}
	return 0
}

func (x *Session) GetClientByteRate() int64 {
	if x != nil {
		return x.ClientByteRate
	}
	return 0
}

func (x *Session) GetServerByteRate() int64 {
	if x != nil {
		return x.ServerByteRate
	}
	return 0
}

func (x *Session) GetPackets() uint64 {
	if x != nil {
		return x.Packets
	}
	return 0
}

func (x *Session) GetClientPackets() uint64 {
	if x != nil {
		return x.ClientPackets
	}
	return 0
}

func (x *Session) GetServerPackets() uint64 {
	if x != nil {
		return x.ServerPackets
	}
	return 0
}

func (x *Session) GetPacketRate() int64 {
	if x != nil {
		return x.PacketRate
	}
	return 0
}

func (x *Session) GetClientPacketRate() int64 {
	if x != nil {
		return x.ClientPacketRate
	}
	return 0
}

func (x *Session) GetServerPacketRate() int64 {
	if x != nil {
		return x.ServerPacketRate
	}
	return 0
}

func (x *Session) GetClientAddress() string {
	if x != nil {
		return x.ClientAddress
	}
	return ""
}

func (x *Session) GetClientAddressNew() string {
	if x != nil {
		return x.ClientAddressNew
	}
	return ""
}

func (x *Session) GetClientInterfaceId() uint32 {
	if x != nil {
		return x.ClientInterfaceId
	}
	return 0
}

func (x *Session) GetClientInterfaceType() uint32 {
	if x != nil {
		return x.ClientInterfaceType
	}
	return 0
}

func (x *Session) GetConntrackId() uint32 {
	if x != nil {
		return x.ConntrackId
	}
	return 0
}

func (x *Session) GetFamily() uint32 {
	if x != nil {
		return x.Family
	}
	return 0
}

func (x *Session) GetIpProtocol() uint32 {
	if x != nil {
		return x.IpProtocol
	}
	return 0
}

func (x *Session) GetMark() uint32 {
	if x != nil {
		return x.Mark
	}
	return 0
}

func (x *Session) GetPriority() uint32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *Session) GetServerAddress() string {
	if x != nil {
		return x.ServerAddress
	}
	return ""
}

func (x *Session) GetServerAddressNew() string {
	if x != nil {
		return x.ServerAddressNew
	}
	return ""
}

func (x *Session) GetServerInterfaceId() uint32 {
	if x != nil {
		return x.ServerInterfaceId
	}
	return 0
}

func (x *Session) GetServerInterfaceType() uint32 {
	if x != nil {
		return x.ServerInterfaceType
	}
	return 0
}

func (x *Session) GetServerPort() uint32 {
	if x != nil {
		return x.ServerPort
	}
	return 0
}

func (x *Session) GetServerPortNew() uint32 {
	if x != nil {
		return x.ServerPortNew
	}
	return 0
}

func (x *Session) GetSessionId() int64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *Session) GetTcpState() uint32 {
	if x != nil {
		return x.TcpState
	}
	return 0
}

func (x *Session) GetTimestampStart() uint64 {
	if x != nil {
		return x.TimestampStart
	}
	return 0
}

func (x *Session) GetWanPolicy() string {
	if x != nil {
		return x.WanPolicy
	}
	return ""
}

func (x *Session) GetApplicationCategory() string {
	if x != nil {
		return x.ApplicationCategory
	}
	return ""
}

func (x *Session) GetApplicationCategoryInferred() string {
	if x != nil {
		return x.ApplicationCategoryInferred
	}
	return ""
}

func (x *Session) GetApplicationConfidence() int32 {
	if x != nil {
		return x.ApplicationConfidence
	}
	return 0
}

func (x *Session) GetApplicationConfidenceInferred() int32 {
	if x != nil {
		return x.ApplicationConfidenceInferred
	}
	return 0
}

func (x *Session) GetApplicationId() string {
	if x != nil {
		return x.ApplicationId
	}
	return ""
}

func (x *Session) GetApplicationIdInferred() string {
	if x != nil {
		return x.ApplicationIdInferred
	}
	return ""
}

func (x *Session) GetApplicationProductivity() int32 {
	if x != nil {
		return x.ApplicationProductivity
	}
	return 0
}

func (x *Session) GetApplicationProductivityInferred() int32 {
	if x != nil {
		return x.ApplicationProductivityInferred
	}
	return 0
}

func (x *Session) GetApplicationProtochain() string {
	if x != nil {
		return x.ApplicationProtochain
	}
	return ""
}

func (x *Session) GetApplicationProtochainInferred() string {
	if x != nil {
		return x.ApplicationProtochainInferred
	}
	return ""
}

func (x *Session) GetApplicationRisk() int32 {
	if x != nil {
		return x.ApplicationRisk
	}
	return 0
}

func (x *Session) GetApplicationRiskInferred() int32 {
	if x != nil {
		return x.ApplicationRiskInferred
	}
	return 0
}

func (x *Session) GetCertDnsNames() string {
	if x != nil {
		return x.CertDnsNames
	}
	return ""
}

func (x *Session) GetCertificateIssuerC() string {
	if x != nil {
		return x.CertificateIssuerC
	}
	return ""
}

func (x *Session) GetCertificateIssuerCn() string {
	if x != nil {
		return x.CertificateIssuerCn
	}
	return ""
}

func (x *Session) GetCertificateIssuerL() string {
	if x != nil {
		return x.CertificateIssuerL
	}
	return ""
}

func (x *Session) GetCertificateIssuerO() string {
	if x != nil {
		return x.CertificateIssuerO
	}
	return ""
}

func (x *Session) GetCertificateIssuerOu() string {
	if x != nil {
		return x.CertificateIssuerOu
	}
	return ""
}

func (x *Session) GetCertificateIssuerP() string {
	if x != nil {
		return x.CertificateIssuerP
	}
	return ""
}

func (x *Session) GetCertificateSubjectCn() string {
	if x != nil {
		return x.CertificateSubjectCn
	}
	return ""
}

func (x *Session) GetCertificateSubjectO() string {
	if x != nil {
		return x.CertificateSubjectO
	}
	return ""
}

func (x *Session) GetCertificateSubjectSan() string {
	if x != nil {
		return x.CertificateSubjectSan
	}
	return ""
}

func (x *Session) GetServerReverseDns() string {
	if x != nil {
		return x.ServerReverseDns
	}
	return ""
}

var File_ActiveSessions_proto protoreflect.FileDescriptor

var file_ActiveSessions_proto_rawDesc = []byte{
	0x0a, 0x14, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x22,
	0x46, 0x0a, 0x0e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x34, 0x0a, 0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x73, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x99, 0x13, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x69, 0x6c, 0x6c, 0x69,
	0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x61,
	0x67, 0x65, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x79,
	0x74, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62,
	0x79, 0x74, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x74, 0x65, 0x52, 0x61, 0x74,
	0x65, 0x12, 0x28, 0x0a, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x74, 0x65,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x42, 0x79, 0x74, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x25, 0x0a, 0x0e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x72, 0x61,
	0x74, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x70,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x6e, 0x65, 0x77, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x4e, 0x65, 0x77, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x11, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x11, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x12,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x13, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x66, 0x61,
	0x6d, 0x69, 0x6c, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x69, 0x70, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x17, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x18, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x19, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x12,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x6e,
	0x65, 0x77, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x12, 0x2e, 0x0a, 0x13, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x15, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x13, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x1d, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x12,
	0x26, 0x0a, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x6e,
	0x65, 0x77, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x50, 0x6f, 0x72, 0x74, 0x4e, 0x65, 0x77, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x63, 0x70, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x20, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x74, 0x63, 0x70, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x21, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x77, 0x61, 0x6e, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x22, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x77, 0x61, 0x6e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x31, 0x0a, 0x14, 0x61,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x18, 0x23, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x42,
	0x0a, 0x1d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x18,
	0x24, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1b, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x66, 0x65, 0x72, 0x72,
	0x65, 0x64, 0x12, 0x35, 0x0a, 0x16, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x25, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x15, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x1f, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65,
	0x6e, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x18, 0x26, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x1d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65,
	0x64, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x27, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x17, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x5f, 0x69, 0x6e, 0x66, 0x65, 0x72,
	0x72, 0x65, 0x64, 0x18, 0x28, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x49, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64,
	0x12, 0x39, 0x0a, 0x18, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x29, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x17, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x4a, 0x0a, 0x21, 0x61,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64,
	0x18, 0x2a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x1f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x49,
	0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x12, 0x35, 0x0a, 0x16, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x18, 0x2b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x46,
	0x0a, 0x1f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65,
	0x64, 0x18, 0x2c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e,
	0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x69, 0x73, 0x6b, 0x18, 0x2d, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x69, 0x73,
	0x6b, 0x12, 0x3a, 0x0a, 0x19, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x72, 0x69, 0x73, 0x6b, 0x5f, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x18, 0x2e,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x17, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x69, 0x73, 0x6b, 0x49, 0x6e, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x12, 0x24, 0x0a,
	0x0e, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x64, 0x6e, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18,
	0x2f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x65, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x5f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x63, 0x18, 0x30, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x12, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x49, 0x73,
	0x73, 0x75, 0x65, 0x72, 0x43, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x63, 0x6e, 0x18, 0x31,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x43, 0x6e, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f,
	0x6c, 0x18, 0x32, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x4c, 0x12, 0x30, 0x0a, 0x14, 0x63,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x72, 0x5f, 0x6f, 0x18, 0x33, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x4f, 0x12, 0x32, 0x0a,
	0x15, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x72, 0x5f, 0x6f, 0x75, 0x18, 0x34, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x4f,
	0x75, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x5f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x70, 0x18, 0x35, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x12, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x49, 0x73, 0x73, 0x75,
	0x65, 0x72, 0x50, 0x12, 0x34, 0x0a, 0x16, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x5f, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x63, 0x6e, 0x18, 0x36, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x14, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x43, 0x6e, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x6f, 0x18, 0x37, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x4f, 0x12, 0x36, 0x0a,
	0x17, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x5f, 0x73, 0x61, 0x6e, 0x18, 0x38, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15,
	0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x53, 0x61, 0x6e, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x5f, 0x64, 0x6e, 0x73, 0x18, 0x39, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65,
	0x44, 0x6e, 0x73, 0x42, 0x4a, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e,
	0x67, 0x2d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73,
	0x2f, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ActiveSessions_proto_rawDescOnce sync.Once
	file_ActiveSessions_proto_rawDescData = file_ActiveSessions_proto_rawDesc
)

func file_ActiveSessions_proto_rawDescGZIP() []byte {
	file_ActiveSessions_proto_rawDescOnce.Do(func() {
		file_ActiveSessions_proto_rawDescData = protoimpl.X.CompressGZIP(file_ActiveSessions_proto_rawDescData)
	})
	return file_ActiveSessions_proto_rawDescData
}

var file_ActiveSessions_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ActiveSessions_proto_goTypes = []interface{}{
	(*ActiveSessions)(nil), // 0: reports.ActiveSessions
	(*Session)(nil),        // 1: reports.Session
}
var file_ActiveSessions_proto_depIdxs = []int32{
	1, // 0: reports.ActiveSessions.sessionsList:type_name -> reports.Session
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ActiveSessions_proto_init() }
func file_ActiveSessions_proto_init() {
	if File_ActiveSessions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ActiveSessions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActiveSessions); i {
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
		file_ActiveSessions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Session); i {
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
			RawDescriptor: file_ActiveSessions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ActiveSessions_proto_goTypes,
		DependencyIndexes: file_ActiveSessions_proto_depIdxs,
		MessageInfos:      file_ActiveSessions_proto_msgTypes,
	}.Build()
	File_ActiveSessions_proto = out.File
	file_ActiveSessions_proto_rawDesc = nil
	file_ActiveSessions_proto_goTypes = nil
	file_ActiveSessions_proto_depIdxs = nil
}

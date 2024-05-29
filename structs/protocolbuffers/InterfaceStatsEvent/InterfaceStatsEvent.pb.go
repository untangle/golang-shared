// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.21.12
// source: InterfaceStatsEvent.proto

package InterfaceStatsEvent

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

type InterfaceStatsEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeStamp              int64   `protobuf:"varint,1,opt,name=timeStamp,proto3" json:"timeStamp,omitempty"`
	InterfaceID            int32   `protobuf:"varint,2,opt,name=interfaceID,proto3" json:"interfaceID,omitempty"`
	InterfaceName          string  `protobuf:"bytes,3,opt,name=interfaceName,proto3" json:"interfaceName,omitempty"`
	DeviceName             string  `protobuf:"bytes,4,opt,name=deviceName,proto3" json:"deviceName,omitempty"`
	IsWan                  bool    `protobuf:"varint,5,opt,name=isWan,proto3" json:"isWan,omitempty"`
	Latency1               float64 `protobuf:"fixed64,6,opt,name=latency1,proto3" json:"latency1,omitempty"`
	Latency5               float64 `protobuf:"fixed64,7,opt,name=latency5,proto3" json:"latency5,omitempty"`
	Latency15              float64 `protobuf:"fixed64,8,opt,name=latency15,proto3" json:"latency15,omitempty"`
	LatencyVariance        float64 `protobuf:"fixed64,9,opt,name=latencyVariance,proto3" json:"latencyVariance,omitempty"`
	PassiveLatency1        float64 `protobuf:"fixed64,10,opt,name=passiveLatency1,proto3" json:"passiveLatency1,omitempty"`
	PassiveLatency5        float64 `protobuf:"fixed64,11,opt,name=passiveLatency5,proto3" json:"passiveLatency5,omitempty"`
	PassiveLatency15       float64 `protobuf:"fixed64,12,opt,name=passiveLatency15,proto3" json:"passiveLatency15,omitempty"`
	PassiveLatencyVariance float64 `protobuf:"fixed64,13,opt,name=passiveLatencyVariance,proto3" json:"passiveLatencyVariance,omitempty"`
	ActiveLatency1         float64 `protobuf:"fixed64,14,opt,name=activeLatency1,proto3" json:"activeLatency1,omitempty"`
	ActiveLatency5         float64 `protobuf:"fixed64,15,opt,name=activeLatency5,proto3" json:"activeLatency5,omitempty"`
	ActiveLatency15        float64 `protobuf:"fixed64,16,opt,name=activeLatency15,proto3" json:"activeLatency15,omitempty"`
	ActiveLatencyVariance  float64 `protobuf:"fixed64,17,opt,name=activeLatencyVariance,proto3" json:"activeLatencyVariance,omitempty"`
	Jitter1                float64 `protobuf:"fixed64,18,opt,name=jitter1,proto3" json:"jitter1,omitempty"`
	Jitter5                float64 `protobuf:"fixed64,19,opt,name=jitter5,proto3" json:"jitter5,omitempty"`
	Jitter15               float64 `protobuf:"fixed64,20,opt,name=jitter15,proto3" json:"jitter15,omitempty"`
	JitterVariance         float64 `protobuf:"fixed64,21,opt,name=jitterVariance,proto3" json:"jitterVariance,omitempty"`
	PingTimeout            uint64  `protobuf:"varint,22,opt,name=pingTimeout,proto3" json:"pingTimeout,omitempty"`
	PingTimeoutRate        uint64  `protobuf:"varint,23,opt,name=pingTimeoutRate,proto3" json:"pingTimeoutRate,omitempty"`
	RxBytes                uint64  `protobuf:"varint,24,opt,name=rxBytes,proto3" json:"rxBytes,omitempty"`
	RxBytesRate            uint64  `protobuf:"varint,25,opt,name=rxBytesRate,proto3" json:"rxBytesRate,omitempty"`
	RxPackets              uint64  `protobuf:"varint,26,opt,name=rxPackets,proto3" json:"rxPackets,omitempty"`
	RxPacketsRate          uint64  `protobuf:"varint,27,opt,name=rxPacketsRate,proto3" json:"rxPacketsRate,omitempty"`
	RxErrs                 uint64  `protobuf:"varint,28,opt,name=rxErrs,proto3" json:"rxErrs,omitempty"`
	RxErrsRate             uint64  `protobuf:"varint,29,opt,name=rxErrsRate,proto3" json:"rxErrsRate,omitempty"`
	RxDrop                 uint64  `protobuf:"varint,30,opt,name=rxDrop,proto3" json:"rxDrop,omitempty"`
	RxDropRate             uint64  `protobuf:"varint,31,opt,name=rxDropRate,proto3" json:"rxDropRate,omitempty"`
	RxFifo                 uint64  `protobuf:"varint,32,opt,name=rxFifo,proto3" json:"rxFifo,omitempty"`
	RxFifoRate             uint64  `protobuf:"varint,33,opt,name=rxFifoRate,proto3" json:"rxFifoRate,omitempty"`
	RxFrame                uint64  `protobuf:"varint,34,opt,name=rxFrame,proto3" json:"rxFrame,omitempty"`
	RxFrameRate            uint64  `protobuf:"varint,35,opt,name=rxFrameRate,proto3" json:"rxFrameRate,omitempty"`
	RxCompressed           uint64  `protobuf:"varint,36,opt,name=rxCompressed,proto3" json:"rxCompressed,omitempty"`
	RxCompressedRate       uint64  `protobuf:"varint,37,opt,name=rxCompressedRate,proto3" json:"rxCompressedRate,omitempty"`
	RxMulticast            uint64  `protobuf:"varint,38,opt,name=rxMulticast,proto3" json:"rxMulticast,omitempty"`
	RxMulticastRate        uint64  `protobuf:"varint,39,opt,name=rxMulticastRate,proto3" json:"rxMulticastRate,omitempty"`
	TxBytes                uint64  `protobuf:"varint,40,opt,name=txBytes,proto3" json:"txBytes,omitempty"`
	TxBytesRate            uint64  `protobuf:"varint,41,opt,name=txBytesRate,proto3" json:"txBytesRate,omitempty"`
	TxPackets              uint64  `protobuf:"varint,42,opt,name=txPackets,proto3" json:"txPackets,omitempty"`
	TxPacketsRate          uint64  `protobuf:"varint,43,opt,name=txPacketsRate,proto3" json:"txPacketsRate,omitempty"`
	TxErrs                 uint64  `protobuf:"varint,44,opt,name=txErrs,proto3" json:"txErrs,omitempty"`
	TxErrsRate             uint64  `protobuf:"varint,45,opt,name=txErrsRate,proto3" json:"txErrsRate,omitempty"`
	TxDrop                 uint64  `protobuf:"varint,46,opt,name=txDrop,proto3" json:"txDrop,omitempty"`
	TxDropRate             uint64  `protobuf:"varint,47,opt,name=txDropRate,proto3" json:"txDropRate,omitempty"`
	TxFifo                 uint64  `protobuf:"varint,48,opt,name=txFifo,proto3" json:"txFifo,omitempty"`
	TxFifoRate             uint64  `protobuf:"varint,49,opt,name=txFifoRate,proto3" json:"txFifoRate,omitempty"`
	TxColls                uint64  `protobuf:"varint,50,opt,name=txColls,proto3" json:"txColls,omitempty"`
	TxCollsRate            uint64  `protobuf:"varint,51,opt,name=txCollsRate,proto3" json:"txCollsRate,omitempty"`
	TxCarrier              uint64  `protobuf:"varint,52,opt,name=txCarrier,proto3" json:"txCarrier,omitempty"`
	TxCarrierRate          uint64  `protobuf:"varint,53,opt,name=txCarrierRate,proto3" json:"txCarrierRate,omitempty"`
	TxCompressed           uint64  `protobuf:"varint,54,opt,name=txCompressed,proto3" json:"txCompressed,omitempty"`
	TxCompressedRate       uint64  `protobuf:"varint,55,opt,name=txCompressedRate,proto3" json:"txCompressedRate,omitempty"`
	Offline                bool    `protobuf:"varint,56,opt,name=offline,proto3" json:"offline,omitempty"`
}

func (x *InterfaceStatsEvent) Reset() {
	*x = InterfaceStatsEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_InterfaceStatsEvent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterfaceStatsEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterfaceStatsEvent) ProtoMessage() {}

func (x *InterfaceStatsEvent) ProtoReflect() protoreflect.Message {
	mi := &file_InterfaceStatsEvent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterfaceStatsEvent.ProtoReflect.Descriptor instead.
func (*InterfaceStatsEvent) Descriptor() ([]byte, []int) {
	return file_InterfaceStatsEvent_proto_rawDescGZIP(), []int{0}
}

func (x *InterfaceStatsEvent) GetTimeStamp() int64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *InterfaceStatsEvent) GetInterfaceID() int32 {
	if x != nil {
		return x.InterfaceID
	}
	return 0
}

func (x *InterfaceStatsEvent) GetInterfaceName() string {
	if x != nil {
		return x.InterfaceName
	}
	return ""
}

func (x *InterfaceStatsEvent) GetDeviceName() string {
	if x != nil {
		return x.DeviceName
	}
	return ""
}

func (x *InterfaceStatsEvent) GetIsWan() bool {
	if x != nil {
		return x.IsWan
	}
	return false
}

func (x *InterfaceStatsEvent) GetLatency1() float64 {
	if x != nil {
		return x.Latency1
	}
	return 0
}

func (x *InterfaceStatsEvent) GetLatency5() float64 {
	if x != nil {
		return x.Latency5
	}
	return 0
}

func (x *InterfaceStatsEvent) GetLatency15() float64 {
	if x != nil {
		return x.Latency15
	}
	return 0
}

func (x *InterfaceStatsEvent) GetLatencyVariance() float64 {
	if x != nil {
		return x.LatencyVariance
	}
	return 0
}

func (x *InterfaceStatsEvent) GetPassiveLatency1() float64 {
	if x != nil {
		return x.PassiveLatency1
	}
	return 0
}

func (x *InterfaceStatsEvent) GetPassiveLatency5() float64 {
	if x != nil {
		return x.PassiveLatency5
	}
	return 0
}

func (x *InterfaceStatsEvent) GetPassiveLatency15() float64 {
	if x != nil {
		return x.PassiveLatency15
	}
	return 0
}

func (x *InterfaceStatsEvent) GetPassiveLatencyVariance() float64 {
	if x != nil {
		return x.PassiveLatencyVariance
	}
	return 0
}

func (x *InterfaceStatsEvent) GetActiveLatency1() float64 {
	if x != nil {
		return x.ActiveLatency1
	}
	return 0
}

func (x *InterfaceStatsEvent) GetActiveLatency5() float64 {
	if x != nil {
		return x.ActiveLatency5
	}
	return 0
}

func (x *InterfaceStatsEvent) GetActiveLatency15() float64 {
	if x != nil {
		return x.ActiveLatency15
	}
	return 0
}

func (x *InterfaceStatsEvent) GetActiveLatencyVariance() float64 {
	if x != nil {
		return x.ActiveLatencyVariance
	}
	return 0
}

func (x *InterfaceStatsEvent) GetJitter1() float64 {
	if x != nil {
		return x.Jitter1
	}
	return 0
}

func (x *InterfaceStatsEvent) GetJitter5() float64 {
	if x != nil {
		return x.Jitter5
	}
	return 0
}

func (x *InterfaceStatsEvent) GetJitter15() float64 {
	if x != nil {
		return x.Jitter15
	}
	return 0
}

func (x *InterfaceStatsEvent) GetJitterVariance() float64 {
	if x != nil {
		return x.JitterVariance
	}
	return 0
}

func (x *InterfaceStatsEvent) GetPingTimeout() uint64 {
	if x != nil {
		return x.PingTimeout
	}
	return 0
}

func (x *InterfaceStatsEvent) GetPingTimeoutRate() uint64 {
	if x != nil {
		return x.PingTimeoutRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxBytes() uint64 {
	if x != nil {
		return x.RxBytes
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxBytesRate() uint64 {
	if x != nil {
		return x.RxBytesRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxPackets() uint64 {
	if x != nil {
		return x.RxPackets
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxPacketsRate() uint64 {
	if x != nil {
		return x.RxPacketsRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxErrs() uint64 {
	if x != nil {
		return x.RxErrs
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxErrsRate() uint64 {
	if x != nil {
		return x.RxErrsRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxDrop() uint64 {
	if x != nil {
		return x.RxDrop
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxDropRate() uint64 {
	if x != nil {
		return x.RxDropRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxFifo() uint64 {
	if x != nil {
		return x.RxFifo
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxFifoRate() uint64 {
	if x != nil {
		return x.RxFifoRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxFrame() uint64 {
	if x != nil {
		return x.RxFrame
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxFrameRate() uint64 {
	if x != nil {
		return x.RxFrameRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxCompressed() uint64 {
	if x != nil {
		return x.RxCompressed
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxCompressedRate() uint64 {
	if x != nil {
		return x.RxCompressedRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxMulticast() uint64 {
	if x != nil {
		return x.RxMulticast
	}
	return 0
}

func (x *InterfaceStatsEvent) GetRxMulticastRate() uint64 {
	if x != nil {
		return x.RxMulticastRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxBytes() uint64 {
	if x != nil {
		return x.TxBytes
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxBytesRate() uint64 {
	if x != nil {
		return x.TxBytesRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxPackets() uint64 {
	if x != nil {
		return x.TxPackets
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxPacketsRate() uint64 {
	if x != nil {
		return x.TxPacketsRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxErrs() uint64 {
	if x != nil {
		return x.TxErrs
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxErrsRate() uint64 {
	if x != nil {
		return x.TxErrsRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxDrop() uint64 {
	if x != nil {
		return x.TxDrop
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxDropRate() uint64 {
	if x != nil {
		return x.TxDropRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxFifo() uint64 {
	if x != nil {
		return x.TxFifo
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxFifoRate() uint64 {
	if x != nil {
		return x.TxFifoRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxColls() uint64 {
	if x != nil {
		return x.TxColls
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxCollsRate() uint64 {
	if x != nil {
		return x.TxCollsRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxCarrier() uint64 {
	if x != nil {
		return x.TxCarrier
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxCarrierRate() uint64 {
	if x != nil {
		return x.TxCarrierRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxCompressed() uint64 {
	if x != nil {
		return x.TxCompressed
	}
	return 0
}

func (x *InterfaceStatsEvent) GetTxCompressedRate() uint64 {
	if x != nil {
		return x.TxCompressedRate
	}
	return 0
}

func (x *InterfaceStatsEvent) GetOffline() bool {
	if x != nil {
		return x.Offline
	}
	return false
}

var File_InterfaceStatsEvent_proto protoreflect.FileDescriptor

var file_InterfaceStatsEvent_proto_rawDesc = []byte{
	0x0a, 0x19, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x73, 0x22, 0xef, 0x0e, 0x0a, 0x13, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x57, 0x61, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x69, 0x73, 0x57, 0x61, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x31, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x31, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x35,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x35,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x35, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x35, 0x12, 0x28,
	0x0a, 0x0f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x63,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x61, 0x73, 0x73,
	0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0f, 0x70, 0x61, 0x73, 0x73, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63,
	0x79, 0x31, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x61, 0x73, 0x73, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74,
	0x65, 0x6e, 0x63, 0x79, 0x35, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x70, 0x61, 0x73,
	0x73, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x35, 0x12, 0x2a, 0x0a, 0x10,
	0x70, 0x61, 0x73, 0x73, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x35,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x70, 0x61, 0x73, 0x73, 0x69, 0x76, 0x65, 0x4c,
	0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x35, 0x12, 0x36, 0x0a, 0x16, 0x70, 0x61, 0x73, 0x73,
	0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e,
	0x63, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x01, 0x52, 0x16, 0x70, 0x61, 0x73, 0x73, 0x69, 0x76,
	0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x63, 0x65,
	0x12, 0x26, 0x0a, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63,
	0x79, 0x31, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x12, 0x26, 0x0a, 0x0e, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x35, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x35,
	0x12, 0x28, 0x0a, 0x0f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63,
	0x79, 0x31, 0x35, 0x18, 0x10, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x31, 0x35, 0x12, 0x34, 0x0a, 0x15, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x56, 0x61, 0x72, 0x69, 0x61,
	0x6e, 0x63, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x01, 0x52, 0x15, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x63, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x31, 0x18, 0x12, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x07, 0x6a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x31, 0x12, 0x18, 0x0a, 0x07, 0x6a, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x35, 0x18, 0x13, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x6a, 0x69, 0x74,
	0x74, 0x65, 0x72, 0x35, 0x12, 0x1a, 0x0a, 0x08, 0x6a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x31, 0x35,
	0x18, 0x14, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x31, 0x35,
	0x12, 0x26, 0x0a, 0x0e, 0x6a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e,
	0x63, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x6a, 0x69, 0x74, 0x74, 0x65, 0x72,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x69, 0x6e, 0x67,
	0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x16, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x70,
	0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x69,
	0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x52, 0x61, 0x74, 0x65, 0x18, 0x17, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0f, 0x70, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x18, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x20,
	0x0a, 0x0b, 0x72, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x61, 0x74, 0x65, 0x18, 0x19, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x72, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x72, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x1a, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x72, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x24,
	0x0a, 0x0d, 0x72, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x52, 0x61, 0x74, 0x65, 0x18,
	0x1b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x72, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x78, 0x45, 0x72, 0x72, 0x73, 0x18, 0x1c,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72, 0x78, 0x45, 0x72, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a,
	0x72, 0x78, 0x45, 0x72, 0x72, 0x73, 0x52, 0x61, 0x74, 0x65, 0x18, 0x1d, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x72, 0x78, 0x45, 0x72, 0x72, 0x73, 0x52, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x72, 0x78, 0x44, 0x72, 0x6f, 0x70, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72, 0x78,
	0x44, 0x72, 0x6f, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x78, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x61,
	0x74, 0x65, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x72, 0x78, 0x44, 0x72, 0x6f, 0x70,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x18, 0x20,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x12, 0x1e, 0x0a, 0x0a,
	0x72, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x52, 0x61, 0x74, 0x65, 0x18, 0x21, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x72, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x72, 0x78, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x18, 0x22, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72,
	0x78, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x78, 0x46, 0x72, 0x61, 0x6d,
	0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x23, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x72, 0x78, 0x46,
	0x72, 0x61, 0x6d, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x78, 0x43, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x18, 0x24, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c,
	0x72, 0x78, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x12, 0x2a, 0x0a, 0x10,
	0x72, 0x78, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x52, 0x61, 0x74, 0x65,
	0x18, 0x25, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x72, 0x78, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x64, 0x52, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x78, 0x4d, 0x75,
	0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x18, 0x26, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x72,
	0x78, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x72, 0x78,
	0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x52, 0x61, 0x74, 0x65, 0x18, 0x27, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0f, 0x72, 0x78, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x28, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x74, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x20,
	0x0a, 0x0b, 0x74, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x61, 0x74, 0x65, 0x18, 0x29, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x74, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x2a, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x24,
	0x0a, 0x0d, 0x74, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x52, 0x61, 0x74, 0x65, 0x18,
	0x2b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x74, 0x78, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x78, 0x45, 0x72, 0x72, 0x73, 0x18, 0x2c,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x74, 0x78, 0x45, 0x72, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a,
	0x74, 0x78, 0x45, 0x72, 0x72, 0x73, 0x52, 0x61, 0x74, 0x65, 0x18, 0x2d, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x74, 0x78, 0x45, 0x72, 0x72, 0x73, 0x52, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x74, 0x78, 0x44, 0x72, 0x6f, 0x70, 0x18, 0x2e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x74, 0x78,
	0x44, 0x72, 0x6f, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x78, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x61,
	0x74, 0x65, 0x18, 0x2f, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x78, 0x44, 0x72, 0x6f, 0x70,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x18, 0x30,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x74, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x12, 0x1e, 0x0a, 0x0a,
	0x74, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x52, 0x61, 0x74, 0x65, 0x18, 0x31, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x74, 0x78, 0x46, 0x69, 0x66, 0x6f, 0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x74, 0x78, 0x43, 0x6f, 0x6c, 0x6c, 0x73, 0x18, 0x32, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x74,
	0x78, 0x43, 0x6f, 0x6c, 0x6c, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x78, 0x43, 0x6f, 0x6c, 0x6c,
	0x73, 0x52, 0x61, 0x74, 0x65, 0x18, 0x33, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x74, 0x78, 0x43,
	0x6f, 0x6c, 0x6c, 0x73, 0x52, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x78, 0x43, 0x61,
	0x72, 0x72, 0x69, 0x65, 0x72, 0x18, 0x34, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x78, 0x43,
	0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x74, 0x78, 0x43, 0x61, 0x72, 0x72,
	0x69, 0x65, 0x72, 0x52, 0x61, 0x74, 0x65, 0x18, 0x35, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x74,
	0x78, 0x43, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x52, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0c,
	0x74, 0x78, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x18, 0x36, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0c, 0x74, 0x78, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64,
	0x12, 0x2a, 0x0a, 0x10, 0x74, 0x78, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64,
	0x52, 0x61, 0x74, 0x65, 0x18, 0x37, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x74, 0x78, 0x43, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x52, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x38, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6f,
	0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x42, 0x4f, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x2f, 0x67, 0x6f,
	0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66,
	0x65, 0x72, 0x73, 0x2f, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_InterfaceStatsEvent_proto_rawDescOnce sync.Once
	file_InterfaceStatsEvent_proto_rawDescData = file_InterfaceStatsEvent_proto_rawDesc
)

func file_InterfaceStatsEvent_proto_rawDescGZIP() []byte {
	file_InterfaceStatsEvent_proto_rawDescOnce.Do(func() {
		file_InterfaceStatsEvent_proto_rawDescData = protoimpl.X.CompressGZIP(file_InterfaceStatsEvent_proto_rawDescData)
	})
	return file_InterfaceStatsEvent_proto_rawDescData
}

var file_InterfaceStatsEvent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_InterfaceStatsEvent_proto_goTypes = []interface{}{
	(*InterfaceStatsEvent)(nil), // 0: reports.InterfaceStatsEvent
}
var file_InterfaceStatsEvent_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_InterfaceStatsEvent_proto_init() }
func file_InterfaceStatsEvent_proto_init() {
	if File_InterfaceStatsEvent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_InterfaceStatsEvent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterfaceStatsEvent); i {
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
			RawDescriptor: file_InterfaceStatsEvent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_InterfaceStatsEvent_proto_goTypes,
		DependencyIndexes: file_InterfaceStatsEvent_proto_depIdxs,
		MessageInfos:      file_InterfaceStatsEvent_proto_msgTypes,
	}.Build()
	File_InterfaceStatsEvent_proto = out.File
	file_InterfaceStatsEvent_proto_rawDesc = nil
	file_InterfaceStatsEvent_proto_goTypes = nil
	file_InterfaceStatsEvent_proto_depIdxs = nil
}

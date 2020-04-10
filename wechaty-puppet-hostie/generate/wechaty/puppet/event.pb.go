// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package puppet

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EventType int32

const (
	EventType_EVENT_TYPE_UNSPECIFIED EventType = 0
	EventType_EVENT_TYPE_HEARTBEAT   EventType = 1
	EventType_EVENT_TYPE_MESSAGE     EventType = 2
	EventType_EVENT_TYPE_DONG        EventType = 3
	EventType_EVENT_TYPE_ERROR       EventType = 16
	EventType_EVENT_TYPE_FRIENDSHIP  EventType = 17
	EventType_EVENT_TYPE_ROOM_INVITE EventType = 18
	EventType_EVENT_TYPE_ROOM_JOIN   EventType = 19
	EventType_EVENT_TYPE_ROOM_LEAVE  EventType = 20
	EventType_EVENT_TYPE_ROOM_TOPIC  EventType = 21
	EventType_EVENT_TYPE_SCAN        EventType = 22
	EventType_EVENT_TYPE_READY       EventType = 23
	EventType_EVENT_TYPE_RESET       EventType = 24
	EventType_EVENT_TYPE_LOGIN       EventType = 25
	EventType_EVENT_TYPE_LOGOUT      EventType = 26
)

var EventType_name = map[int32]string{
	0:  "EVENT_TYPE_UNSPECIFIED",
	1:  "EVENT_TYPE_HEARTBEAT",
	2:  "EVENT_TYPE_MESSAGE",
	3:  "EVENT_TYPE_DONG",
	16: "EVENT_TYPE_ERROR",
	17: "EVENT_TYPE_FRIENDSHIP",
	18: "EVENT_TYPE_ROOM_INVITE",
	19: "EVENT_TYPE_ROOM_JOIN",
	20: "EVENT_TYPE_ROOM_LEAVE",
	21: "EVENT_TYPE_ROOM_TOPIC",
	22: "EVENT_TYPE_SCAN",
	23: "EVENT_TYPE_READY",
	24: "EVENT_TYPE_RESET",
	25: "EVENT_TYPE_LOGIN",
	26: "EVENT_TYPE_LOGOUT",
}

var EventType_value = map[string]int32{
	"EVENT_TYPE_UNSPECIFIED": 0,
	"EVENT_TYPE_HEARTBEAT":   1,
	"EVENT_TYPE_MESSAGE":     2,
	"EVENT_TYPE_DONG":        3,
	"EVENT_TYPE_ERROR":       16,
	"EVENT_TYPE_FRIENDSHIP":  17,
	"EVENT_TYPE_ROOM_INVITE": 18,
	"EVENT_TYPE_ROOM_JOIN":   19,
	"EVENT_TYPE_ROOM_LEAVE":  20,
	"EVENT_TYPE_ROOM_TOPIC":  21,
	"EVENT_TYPE_SCAN":        22,
	"EVENT_TYPE_READY":       23,
	"EVENT_TYPE_RESET":       24,
	"EVENT_TYPE_LOGIN":       25,
	"EVENT_TYPE_LOGOUT":      26,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}

func (EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

type EventRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventRequest) Reset()         { *m = EventRequest{} }
func (m *EventRequest) String() string { return proto.CompactTextString(m) }
func (*EventRequest) ProtoMessage()    {}
func (*EventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

func (m *EventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventRequest.Unmarshal(m, b)
}
func (m *EventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventRequest.Marshal(b, m, deterministic)
}
func (m *EventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventRequest.Merge(m, src)
}
func (m *EventRequest) XXX_Size() int {
	return xxx_messageInfo_EventRequest.Size(m)
}
func (m *EventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EventRequest proto.InternalMessageInfo

type EventResponse struct {
	Type EventType `protobuf:"varint,1,opt,name=type,proto3,enum=wechaty.puppet.EventType" json:"type,omitempty"`
	// TODO: Huan(202002) consider to use a PB Map?
	Payload              string   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventResponse) Reset()         { *m = EventResponse{} }
func (m *EventResponse) String() string { return proto.CompactTextString(m) }
func (*EventResponse) ProtoMessage()    {}
func (*EventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{1}
}

func (m *EventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventResponse.Unmarshal(m, b)
}
func (m *EventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventResponse.Marshal(b, m, deterministic)
}
func (m *EventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventResponse.Merge(m, src)
}
func (m *EventResponse) XXX_Size() int {
	return xxx_messageInfo_EventResponse.Size(m)
}
func (m *EventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EventResponse proto.InternalMessageInfo

func (m *EventResponse) GetType() EventType {
	if m != nil {
		return m.Type
	}
	return EventType_EVENT_TYPE_UNSPECIFIED
}

func (m *EventResponse) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func init() {
	proto.RegisterEnum("wechaty.puppet.EventType", EventType_name, EventType_value)
	proto.RegisterType((*EventRequest)(nil), "wechaty.puppet.EventRequest")
	proto.RegisterType((*EventResponse)(nil), "wechaty.puppet.EventResponse")
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_2d17a9d3f0ddf27e) }

var fileDescriptor_2d17a9d3f0ddf27e = []byte{
	// 366 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xdf, 0x8f, 0xd2, 0x40,
	0x10, 0xc7, 0x2d, 0x12, 0xf5, 0x56, 0xe5, 0xe6, 0xf6, 0x00, 0xcb, 0x3d, 0x5d, 0x78, 0xba, 0x98,
	0xb4, 0x24, 0xfa, 0x17, 0xf4, 0xe8, 0x1c, 0xb7, 0x0a, 0xbb, 0xcd, 0x76, 0x21, 0xe2, 0x0b, 0x29,
	0xb8, 0x01, 0x12, 0xa5, 0x2b, 0x5d, 0x34, 0xfd, 0xeb, 0x35, 0x96, 0x1f, 0x29, 0xf5, 0x9e, 0x76,
	0xe7, 0xf3, 0x99, 0xd9, 0x7c, 0x93, 0x59, 0xf2, 0x5a, 0xff, 0xd2, 0x1b, 0xeb, 0x9b, 0x6d, 0x6a,
	0x53, 0xda, 0xf8, 0xad, 0x17, 0xab, 0xc4, 0xe6, 0xbe, 0xd9, 0x19, 0xa3, 0x6d, 0xb7, 0x41, 0xde,
	0xe0, 0x3f, 0x2d, 0xf5, 0xcf, 0x9d, 0xce, 0x6c, 0xf7, 0x0b, 0x79, 0x7b, 0xa8, 0x33, 0x93, 0x6e,
	0x32, 0x4d, 0x3d, 0x52, 0xb7, 0xb9, 0xd1, 0xae, 0x73, 0xeb, 0xdc, 0x35, 0x3e, 0x74, 0xfc, 0xf3,
	0x79, 0xbf, 0x68, 0x56, 0xb9, 0xd1, 0xb2, 0x68, 0xa3, 0x2e, 0x79, 0x69, 0x92, 0xfc, 0x7b, 0x9a,
	0x7c, 0x73, 0x6b, 0xb7, 0xce, 0xdd, 0x85, 0x3c, 0x96, 0xef, 0xff, 0xd4, 0xc8, 0xc5, 0xa9, 0x9b,
	0xde, 0x90, 0x36, 0x4e, 0x90, 0xab, 0x99, 0x9a, 0x46, 0x38, 0x1b, 0xf3, 0x38, 0xc2, 0x3e, 0x7b,
	0x60, 0x18, 0xc2, 0x33, 0xea, 0x92, 0x66, 0xc9, 0x3d, 0x62, 0x20, 0xd5, 0x3d, 0x06, 0x0a, 0x1c,
	0xda, 0x26, 0xb4, 0x64, 0x46, 0x18, 0xc7, 0xc1, 0x00, 0xa1, 0x46, 0xaf, 0xc9, 0x65, 0x89, 0x87,
	0x82, 0x0f, 0xe0, 0x39, 0x6d, 0x12, 0x28, 0x41, 0x94, 0x52, 0x48, 0x00, 0xda, 0x21, 0xad, 0x12,
	0x7d, 0x90, 0x0c, 0x79, 0x18, 0x3f, 0xb2, 0x08, 0xae, 0x2a, 0x99, 0xa4, 0x10, 0xa3, 0x19, 0xe3,
	0x13, 0xa6, 0x10, 0x68, 0x25, 0x53, 0xe1, 0x3e, 0x09, 0xc6, 0xe1, 0xba, 0xf2, 0x60, 0x61, 0x86,
	0x18, 0x4c, 0x10, 0x9a, 0x4f, 0x29, 0x25, 0x22, 0xd6, 0x87, 0x56, 0x25, 0x71, 0xdc, 0x0f, 0x38,
	0xb4, 0x2b, 0x89, 0x25, 0x06, 0xe1, 0x14, 0xde, 0xfd, 0x47, 0x63, 0x54, 0xe0, 0x56, 0xe8, 0x50,
	0x0c, 0x18, 0x87, 0x0e, 0x6d, 0x91, 0xab, 0x73, 0x2a, 0xc6, 0x0a, 0x6e, 0xba, 0xf5, 0x57, 0x75,
	0xb8, 0xbc, 0x1f, 0x7d, 0xfd, 0xbc, 0x5c, 0xdb, 0xd5, 0x6e, 0xee, 0x2f, 0xd2, 0x1f, 0xbd, 0xc3,
	0x22, 0x7b, 0xcb, 0xd4, 0x3b, 0x5e, 0x0f, 0xa7, 0xb7, 0xdf, 0xad, 0xb7, 0x4a, 0x33, 0xbb, 0xd6,
	0xbd, 0xa5, 0xde, 0xe8, 0x6d, 0x62, 0xf5, 0x69, 0x62, 0xaf, 0xe7, 0x2f, 0x8a, 0x1f, 0xf5, 0xf1,
	0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x83, 0x86, 0xfe, 0x58, 0x60, 0x02, 0x00, 0x00,
}

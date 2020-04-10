// Code generated by protoc-gen-go. DO NOT EDIT.
// source: room_member.proto

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

type RoomMemberPayloadRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MemberId             string   `protobuf:"bytes,2,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomMemberPayloadRequest) Reset()         { *m = RoomMemberPayloadRequest{} }
func (m *RoomMemberPayloadRequest) String() string { return proto.CompactTextString(m) }
func (*RoomMemberPayloadRequest) ProtoMessage()    {}
func (*RoomMemberPayloadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_63e6562082d7765e, []int{0}
}

func (m *RoomMemberPayloadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomMemberPayloadRequest.Unmarshal(m, b)
}
func (m *RoomMemberPayloadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomMemberPayloadRequest.Marshal(b, m, deterministic)
}
func (m *RoomMemberPayloadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomMemberPayloadRequest.Merge(m, src)
}
func (m *RoomMemberPayloadRequest) XXX_Size() int {
	return xxx_messageInfo_RoomMemberPayloadRequest.Size(m)
}
func (m *RoomMemberPayloadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomMemberPayloadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RoomMemberPayloadRequest proto.InternalMessageInfo

func (m *RoomMemberPayloadRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RoomMemberPayloadRequest) GetMemberId() string {
	if m != nil {
		return m.MemberId
	}
	return ""
}

type RoomMemberPayloadResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RoomAlias            string   `protobuf:"bytes,2,opt,name=room_alias,json=roomAlias,proto3" json:"room_alias,omitempty"`
	InviterId            string   `protobuf:"bytes,3,opt,name=inviter_id,json=inviterId,proto3" json:"inviter_id,omitempty"`
	Avatar               string   `protobuf:"bytes,4,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Name                 string   `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomMemberPayloadResponse) Reset()         { *m = RoomMemberPayloadResponse{} }
func (m *RoomMemberPayloadResponse) String() string { return proto.CompactTextString(m) }
func (*RoomMemberPayloadResponse) ProtoMessage()    {}
func (*RoomMemberPayloadResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_63e6562082d7765e, []int{1}
}

func (m *RoomMemberPayloadResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomMemberPayloadResponse.Unmarshal(m, b)
}
func (m *RoomMemberPayloadResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomMemberPayloadResponse.Marshal(b, m, deterministic)
}
func (m *RoomMemberPayloadResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomMemberPayloadResponse.Merge(m, src)
}
func (m *RoomMemberPayloadResponse) XXX_Size() int {
	return xxx_messageInfo_RoomMemberPayloadResponse.Size(m)
}
func (m *RoomMemberPayloadResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomMemberPayloadResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RoomMemberPayloadResponse proto.InternalMessageInfo

func (m *RoomMemberPayloadResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RoomMemberPayloadResponse) GetRoomAlias() string {
	if m != nil {
		return m.RoomAlias
	}
	return ""
}

func (m *RoomMemberPayloadResponse) GetInviterId() string {
	if m != nil {
		return m.InviterId
	}
	return ""
}

func (m *RoomMemberPayloadResponse) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *RoomMemberPayloadResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RoomMemberListRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomMemberListRequest) Reset()         { *m = RoomMemberListRequest{} }
func (m *RoomMemberListRequest) String() string { return proto.CompactTextString(m) }
func (*RoomMemberListRequest) ProtoMessage()    {}
func (*RoomMemberListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_63e6562082d7765e, []int{2}
}

func (m *RoomMemberListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomMemberListRequest.Unmarshal(m, b)
}
func (m *RoomMemberListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomMemberListRequest.Marshal(b, m, deterministic)
}
func (m *RoomMemberListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomMemberListRequest.Merge(m, src)
}
func (m *RoomMemberListRequest) XXX_Size() int {
	return xxx_messageInfo_RoomMemberListRequest.Size(m)
}
func (m *RoomMemberListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomMemberListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RoomMemberListRequest proto.InternalMessageInfo

func (m *RoomMemberListRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type RoomMemberListResponse struct {
	MemberIds            []string `protobuf:"bytes,1,rep,name=member_ids,json=memberIds,proto3" json:"member_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomMemberListResponse) Reset()         { *m = RoomMemberListResponse{} }
func (m *RoomMemberListResponse) String() string { return proto.CompactTextString(m) }
func (*RoomMemberListResponse) ProtoMessage()    {}
func (*RoomMemberListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_63e6562082d7765e, []int{3}
}

func (m *RoomMemberListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomMemberListResponse.Unmarshal(m, b)
}
func (m *RoomMemberListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomMemberListResponse.Marshal(b, m, deterministic)
}
func (m *RoomMemberListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomMemberListResponse.Merge(m, src)
}
func (m *RoomMemberListResponse) XXX_Size() int {
	return xxx_messageInfo_RoomMemberListResponse.Size(m)
}
func (m *RoomMemberListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomMemberListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RoomMemberListResponse proto.InternalMessageInfo

func (m *RoomMemberListResponse) GetMemberIds() []string {
	if m != nil {
		return m.MemberIds
	}
	return nil
}

func init() {
	proto.RegisterType((*RoomMemberPayloadRequest)(nil), "wechaty.puppet.RoomMemberPayloadRequest")
	proto.RegisterType((*RoomMemberPayloadResponse)(nil), "wechaty.puppet.RoomMemberPayloadResponse")
	proto.RegisterType((*RoomMemberListRequest)(nil), "wechaty.puppet.RoomMemberListRequest")
	proto.RegisterType((*RoomMemberListResponse)(nil), "wechaty.puppet.RoomMemberListResponse")
}

func init() { proto.RegisterFile("room_member.proto", fileDescriptor_63e6562082d7765e) }

var fileDescriptor_63e6562082d7765e = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x3f, 0x6b, 0xc3, 0x30,
	0x10, 0xc5, 0x71, 0x92, 0x86, 0xfa, 0x86, 0x40, 0x05, 0x0d, 0x2a, 0x25, 0x10, 0xbc, 0xb4, 0x8b,
	0xed, 0xa1, 0x43, 0xe7, 0x76, 0x29, 0xa1, 0x0d, 0x14, 0x8f, 0x5d, 0x82, 0x1c, 0x1d, 0xb6, 0x20,
	0xf2, 0xb9, 0x96, 0x9c, 0x92, 0x0f, 0xd2, 0xef, 0x5b, 0x2c, 0x29, 0x29, 0xf4, 0xcf, 0xe4, 0xf3,
	0xef, 0x9d, 0x9e, 0xde, 0x9d, 0xe0, 0xa2, 0x23, 0xd2, 0x1b, 0x8d, 0xba, 0xc4, 0x2e, 0x6b, 0x3b,
	0xb2, 0xc4, 0x66, 0x1f, 0xb8, 0xad, 0x85, 0x3d, 0x64, 0x6d, 0xdf, 0xb6, 0x68, 0x93, 0x27, 0xe0,
	0x05, 0x91, 0x5e, 0xbb, 0x9e, 0x57, 0x71, 0xd8, 0x91, 0x90, 0x05, 0xbe, 0xf7, 0x68, 0x2c, 0x9b,
	0xc1, 0x48, 0x49, 0x1e, 0x2d, 0xa3, 0xdb, 0xb8, 0x18, 0x29, 0xc9, 0xae, 0x21, 0xf6, 0x5e, 0x1b,
	0x25, 0xf9, 0xc8, 0xe1, 0x73, 0x0f, 0x56, 0x32, 0xf9, 0x8c, 0xe0, 0xea, 0x0f, 0x27, 0xd3, 0x52,
	0x63, 0xf0, 0x97, 0xd5, 0x02, 0xc0, 0x65, 0x13, 0x3b, 0x25, 0x4c, 0xf0, 0x8a, 0x07, 0xf2, 0x30,
	0x80, 0x41, 0x56, 0xcd, 0x5e, 0x59, 0x7f, 0xd5, 0xd8, 0xcb, 0x81, 0xac, 0x24, 0x9b, 0xc3, 0x54,
	0xec, 0x85, 0x15, 0x1d, 0x9f, 0x38, 0x29, 0xfc, 0x31, 0x06, 0x93, 0x46, 0x68, 0xe4, 0x67, 0x8e,
	0xba, 0x3a, 0xb9, 0x81, 0xcb, 0xef, 0x58, 0x2f, 0xca, 0xd8, 0x7f, 0xa6, 0x4b, 0xee, 0x61, 0xfe,
	0xb3, 0x31, 0x84, 0x5f, 0x00, 0x9c, 0xe6, 0x36, 0x3c, 0x5a, 0x8e, 0x87, 0x34, 0xc7, 0xc1, 0xcd,
	0xe3, 0xfa, 0xed, 0xb9, 0x52, 0xb6, 0xee, 0xcb, 0x6c, 0x4b, 0x3a, 0x0f, 0xfb, 0xcd, 0x2b, 0x4a,
	0x8f, 0x65, 0xf8, 0xa6, 0x7e, 0xe5, 0x69, 0x4d, 0xc6, 0x2a, 0xcc, 0x2b, 0x6c, 0xb0, 0x13, 0x16,
	0x4f, 0x27, 0xbc, 0x5c, 0x4e, 0xdd, 0x43, 0xdd, 0x7d, 0x05, 0x00, 0x00, 0xff, 0xff, 0x62, 0x88,
	0xbb, 0xf9, 0xbd, 0x01, 0x00, 0x00,
}

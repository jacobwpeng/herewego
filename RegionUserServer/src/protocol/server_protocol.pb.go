// Code generated by protoc-gen-go.
// source: server_protocol.proto
// DO NOT EDIT!

/*
Package regionserverprotocol is a generated protocol buffer package.

It is generated from these files:
	server_protocol.proto

It has these top-level messages:
	Message
	UpdateRegionRequest
	PickFromRegionRequest
	PickFromRegionResponse
	RegionStatus
	RegionStatusRequest
	RegionStatusReponse
*/
package regionserverprotocol

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type MessageCommand int32

const (
	MessageCommand_UpdateRegionRequestCmd    MessageCommand = 1000
	MessageCommand_PickFromRegionRequestCmd  MessageCommand = 1010
	MessageCommand_PickFromRegionResponseCmd MessageCommand = 1011
	MessageCommand_RegionStatusRequestCmd    MessageCommand = 1020
	MessageCommand_RegionStatusReponseCmd    MessageCommand = 1021
)

var MessageCommand_name = map[int32]string{
	1000: "UpdateRegionRequestCmd",
	1010: "PickFromRegionRequestCmd",
	1011: "PickFromRegionResponseCmd",
	1020: "RegionStatusRequestCmd",
	1021: "RegionStatusReponseCmd",
}
var MessageCommand_value = map[string]int32{
	"UpdateRegionRequestCmd":    1000,
	"PickFromRegionRequestCmd":  1010,
	"PickFromRegionResponseCmd": 1011,
	"RegionStatusRequestCmd":    1020,
	"RegionStatusReponseCmd":    1021,
}

func (x MessageCommand) Enum() *MessageCommand {
	p := new(MessageCommand)
	*p = x
	return p
}
func (x MessageCommand) String() string {
	return proto.EnumName(MessageCommand_name, int32(x))
}
func (x *MessageCommand) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageCommand_value, data, "MessageCommand")
	if err != nil {
		return err
	}
	*x = MessageCommand(value)
	return nil
}
func (MessageCommand) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Message struct {
	Cmd              *uint32 `protobuf:"varint,1,opt,name=cmd" json:"cmd,omitempty"`
	Ctx              *uint32 `protobuf:"varint,2,opt,name=ctx" json:"ctx,omitempty"`
	Payload          []byte  `protobuf:"bytes,3,opt,name=payload" json:"payload,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetCmd() uint32 {
	if m != nil && m.Cmd != nil {
		return *m.Cmd
	}
	return 0
}

func (m *Message) GetCtx() uint32 {
	if m != nil && m.Ctx != nil {
		return *m.Ctx
	}
	return 0
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type UpdateRegionRequest struct {
	Uin              *uint32 `protobuf:"varint,1,opt,name=uin" json:"uin,omitempty"`
	Level            *uint32 `protobuf:"varint,2,opt,name=level" json:"level,omitempty"`
	Region           *uint32 `protobuf:"varint,3,opt,name=region" json:"region,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *UpdateRegionRequest) Reset()                    { *m = UpdateRegionRequest{} }
func (m *UpdateRegionRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateRegionRequest) ProtoMessage()               {}
func (*UpdateRegionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UpdateRegionRequest) GetUin() uint32 {
	if m != nil && m.Uin != nil {
		return *m.Uin
	}
	return 0
}

func (m *UpdateRegionRequest) GetLevel() uint32 {
	if m != nil && m.Level != nil {
		return *m.Level
	}
	return 0
}

func (m *UpdateRegionRequest) GetRegion() uint32 {
	if m != nil && m.Region != nil {
		return *m.Region
	}
	return 0
}

type PickFromRegionRequest struct {
	SelfUin          *uint32 `protobuf:"varint,1,opt,name=self_uin" json:"self_uin,omitempty"`
	SelfLevel        *uint32 `protobuf:"varint,2,opt,name=self_level" json:"self_level,omitempty"`
	ExpectRegion     *uint32 `protobuf:"varint,3,opt,name=expect_region" json:"expect_region,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PickFromRegionRequest) Reset()                    { *m = PickFromRegionRequest{} }
func (m *PickFromRegionRequest) String() string            { return proto.CompactTextString(m) }
func (*PickFromRegionRequest) ProtoMessage()               {}
func (*PickFromRegionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PickFromRegionRequest) GetSelfUin() uint32 {
	if m != nil && m.SelfUin != nil {
		return *m.SelfUin
	}
	return 0
}

func (m *PickFromRegionRequest) GetSelfLevel() uint32 {
	if m != nil && m.SelfLevel != nil {
		return *m.SelfLevel
	}
	return 0
}

func (m *PickFromRegionRequest) GetExpectRegion() uint32 {
	if m != nil && m.ExpectRegion != nil {
		return *m.ExpectRegion
	}
	return 0
}

type PickFromRegionResponse struct {
	Uin              []uint32 `protobuf:"varint,1,rep,name=uin" json:"uin,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PickFromRegionResponse) Reset()                    { *m = PickFromRegionResponse{} }
func (m *PickFromRegionResponse) String() string            { return proto.CompactTextString(m) }
func (*PickFromRegionResponse) ProtoMessage()               {}
func (*PickFromRegionResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PickFromRegionResponse) GetUin() []uint32 {
	if m != nil {
		return m.Uin
	}
	return nil
}

type RegionStatus struct {
	Region           *uint32 `protobuf:"varint,1,opt,name=region" json:"region,omitempty"`
	Num              *uint32 `protobuf:"varint,2,opt,name=num" json:"num,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RegionStatus) Reset()                    { *m = RegionStatus{} }
func (m *RegionStatus) String() string            { return proto.CompactTextString(m) }
func (*RegionStatus) ProtoMessage()               {}
func (*RegionStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RegionStatus) GetRegion() uint32 {
	if m != nil && m.Region != nil {
		return *m.Region
	}
	return 0
}

func (m *RegionStatus) GetNum() uint32 {
	if m != nil && m.Num != nil {
		return *m.Num
	}
	return 0
}

type RegionStatusRequest struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RegionStatusRequest) Reset()                    { *m = RegionStatusRequest{} }
func (m *RegionStatusRequest) String() string            { return proto.CompactTextString(m) }
func (*RegionStatusRequest) ProtoMessage()               {}
func (*RegionStatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type RegionStatusReponse struct {
	Status           []*RegionStatus `protobuf:"bytes,1,rep,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *RegionStatusReponse) Reset()                    { *m = RegionStatusReponse{} }
func (m *RegionStatusReponse) String() string            { return proto.CompactTextString(m) }
func (*RegionStatusReponse) ProtoMessage()               {}
func (*RegionStatusReponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RegionStatusReponse) GetStatus() []*RegionStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "regionserverprotocol.Message")
	proto.RegisterType((*UpdateRegionRequest)(nil), "regionserverprotocol.UpdateRegionRequest")
	proto.RegisterType((*PickFromRegionRequest)(nil), "regionserverprotocol.PickFromRegionRequest")
	proto.RegisterType((*PickFromRegionResponse)(nil), "regionserverprotocol.PickFromRegionResponse")
	proto.RegisterType((*RegionStatus)(nil), "regionserverprotocol.RegionStatus")
	proto.RegisterType((*RegionStatusRequest)(nil), "regionserverprotocol.RegionStatusRequest")
	proto.RegisterType((*RegionStatusReponse)(nil), "regionserverprotocol.RegionStatusReponse")
	proto.RegisterEnum("regionserverprotocol.MessageCommand", MessageCommand_name, MessageCommand_value)
}

var fileDescriptor0 = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x91, 0xc1, 0x4a, 0xfb, 0x40,
	0x10, 0xc6, 0xe9, 0xbf, 0xfc, 0xbb, 0x32, 0x6d, 0x6a, 0xd9, 0x9a, 0x12, 0x11, 0x45, 0x16, 0x04,
	0x51, 0xe8, 0xa1, 0xf8, 0x02, 0x52, 0x10, 0x3c, 0x08, 0x52, 0xf5, 0x1c, 0x96, 0x64, 0x2c, 0xc1,
	0x6c, 0x36, 0x66, 0x37, 0xa5, 0xbe, 0x90, 0xcf, 0xe6, 0x59, 0xaf, 0x0a, 0x76, 0x37, 0xb1, 0x4d,
	0xcb, 0xde, 0xb2, 0xf3, 0xcd, 0xfc, 0xe6, 0x9b, 0x2f, 0xe0, 0x2b, 0x2c, 0x16, 0x58, 0x84, 0x79,
	0x21, 0xb5, 0x8c, 0x64, 0x3a, 0xb6, 0x1f, 0xf4, 0xa0, 0xc0, 0x79, 0x22, 0xb3, 0x4a, 0xfc, 0xd3,
	0xd8, 0x15, 0x90, 0x3b, 0x54, 0x8a, 0xcf, 0x91, 0x76, 0xa1, 0x1d, 0x89, 0x38, 0x68, 0x9d, 0xb6,
	0xce, 0x3d, 0xfb, 0xd0, 0xcb, 0xe0, 0x9f, 0x7d, 0xec, 0x03, 0xc9, 0xf9, 0x5b, 0x2a, 0x79, 0x1c,
	0xb4, 0x57, 0x85, 0x1e, 0xbb, 0x86, 0xe1, 0x53, 0x1e, 0x73, 0x8d, 0x33, 0xcb, 0x9c, 0xe1, 0x6b,
	0x89, 0x4a, 0x9b, 0xa1, 0x32, 0xc9, 0x6a, 0x82, 0x07, 0xff, 0x53, 0x5c, 0x60, 0x5a, 0x33, 0xfa,
	0xd0, 0xa9, 0x0c, 0x58, 0x84, 0xc7, 0x1e, 0xc1, 0xbf, 0x4f, 0xa2, 0x97, 0x9b, 0x42, 0x8a, 0x6d,
	0xc8, 0x00, 0xf6, 0x14, 0xa6, 0xcf, 0xe1, 0x86, 0x44, 0x01, 0x6c, 0xa5, 0x89, 0xf3, 0xc1, 0xc3,
	0x65, 0x8e, 0x91, 0x0e, 0xb7, 0xa8, 0x67, 0x30, 0xda, 0xa5, 0xaa, 0xdc, 0x1c, 0xbd, 0xf1, 0xd6,
	0x5e, 0xb5, 0x5d, 0x42, 0xaf, 0x92, 0x1f, 0x34, 0xd7, 0xa5, 0x6a, 0x98, 0x5b, 0x5f, 0x9f, 0x95,
	0xa2, 0x5a, 0xc5, 0x7c, 0x18, 0x36, 0x9b, 0x6b, 0x9f, 0xec, 0x76, 0xb7, 0x5c, 0xed, 0x99, 0x40,
	0x47, 0xd9, 0x82, 0x5d, 0xd5, 0x9d, 0xb0, 0xb1, 0x2b, 0xf7, 0x71, 0x73, 0xf4, 0xe2, 0xbd, 0x05,
	0xfd, 0xfa, 0x2f, 0x4c, 0xa5, 0x10, 0x3c, 0x8b, 0xe9, 0x11, 0x8c, 0x1c, 0x09, 0x4f, 0x45, 0x3c,
	0xf8, 0x20, 0xf4, 0x18, 0x02, 0x67, 0x76, 0x46, 0xfe, 0x24, 0xf4, 0x04, 0x0e, 0xdd, 0x21, 0x18,
	0xfd, 0x8b, 0x18, 0xb6, 0xe3, 0x20, 0x23, 0x7e, 0x3b, 0xc4, 0xf5, 0xe4, 0x0f, 0xf9, 0x0d, 0x00,
	0x00, 0xff, 0xff, 0x64, 0x82, 0xa1, 0x9c, 0x5b, 0x02, 0x00, 0x00,
}
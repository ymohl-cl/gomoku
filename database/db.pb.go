// Code generated by protoc-gen-go. DO NOT EDIT.
// source: db.proto

/*
Package database is a generated protocol buffer package.

It is generated from these files:
	db.proto

It has these top-level messages:
	Player
	Session
	Stat
	Move
	Data
*/
package database

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Player struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Stats *Stat  `protobuf:"bytes,2,opt,name=stats" json:"stats,omitempty"`
	Saves int32  `protobuf:"varint,3,opt,name=saves" json:"saves,omitempty"`
}

func (m *Player) Reset()                    { *m = Player{} }
func (m *Player) String() string            { return proto.CompactTextString(m) }
func (*Player) ProtoMessage()               {}
func (*Player) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Player) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Player) GetStats() *Stat {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *Player) GetSaves() int32 {
	if m != nil {
		return m.Saves
	}
	return 0
}

type Session struct {
	P1          *Player `protobuf:"bytes,1,opt,name=p1" json:"p1,omitempty"`
	P2          *Player `protobuf:"bytes,2,opt,name=p2" json:"p2,omitempty"`
	HistoricP1  []*Move `protobuf:"bytes,3,rep,name=historicP1" json:"historicP1,omitempty"`
	HistoricP2  []*Move `protobuf:"bytes,4,rep,name=historicP2" json:"historicP2,omitempty"`
	NbCaptureP1 int32   `protobuf:"varint,5,opt,name=nbCaptureP1" json:"nbCaptureP1,omitempty"`
	NbCaptureP2 int32   `protobuf:"varint,6,opt,name=nbCaptureP2" json:"nbCaptureP2,omitempty"`
}

func (m *Session) Reset()                    { *m = Session{} }
func (m *Session) String() string            { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()               {}
func (*Session) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Session) GetP1() *Player {
	if m != nil {
		return m.P1
	}
	return nil
}

func (m *Session) GetP2() *Player {
	if m != nil {
		return m.P2
	}
	return nil
}

func (m *Session) GetHistoricP1() []*Move {
	if m != nil {
		return m.HistoricP1
	}
	return nil
}

func (m *Session) GetHistoricP2() []*Move {
	if m != nil {
		return m.HistoricP2
	}
	return nil
}

func (m *Session) GetNbCaptureP1() int32 {
	if m != nil {
		return m.NbCaptureP1
	}
	return 0
}

func (m *Session) GetNbCaptureP2() int32 {
	if m != nil {
		return m.NbCaptureP2
	}
	return 0
}

type Stat struct {
	NbGame int32 `protobuf:"varint,1,opt,name=nbGame" json:"nbGame,omitempty"`
	NbWin  int32 `protobuf:"varint,2,opt,name=nbWin" json:"nbWin,omitempty"`
}

func (m *Stat) Reset()                    { *m = Stat{} }
func (m *Stat) String() string            { return proto.CompactTextString(m) }
func (*Stat) ProtoMessage()               {}
func (*Stat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Stat) GetNbGame() int32 {
	if m != nil {
		return m.NbGame
	}
	return 0
}

func (m *Stat) GetNbWin() int32 {
	if m != nil {
		return m.NbWin
	}
	return 0
}

type Move struct {
	X      int32 `protobuf:"varint,1,opt,name=x" json:"x,omitempty"`
	Y      int32 `protobuf:"varint,2,opt,name=y" json:"y,omitempty"`
	Reflex int64 `protobuf:"varint,3,opt,name=reflex" json:"reflex,omitempty"`
}

func (m *Move) Reset()                    { *m = Move{} }
func (m *Move) String() string            { return proto.CompactTextString(m) }
func (*Move) ProtoMessage()               {}
func (*Move) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Move) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Move) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Move) GetReflex() int64 {
	if m != nil {
		return m.Reflex
	}
	return 0
}

type Data struct {
	Players  []*Player  `protobuf:"bytes,1,rep,name=players" json:"players,omitempty"`
	Sessions []*Session `protobuf:"bytes,2,rep,name=sessions" json:"sessions,omitempty"`
	Current  *Session   `protobuf:"bytes,3,opt,name=current" json:"current,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Data) GetPlayers() []*Player {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *Data) GetSessions() []*Session {
	if m != nil {
		return m.Sessions
	}
	return nil
}

func (m *Data) GetCurrent() *Session {
	if m != nil {
		return m.Current
	}
	return nil
}

func init() {
	proto.RegisterType((*Player)(nil), "database.Player")
	proto.RegisterType((*Session)(nil), "database.Session")
	proto.RegisterType((*Stat)(nil), "database.Stat")
	proto.RegisterType((*Move)(nil), "database.Move")
	proto.RegisterType((*Data)(nil), "database.Data")
}

func init() { proto.RegisterFile("db.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcd, 0x4a, 0x33, 0x31,
	0x14, 0x86, 0x49, 0xe7, 0xa7, 0xfd, 0x4e, 0x3f, 0x44, 0x0f, 0x22, 0x59, 0x0e, 0x83, 0x8b, 0xa2,
	0x38, 0x30, 0xd1, 0x95, 0x5b, 0x05, 0x57, 0x42, 0x49, 0x17, 0xba, 0xcd, 0xb4, 0x11, 0x07, 0x6a,
	0x66, 0x48, 0xd2, 0xd2, 0xde, 0x84, 0xf7, 0xea, 0x1d, 0x48, 0x92, 0x19, 0x6d, 0xb5, 0xe0, 0x6e,
	0xde, 0x39, 0xcf, 0xf9, 0x7b, 0x4f, 0x60, 0xb4, 0xa8, 0x8a, 0x56, 0x37, 0xb6, 0xc1, 0xd1, 0x42,
	0x58, 0x51, 0x09, 0x23, 0xf3, 0x67, 0x48, 0xa7, 0x4b, 0xb1, 0x95, 0x1a, 0x11, 0x62, 0x25, 0xde,
	0x24, 0x25, 0x19, 0x99, 0xfc, 0xe3, 0xfe, 0x1b, 0xcf, 0x21, 0x31, 0x56, 0x58, 0x43, 0x07, 0x19,
	0x99, 0x8c, 0xd9, 0x51, 0xd1, 0xe7, 0x15, 0x33, 0x2b, 0x2c, 0x0f, 0x41, 0x3c, 0x85, 0xc4, 0x88,
	0xb5, 0x34, 0x34, 0xca, 0xc8, 0x24, 0xe1, 0x41, 0xe4, 0x1f, 0x04, 0x86, 0x33, 0x69, 0x4c, 0xdd,
	0x28, 0xcc, 0x60, 0xd0, 0x96, 0xbe, 0xf2, 0x98, 0x1d, 0x7f, 0x17, 0x09, 0x9d, 0xf9, 0xa0, 0x2d,
	0x3d, 0xc1, 0xba, 0x36, 0x87, 0x08, 0x86, 0x05, 0xc0, 0x6b, 0x6d, 0x6c, 0xa3, 0xeb, 0xf9, 0xb4,
	0xa4, 0x51, 0x16, 0xed, 0x0f, 0xf4, 0xd8, 0xac, 0x25, 0xdf, 0x21, 0xf6, 0x78, 0x46, 0xe3, 0x3f,
	0x78, 0x86, 0x19, 0x8c, 0x55, 0x75, 0x27, 0x5a, 0xbb, 0xd2, 0x72, 0x5a, 0xd2, 0xc4, 0xef, 0xb2,
	0xfb, 0x6b, 0x9f, 0x60, 0x34, 0xfd, 0x49, 0xb0, 0xfc, 0x06, 0x62, 0x67, 0x0c, 0x9e, 0x41, 0xaa,
	0xaa, 0x87, 0xde, 0xcd, 0x84, 0x77, 0xca, 0x39, 0xa5, 0xaa, 0xa7, 0x5a, 0xf9, 0x45, 0x13, 0x1e,
	0x44, 0x7e, 0x0b, 0xb1, 0x9b, 0x06, 0xff, 0x03, 0xd9, 0x74, 0x09, 0x64, 0xe3, 0xd4, 0xb6, 0xe3,
	0xc8, 0xd6, 0x55, 0xd4, 0xf2, 0x65, 0x29, 0x37, 0xde, 0xe4, 0x88, 0x77, 0x2a, 0x7f, 0x27, 0x10,
	0xdf, 0x0b, 0x2b, 0xf0, 0x02, 0x86, 0xad, 0x37, 0xcb, 0x50, 0xe2, 0x77, 0xfd, 0xed, 0x62, 0x0f,
	0xe0, 0x15, 0x8c, 0x4c, 0xb8, 0x8c, 0xbb, 0xac, 0x83, 0x4f, 0x76, 0x2e, 0x1b, 0x22, 0xfc, 0x0b,
	0xc1, 0x4b, 0x18, 0xce, 0x57, 0x5a, 0x4b, 0x65, 0x7d, 0xf3, 0x83, 0x74, 0x4f, 0x54, 0xa9, 0x7f,
	0x61, 0xd7, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf5, 0x5d, 0xfb, 0xfd, 0x6d, 0x02, 0x00, 0x00,
}

// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/tests/integration/data/storage/etcd/store_test.proto
// DO NOT EDIT!

/*
Package etcd is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/tests/integration/data/storage/etcd/store_test.proto

It has these top-level messages:
	TestMessage
*/
package etcd

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

type TestMessage struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *TestMessage) Reset()                    { *m = TestMessage{} }
func (m *TestMessage) String() string            { return proto.CompactTextString(m) }
func (*TestMessage) ProtoMessage()               {}
func (*TestMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TestMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TestMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*TestMessage)(nil), "etcd.TestMessage")
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/tests/integration/data/storage/etcd/store_test.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x1c, 0xcb, 0xc1, 0x0a, 0xc2, 0x30,
	0x10, 0x04, 0x50, 0x5a, 0x8a, 0x60, 0x04, 0x0f, 0x39, 0xf5, 0x28, 0x9e, 0x3c, 0x75, 0x11, 0xbf,
	0x43, 0x10, 0xf1, 0x2e, 0xdb, 0x66, 0x89, 0x01, 0x9b, 0x0d, 0xd9, 0xf1, 0xff, 0xa5, 0xb9, 0xcd,
	0x3c, 0x66, 0xdc, 0x23, 0x26, 0x7c, 0x7e, 0xf3, 0xb4, 0xe8, 0x4a, 0x5c, 0xca, 0x22, 0x5f, 0xa9,
	0x0c, 0xad, 0xc4, 0x6b, 0x21, 0x88, 0xc1, 0x28, 0x65, 0x48, 0xac, 0x8c, 0xa4, 0x99, 0x02, 0x83,
	0xc9, 0xa0, 0x95, 0xa3, 0x90, 0x60, 0x09, 0xad, 0xc8, 0x7b, 0x5b, 0x4e, 0xa5, 0x2a, 0xd4, 0x0f,
	0x1b, 0x9f, 0xaf, 0xee, 0xf0, 0x12, 0xc3, 0x5d, 0xcc, 0x38, 0x8a, 0x3f, 0xba, 0x3e, 0x85, 0xb1,
	0x3b, 0x75, 0x97, 0xfd, 0xb3, 0x4f, 0xc1, 0x7b, 0x37, 0x64, 0x5e, 0x65, 0xec, 0x9b, 0xb4, 0x3c,
	0xef, 0xda, 0xff, 0xf6, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x0f, 0x7d, 0xf0, 0x93, 0x00, 0x00,
	0x00,
}

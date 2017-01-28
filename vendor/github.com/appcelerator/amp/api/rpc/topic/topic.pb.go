// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/api/rpc/topic/topic.proto
// DO NOT EDIT!

/*
Package topic is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/api/rpc/topic/topic.proto

It has these top-level messages:
	TopicEntry
	CreateRequest
	CreateReply
	ListRequest
	ListReply
	DeleteRequest
	DeleteReply
*/
package topic

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TopicEntry struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *TopicEntry) Reset()                    { *m = TopicEntry{} }
func (m *TopicEntry) String() string            { return proto.CompactTextString(m) }
func (*TopicEntry) ProtoMessage()               {}
func (*TopicEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TopicEntry) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TopicEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateRequest struct {
	Topic *TopicEntry `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateRequest) GetTopic() *TopicEntry {
	if m != nil {
		return m.Topic
	}
	return nil
}

type CreateReply struct {
	Topic *TopicEntry `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
}

func (m *CreateReply) Reset()                    { *m = CreateReply{} }
func (m *CreateReply) String() string            { return proto.CompactTextString(m) }
func (*CreateReply) ProtoMessage()               {}
func (*CreateReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CreateReply) GetTopic() *TopicEntry {
	if m != nil {
		return m.Topic
	}
	return nil
}

type ListRequest struct {
}

func (m *ListRequest) Reset()                    { *m = ListRequest{} }
func (m *ListRequest) String() string            { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()               {}
func (*ListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ListReply struct {
	Topics []*TopicEntry `protobuf:"bytes,1,rep,name=topics" json:"topics,omitempty"`
}

func (m *ListReply) Reset()                    { *m = ListReply{} }
func (m *ListReply) String() string            { return proto.CompactTextString(m) }
func (*ListReply) ProtoMessage()               {}
func (*ListReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ListReply) GetTopics() []*TopicEntry {
	if m != nil {
		return m.Topics
	}
	return nil
}

type DeleteRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteReply struct {
	Topic *TopicEntry `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
}

func (m *DeleteReply) Reset()                    { *m = DeleteReply{} }
func (m *DeleteReply) String() string            { return proto.CompactTextString(m) }
func (*DeleteReply) ProtoMessage()               {}
func (*DeleteReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *DeleteReply) GetTopic() *TopicEntry {
	if m != nil {
		return m.Topic
	}
	return nil
}

func init() {
	proto.RegisterType((*TopicEntry)(nil), "topic.TopicEntry")
	proto.RegisterType((*CreateRequest)(nil), "topic.CreateRequest")
	proto.RegisterType((*CreateReply)(nil), "topic.CreateReply")
	proto.RegisterType((*ListRequest)(nil), "topic.ListRequest")
	proto.RegisterType((*ListReply)(nil), "topic.ListReply")
	proto.RegisterType((*DeleteRequest)(nil), "topic.DeleteRequest")
	proto.RegisterType((*DeleteReply)(nil), "topic.DeleteReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Topic service

type TopicClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type topicClient struct {
	cc *grpc.ClientConn
}

func NewTopicClient(cc *grpc.ClientConn) TopicClient {
	return &topicClient{cc}
}

func (c *topicClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error) {
	out := new(CreateReply)
	err := grpc.Invoke(ctx, "/topic.Topic/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error) {
	out := new(ListReply)
	err := grpc.Invoke(ctx, "/topic.Topic/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := grpc.Invoke(ctx, "/topic.Topic/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Topic service

type TopicServer interface {
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	List(context.Context, *ListRequest) (*ListReply, error)
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

func RegisterTopicServer(s *grpc.Server, srv TopicServer) {
	s.RegisterService(&_Topic_serviceDesc, srv)
}

func _Topic_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopicServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/topic.Topic/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopicServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Topic_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopicServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/topic.Topic/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopicServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Topic_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopicServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/topic.Topic/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopicServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Topic_serviceDesc = grpc.ServiceDesc{
	ServiceName: "topic.Topic",
	HandlerType: (*TopicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Topic_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Topic_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Topic_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/appcelerator/amp/api/rpc/topic/topic.proto",
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/api/rpc/topic/topic.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x92, 0xcf, 0x4a, 0x33, 0x31,
	0x14, 0xc5, 0x99, 0xf9, 0xda, 0x81, 0xde, 0xa1, 0xa5, 0xbd, 0x94, 0x8f, 0x52, 0x04, 0x4b, 0x36,
	0x6a, 0x17, 0x8d, 0x56, 0x2c, 0xe2, 0xc6, 0x85, 0x0a, 0x22, 0xae, 0x8a, 0x2f, 0x90, 0xb6, 0xa1,
	0x06, 0xa6, 0x93, 0x98, 0x49, 0x85, 0x22, 0x6e, 0x7c, 0x05, 0x1f, 0xcd, 0xbd, 0x2b, 0x1f, 0x44,
	0x26, 0x99, 0x3f, 0x9d, 0xe2, 0x42, 0x37, 0xc3, 0xe4, 0xce, 0x39, 0xbf, 0x9c, 0x93, 0x09, 0x9c,
	0x2d, 0x85, 0x79, 0x5c, 0xcf, 0x46, 0x73, 0xb9, 0xa2, 0x4c, 0xa9, 0x39, 0x8f, 0xb8, 0x66, 0x46,
	0x6a, 0xca, 0x56, 0x8a, 0x32, 0x25, 0xa8, 0x56, 0x73, 0x6a, 0xa4, 0x12, 0xd9, 0x73, 0xa4, 0xb4,
	0x34, 0x12, 0xeb, 0x76, 0xd1, 0xdf, 0x5b, 0x4a, 0xb9, 0x8c, 0xb8, 0x15, 0xb2, 0x38, 0x96, 0x86,
	0x19, 0x21, 0xe3, 0xc4, 0x89, 0xc8, 0x31, 0xc0, 0x43, 0x2a, 0xbb, 0x89, 0x8d, 0xde, 0x60, 0x0b,
	0x7c, 0xb1, 0xe8, 0x79, 0x03, 0xef, 0xb0, 0x31, 0xf5, 0xc5, 0x02, 0x11, 0x6a, 0x31, 0x5b, 0xf1,
	0x9e, 0x6f, 0x27, 0xf6, 0x9d, 0x9c, 0x43, 0xf3, 0x4a, 0x73, 0x66, 0xf8, 0x94, 0x3f, 0xad, 0x79,
	0x62, 0xf0, 0x00, 0xdc, 0x4e, 0xd6, 0x17, 0x8e, 0x3b, 0x23, 0x17, 0xa2, 0xc4, 0x4e, 0xdd, 0x77,
	0x32, 0x81, 0x30, 0x77, 0xaa, 0x68, 0xf3, 0x7b, 0x5f, 0x13, 0xc2, 0x7b, 0x91, 0x98, 0x6c, 0x3f,
	0x32, 0x81, 0x86, 0x5b, 0xa6, 0x90, 0x23, 0x08, 0xac, 0x28, 0xe9, 0x79, 0x83, 0x7f, 0x3f, 0x53,
	0x32, 0x01, 0xd9, 0x87, 0xe6, 0x35, 0x8f, 0x78, 0x19, 0x7c, 0xa7, 0x6d, 0x9a, 0x2f, 0x17, 0xfc,
	0x25, 0xdf, 0xf8, 0xd3, 0x83, 0xba, 0x9d, 0xe2, 0x2d, 0x04, 0xae, 0x21, 0x76, 0x33, 0x75, 0xe5,
	0xa8, 0xfa, 0xb8, 0x33, 0x55, 0xd1, 0x86, 0x74, 0xdf, 0x3e, 0xbe, 0xde, 0xfd, 0x16, 0x69, 0xd0,
	0xe7, 0x13, 0xf7, 0xff, 0x2e, 0xbc, 0x21, 0x5e, 0x42, 0x2d, 0x2d, 0x89, 0xb9, 0x63, 0xeb, 0x00,
	0xfa, 0xed, 0xca, 0x2c, 0x65, 0x74, 0x2c, 0x23, 0xc4, 0x92, 0x81, 0x77, 0x10, 0xb8, 0x32, 0x45,
	0x94, 0x4a, 0xf9, 0x22, 0xca, 0x56, 0x63, 0xf2, 0xdf, 0x62, 0xda, 0xc3, 0x56, 0x81, 0xa1, 0x2f,
	0x62, 0xf1, 0x3a, 0x0b, 0xec, 0x5d, 0x39, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xea, 0x35, 0x52,
	0xb4, 0x89, 0x02, 0x00, 0x00,
}

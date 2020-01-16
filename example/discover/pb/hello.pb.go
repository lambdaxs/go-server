// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hello.proto

package hello

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type SayHelloReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayHelloReq) Reset()         { *m = SayHelloReq{} }
func (m *SayHelloReq) String() string { return proto.CompactTextString(m) }
func (*SayHelloReq) ProtoMessage()    {}
func (*SayHelloReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{0}
}

func (m *SayHelloReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayHelloReq.Unmarshal(m, b)
}
func (m *SayHelloReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayHelloReq.Marshal(b, m, deterministic)
}
func (m *SayHelloReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayHelloReq.Merge(m, src)
}
func (m *SayHelloReq) XXX_Size() int {
	return xxx_messageInfo_SayHelloReq.Size(m)
}
func (m *SayHelloReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SayHelloReq.DiscardUnknown(m)
}

var xxx_messageInfo_SayHelloReq proto.InternalMessageInfo

func (m *SayHelloReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SayHelloResp struct {
	Reply                string   `protobuf:"bytes,1,opt,name=reply,proto3" json:"reply,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayHelloResp) Reset()         { *m = SayHelloResp{} }
func (m *SayHelloResp) String() string { return proto.CompactTextString(m) }
func (*SayHelloResp) ProtoMessage()    {}
func (*SayHelloResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{1}
}

func (m *SayHelloResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayHelloResp.Unmarshal(m, b)
}
func (m *SayHelloResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayHelloResp.Marshal(b, m, deterministic)
}
func (m *SayHelloResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayHelloResp.Merge(m, src)
}
func (m *SayHelloResp) XXX_Size() int {
	return xxx_messageInfo_SayHelloResp.Size(m)
}
func (m *SayHelloResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SayHelloResp.DiscardUnknown(m)
}

var xxx_messageInfo_SayHelloResp proto.InternalMessageInfo

func (m *SayHelloResp) GetReply() string {
	if m != nil {
		return m.Reply
	}
	return ""
}

func init() {
	proto.RegisterType((*SayHelloReq)(nil), "SayHelloReq")
	proto.RegisterType((*SayHelloResp)(nil), "SayHelloResp")
}

func init() { proto.RegisterFile("hello.proto", fileDescriptor_61ef911816e0a8ce) }

var fileDescriptor_61ef911816e0a8ce = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe4, 0xe2, 0x0e, 0x4e, 0xac, 0xf4, 0x00,
	0x89, 0x04, 0xa5, 0x16, 0x0a, 0x09, 0x71, 0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30,
	0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0x2a, 0x5c, 0x3c, 0x08, 0x25, 0xc5, 0x05, 0x42, 0x22, 0x5c,
	0xac, 0x45, 0xa9, 0x05, 0x39, 0x95, 0x50, 0x45, 0x10, 0x8e, 0x91, 0x05, 0x17, 0x37, 0x58, 0x49,
	0x70, 0x6a, 0x51, 0x59, 0x6a, 0x91, 0x90, 0x26, 0x17, 0x07, 0x4c, 0x93, 0x10, 0x8f, 0x1e, 0x92,
	0x15, 0x52, 0xbc, 0x7a, 0xc8, 0xa6, 0x29, 0x31, 0x24, 0xb1, 0x81, 0x5d, 0x62, 0x0c, 0x08, 0x00,
	0x00, 0xff, 0xff, 0xb9, 0x42, 0x0a, 0xfb, 0x98, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HelloServerClient is the client API for HelloServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HelloServerClient interface {
	SayHello(ctx context.Context, in *SayHelloReq, opts ...grpc.CallOption) (*SayHelloResp, error)
}

type helloServerClient struct {
	cc *grpc.ClientConn
}

func NewHelloServerClient(cc *grpc.ClientConn) HelloServerClient {
	return &helloServerClient{cc}
}

func (c *helloServerClient) SayHello(ctx context.Context, in *SayHelloReq, opts ...grpc.CallOption) (*SayHelloResp, error) {
	out := new(SayHelloResp)
	err := c.cc.Invoke(ctx, "/HelloServer/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloServerServer is the server API for HelloServer service.
type HelloServerServer interface {
	SayHello(context.Context, *SayHelloReq) (*SayHelloResp, error)
}

// UnimplementedHelloServerServer can be embedded to have forward compatible implementations.
type UnimplementedHelloServerServer struct {
}

func (*UnimplementedHelloServerServer) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

func RegisterHelloServerServer(s *grpc.Server, srv HelloServerServer) {
	s.RegisterService(&_HelloServer_serviceDesc, srv)
}

func _HelloServer_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayHelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServerServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HelloServer/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServerServer).SayHello(ctx, req.(*SayHelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _HelloServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "HelloServer",
	HandlerType: (*HelloServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _HelloServer_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hello.proto",
}

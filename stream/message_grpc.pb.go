// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: stream/message.proto

package stream

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Cast_CreateStream_FullMethodName = "/protobuf.Cast/CreateStream"
	Cast_CastMessage_FullMethodName  = "/protobuf.Cast/CastMessage"
)

// CastClient is the client API for Cast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CastClient interface {
	CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Cast_CreateStreamClient, error)
	CastMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Close, error)
}

type castClient struct {
	cc grpc.ClientConnInterface
}

func NewCastClient(cc grpc.ClientConnInterface) CastClient {
	return &castClient{cc}
}

func (c *castClient) CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Cast_CreateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Cast_ServiceDesc.Streams[0], Cast_CreateStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &castCreateStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Cast_CreateStreamClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type castCreateStreamClient struct {
	grpc.ClientStream
}

func (x *castCreateStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *castClient) CastMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Close, error) {
	out := new(Close)
	err := c.cc.Invoke(ctx, Cast_CastMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CastServer is the server API for Cast service.
// All implementations must embed UnimplementedCastServer
// for forward compatibility
type CastServer interface {
	CreateStream(*Connect, Cast_CreateStreamServer) error
	CastMessage(context.Context, *Message) (*Close, error)
	mustEmbedUnimplementedCastServer()
}

// UnimplementedCastServer must be embedded to have forward compatible implementations.
type UnimplementedCastServer struct {
}

func (UnimplementedCastServer) CreateStream(*Connect, Cast_CreateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}
func (UnimplementedCastServer) CastMessage(context.Context, *Message) (*Close, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CastMessage not implemented")
}
func (UnimplementedCastServer) mustEmbedUnimplementedCastServer() {}

// UnsafeCastServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CastServer will
// result in compilation errors.
type UnsafeCastServer interface {
	mustEmbedUnimplementedCastServer()
}

func RegisterCastServer(s grpc.ServiceRegistrar, srv CastServer) {
	s.RegisterService(&Cast_ServiceDesc, srv)
}

func _Cast_CreateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Connect)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CastServer).CreateStream(m, &castCreateStreamServer{stream})
}

type Cast_CreateStreamServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type castCreateStreamServer struct {
	grpc.ServerStream
}

func (x *castCreateStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _Cast_CastMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CastServer).CastMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cast_CastMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CastServer).CastMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// Cast_ServiceDesc is the grpc.ServiceDesc for Cast service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cast_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Cast",
	HandlerType: (*CastServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CastMessage",
			Handler:    _Cast_CastMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateStream",
			Handler:       _Cast_CreateStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "stream/message.proto",
}
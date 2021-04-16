// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package syncchainpb

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

// SyncchainClient is the client API for Syncchain service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SyncchainClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type syncchainClient struct {
	cc grpc.ClientConnInterface
}

func NewSyncchainClient(cc grpc.ClientConnInterface) SyncchainClient {
	return &syncchainClient{cc}
}

func (c *syncchainClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/syncchainpb.Syncchain/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SyncchainServer is the server API for Syncchain service.
// All implementations must embed UnimplementedSyncchainServer
// for forward compatibility
type SyncchainServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedSyncchainServer()
}

// UnimplementedSyncchainServer must be embedded to have forward compatible implementations.
type UnimplementedSyncchainServer struct {
}

func (UnimplementedSyncchainServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSyncchainServer) mustEmbedUnimplementedSyncchainServer() {}

// UnsafeSyncchainServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SyncchainServer will
// result in compilation errors.
type UnsafeSyncchainServer interface {
	mustEmbedUnimplementedSyncchainServer()
}

func RegisterSyncchainServer(s grpc.ServiceRegistrar, srv SyncchainServer) {
	s.RegisterService(&Syncchain_ServiceDesc, srv)
}

func _Syncchain_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncchainServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syncchainpb.Syncchain/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncchainServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Syncchain_ServiceDesc is the grpc.ServiceDesc for Syncchain service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Syncchain_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syncchainpb.Syncchain",
	HandlerType: (*SyncchainServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Syncchain_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syncchain.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/indexer.proto

package pb

import (
	context "context"
	pb "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IndexerServiceClient is the client API for IndexerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IndexerServiceClient interface {
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (IndexerService_SubscribeClient, error)
	Unsubscribe(ctx context.Context, in *pb.UnsubscribeRequest, opts ...grpc.CallOption) (*pb.UnsubscribeResponse, error)
	JSONSchemaForClass(ctx context.Context, in *Bytes, opts ...grpc.CallOption) (*Bytes, error)
	JSONSchemaForContract(ctx context.Context, in *Bytes, opts ...grpc.CallOption) (*Bytes, error)
}

type indexerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIndexerServiceClient(cc grpc.ClientConnInterface) IndexerServiceClient {
	return &indexerServiceClient{cc}
}

func (c *indexerServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (IndexerService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &IndexerService_ServiceDesc.Streams[0], "/proto.IndexerService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &indexerServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type IndexerService_SubscribeClient interface {
	Recv() (*Subscription, error)
	grpc.ClientStream
}

type indexerServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *indexerServiceSubscribeClient) Recv() (*Subscription, error) {
	m := new(Subscription)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *indexerServiceClient) Unsubscribe(ctx context.Context, in *pb.UnsubscribeRequest, opts ...grpc.CallOption) (*pb.UnsubscribeResponse, error) {
	out := new(pb.UnsubscribeResponse)
	err := c.cc.Invoke(ctx, "/proto.IndexerService/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexerServiceClient) JSONSchemaForClass(ctx context.Context, in *Bytes, opts ...grpc.CallOption) (*Bytes, error) {
	out := new(Bytes)
	err := c.cc.Invoke(ctx, "/proto.IndexerService/JSONSchemaForClass", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexerServiceClient) JSONSchemaForContract(ctx context.Context, in *Bytes, opts ...grpc.CallOption) (*Bytes, error) {
	out := new(Bytes)
	err := c.cc.Invoke(ctx, "/proto.IndexerService/JSONSchemaForContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexerServiceServer is the server API for IndexerService service.
// All implementations must embed UnimplementedIndexerServiceServer
// for forward compatibility
type IndexerServiceServer interface {
	Subscribe(*SubscribeRequest, IndexerService_SubscribeServer) error
	Unsubscribe(context.Context, *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error)
	JSONSchemaForClass(context.Context, *Bytes) (*Bytes, error)
	JSONSchemaForContract(context.Context, *Bytes) (*Bytes, error)
	mustEmbedUnimplementedIndexerServiceServer()
}

// UnimplementedIndexerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIndexerServiceServer struct {
}

func (UnimplementedIndexerServiceServer) Subscribe(*SubscribeRequest, IndexerService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedIndexerServiceServer) Unsubscribe(context.Context, *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (UnimplementedIndexerServiceServer) JSONSchemaForClass(context.Context, *Bytes) (*Bytes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JSONSchemaForClass not implemented")
}
func (UnimplementedIndexerServiceServer) JSONSchemaForContract(context.Context, *Bytes) (*Bytes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JSONSchemaForContract not implemented")
}
func (UnimplementedIndexerServiceServer) mustEmbedUnimplementedIndexerServiceServer() {}

// UnsafeIndexerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IndexerServiceServer will
// result in compilation errors.
type UnsafeIndexerServiceServer interface {
	mustEmbedUnimplementedIndexerServiceServer()
}

func RegisterIndexerServiceServer(s grpc.ServiceRegistrar, srv IndexerServiceServer) {
	s.RegisterService(&IndexerService_ServiceDesc, srv)
}

func _IndexerService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IndexerServiceServer).Subscribe(m, &indexerServiceSubscribeServer{stream})
}

type IndexerService_SubscribeServer interface {
	Send(*Subscription) error
	grpc.ServerStream
}

type indexerServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *indexerServiceSubscribeServer) Send(m *Subscription) error {
	return x.ServerStream.SendMsg(m)
}

func _IndexerService_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(pb.UnsubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexerServiceServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.IndexerService/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexerServiceServer).Unsubscribe(ctx, req.(*pb.UnsubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexerService_JSONSchemaForClass_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Bytes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexerServiceServer).JSONSchemaForClass(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.IndexerService/JSONSchemaForClass",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexerServiceServer).JSONSchemaForClass(ctx, req.(*Bytes))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexerService_JSONSchemaForContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Bytes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexerServiceServer).JSONSchemaForContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.IndexerService/JSONSchemaForContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexerServiceServer).JSONSchemaForContract(ctx, req.(*Bytes))
	}
	return interceptor(ctx, in, info, handler)
}

// IndexerService_ServiceDesc is the grpc.ServiceDesc for IndexerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IndexerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.IndexerService",
	HandlerType: (*IndexerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unsubscribe",
			Handler:    _IndexerService_Unsubscribe_Handler,
		},
		{
			MethodName: "JSONSchemaForClass",
			Handler:    _IndexerService_JSONSchemaForClass_Handler,
		},
		{
			MethodName: "JSONSchemaForContract",
			Handler:    _IndexerService_JSONSchemaForContract_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _IndexerService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/indexer.proto",
}

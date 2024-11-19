// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: terraformer.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TerraformerService_BuildInfrastructure_FullMethodName   = "/claudie.TerraformerService/BuildInfrastructure"
	TerraformerService_DestroyInfrastructure_FullMethodName = "/claudie.TerraformerService/DestroyInfrastructure"
)

// TerraformerServiceClient is the client API for TerraformerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TerraformerServiceClient interface {
	// BuildInfrastructure builds the infrastructure based on the provided desired state (includes addition/deletion of *stuff*).
	BuildInfrastructure(ctx context.Context, in *BuildInfrastructureRequest, opts ...grpc.CallOption) (*BuildInfrastructureResponse, error)
	// DestroyInfrastructure destroys the infrastructure completely.
	DestroyInfrastructure(ctx context.Context, in *DestroyInfrastructureRequest, opts ...grpc.CallOption) (*DestroyInfrastructureResponse, error)
}

type terraformerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTerraformerServiceClient(cc grpc.ClientConnInterface) TerraformerServiceClient {
	return &terraformerServiceClient{cc}
}

func (c *terraformerServiceClient) BuildInfrastructure(ctx context.Context, in *BuildInfrastructureRequest, opts ...grpc.CallOption) (*BuildInfrastructureResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BuildInfrastructureResponse)
	err := c.cc.Invoke(ctx, TerraformerService_BuildInfrastructure_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terraformerServiceClient) DestroyInfrastructure(ctx context.Context, in *DestroyInfrastructureRequest, opts ...grpc.CallOption) (*DestroyInfrastructureResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DestroyInfrastructureResponse)
	err := c.cc.Invoke(ctx, TerraformerService_DestroyInfrastructure_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TerraformerServiceServer is the server API for TerraformerService service.
// All implementations must embed UnimplementedTerraformerServiceServer
// for forward compatibility.
type TerraformerServiceServer interface {
	// BuildInfrastructure builds the infrastructure based on the provided desired state (includes addition/deletion of *stuff*).
	BuildInfrastructure(context.Context, *BuildInfrastructureRequest) (*BuildInfrastructureResponse, error)
	// DestroyInfrastructure destroys the infrastructure completely.
	DestroyInfrastructure(context.Context, *DestroyInfrastructureRequest) (*DestroyInfrastructureResponse, error)
	mustEmbedUnimplementedTerraformerServiceServer()
}

// UnimplementedTerraformerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTerraformerServiceServer struct{}

func (UnimplementedTerraformerServiceServer) BuildInfrastructure(context.Context, *BuildInfrastructureRequest) (*BuildInfrastructureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildInfrastructure not implemented")
}
func (UnimplementedTerraformerServiceServer) DestroyInfrastructure(context.Context, *DestroyInfrastructureRequest) (*DestroyInfrastructureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroyInfrastructure not implemented")
}
func (UnimplementedTerraformerServiceServer) mustEmbedUnimplementedTerraformerServiceServer() {}
func (UnimplementedTerraformerServiceServer) testEmbeddedByValue()                            {}

// UnsafeTerraformerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TerraformerServiceServer will
// result in compilation errors.
type UnsafeTerraformerServiceServer interface {
	mustEmbedUnimplementedTerraformerServiceServer()
}

func RegisterTerraformerServiceServer(s grpc.ServiceRegistrar, srv TerraformerServiceServer) {
	// If the following call pancis, it indicates UnimplementedTerraformerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TerraformerService_ServiceDesc, srv)
}

func _TerraformerService_BuildInfrastructure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuildInfrastructureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerraformerServiceServer).BuildInfrastructure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerraformerService_BuildInfrastructure_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerraformerServiceServer).BuildInfrastructure(ctx, req.(*BuildInfrastructureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerraformerService_DestroyInfrastructure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DestroyInfrastructureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerraformerServiceServer).DestroyInfrastructure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerraformerService_DestroyInfrastructure_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerraformerServiceServer).DestroyInfrastructure(ctx, req.(*DestroyInfrastructureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TerraformerService_ServiceDesc is the grpc.ServiceDesc for TerraformerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TerraformerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "claudie.TerraformerService",
	HandlerType: (*TerraformerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BuildInfrastructure",
			Handler:    _TerraformerService_BuildInfrastructure_Handler,
		},
		{
			MethodName: "DestroyInfrastructure",
			Handler:    _TerraformerService_DestroyInfrastructure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "terraformer.proto",
}

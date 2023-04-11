// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: proto/contextBox.proto

package pb

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

// ContextBoxServiceClient is the client API for ContextBoxService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContextBoxServiceClient interface {
	// SaveConfigFrontEnd saves the config parsed by Frontend.
	SaveConfigFrontEnd(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error)
	// SaveConfigScheduler saves the config parsed by Scheduler.
	SaveConfigScheduler(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error)
	// SaveConfigBuilder saves the config parsed by Builder.
	SaveConfigBuilder(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error)
	// GetConfigFromDB gets a single config from the database.
	GetConfigFromDB(ctx context.Context, in *GetConfigFromDBRequest, opts ...grpc.CallOption) (*GetConfigFromDBResponse, error)
	// GetConfigScheduler gets a config from Scheduler's queue of pending configs.
	GetConfigScheduler(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
	// GetConfigBuilder gets a config from Builder's queue of pending configs.
	GetConfigBuilder(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
	// GetAllConfigs gets all configs from the database.
	GetAllConfigs(ctx context.Context, in *GetAllConfigsRequest, opts ...grpc.CallOption) (*GetAllConfigsResponse, error)
	// DeleteConfig sets the manifest to null, effectively forcing the deletion of the infrastructure
	// defined by the manifest on the very next config (diff) check.
	DeleteConfig(ctx context.Context, in *DeleteConfigRequest, opts ...grpc.CallOption) (*DeleteConfigResponse, error)
	// DeleteConfigFromDB deletes the config from the database.
	DeleteConfigFromDB(ctx context.Context, in *DeleteConfigRequest, opts ...grpc.CallOption) (*DeleteConfigResponse, error)
	// UpdateNodepool updates specific nodepool from the config. Used mainly for autoscaling.
	UpdateNodepool(ctx context.Context, in *UpdateNodepoolRequest, opts ...grpc.CallOption) (*UpdateNodepoolResponse, error)
}

type contextBoxServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewContextBoxServiceClient(cc grpc.ClientConnInterface) ContextBoxServiceClient {
	return &contextBoxServiceClient{cc}
}

func (c *contextBoxServiceClient) SaveConfigFrontEnd(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error) {
	out := new(SaveConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/SaveConfigFrontEnd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) SaveConfigScheduler(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error) {
	out := new(SaveConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/SaveConfigScheduler", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) SaveConfigBuilder(ctx context.Context, in *SaveConfigRequest, opts ...grpc.CallOption) (*SaveConfigResponse, error) {
	out := new(SaveConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/SaveConfigBuilder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) GetConfigFromDB(ctx context.Context, in *GetConfigFromDBRequest, opts ...grpc.CallOption) (*GetConfigFromDBResponse, error) {
	out := new(GetConfigFromDBResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/GetConfigFromDB", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) GetConfigScheduler(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/GetConfigScheduler", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) GetConfigBuilder(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/GetConfigBuilder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) GetAllConfigs(ctx context.Context, in *GetAllConfigsRequest, opts ...grpc.CallOption) (*GetAllConfigsResponse, error) {
	out := new(GetAllConfigsResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/GetAllConfigs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) DeleteConfig(ctx context.Context, in *DeleteConfigRequest, opts ...grpc.CallOption) (*DeleteConfigResponse, error) {
	out := new(DeleteConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/DeleteConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) DeleteConfigFromDB(ctx context.Context, in *DeleteConfigRequest, opts ...grpc.CallOption) (*DeleteConfigResponse, error) {
	out := new(DeleteConfigResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/DeleteConfigFromDB", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextBoxServiceClient) UpdateNodepool(ctx context.Context, in *UpdateNodepoolRequest, opts ...grpc.CallOption) (*UpdateNodepoolResponse, error) {
	out := new(UpdateNodepoolResponse)
	err := c.cc.Invoke(ctx, "/claudie.ContextBoxService/UpdateNodepool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContextBoxServiceServer is the server API for ContextBoxService service.
// All implementations must embed UnimplementedContextBoxServiceServer
// for forward compatibility
type ContextBoxServiceServer interface {
	// SaveConfigFrontEnd saves the config parsed by Frontend.
	SaveConfigFrontEnd(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error)
	// SaveConfigScheduler saves the config parsed by Scheduler.
	SaveConfigScheduler(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error)
	// SaveConfigBuilder saves the config parsed by Builder.
	SaveConfigBuilder(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error)
	// GetConfigFromDB gets a single config from the database.
	GetConfigFromDB(context.Context, *GetConfigFromDBRequest) (*GetConfigFromDBResponse, error)
	// GetConfigScheduler gets a config from Scheduler's queue of pending configs.
	GetConfigScheduler(context.Context, *GetConfigRequest) (*GetConfigResponse, error)
	// GetConfigBuilder gets a config from Builder's queue of pending configs.
	GetConfigBuilder(context.Context, *GetConfigRequest) (*GetConfigResponse, error)
	// GetAllConfigs gets all configs from the database.
	GetAllConfigs(context.Context, *GetAllConfigsRequest) (*GetAllConfigsResponse, error)
	// DeleteConfig sets the manifest to null, effectively forcing the deletion of the infrastructure
	// defined by the manifest on the very next config (diff) check.
	DeleteConfig(context.Context, *DeleteConfigRequest) (*DeleteConfigResponse, error)
	// DeleteConfigFromDB deletes the config from the database.
	DeleteConfigFromDB(context.Context, *DeleteConfigRequest) (*DeleteConfigResponse, error)
	// UpdateNodepool updates specific nodepool from the config. Used mainly for autoscaling.
	UpdateNodepool(context.Context, *UpdateNodepoolRequest) (*UpdateNodepoolResponse, error)
	mustEmbedUnimplementedContextBoxServiceServer()
}

// UnimplementedContextBoxServiceServer must be embedded to have forward compatible implementations.
type UnimplementedContextBoxServiceServer struct {
}

func (UnimplementedContextBoxServiceServer) SaveConfigFrontEnd(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveConfigFrontEnd not implemented")
}
func (UnimplementedContextBoxServiceServer) SaveConfigScheduler(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveConfigScheduler not implemented")
}
func (UnimplementedContextBoxServiceServer) SaveConfigBuilder(context.Context, *SaveConfigRequest) (*SaveConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveConfigBuilder not implemented")
}
func (UnimplementedContextBoxServiceServer) GetConfigFromDB(context.Context, *GetConfigFromDBRequest) (*GetConfigFromDBResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigFromDB not implemented")
}
func (UnimplementedContextBoxServiceServer) GetConfigScheduler(context.Context, *GetConfigRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigScheduler not implemented")
}
func (UnimplementedContextBoxServiceServer) GetConfigBuilder(context.Context, *GetConfigRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigBuilder not implemented")
}
func (UnimplementedContextBoxServiceServer) GetAllConfigs(context.Context, *GetAllConfigsRequest) (*GetAllConfigsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllConfigs not implemented")
}
func (UnimplementedContextBoxServiceServer) DeleteConfig(context.Context, *DeleteConfigRequest) (*DeleteConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConfig not implemented")
}
func (UnimplementedContextBoxServiceServer) DeleteConfigFromDB(context.Context, *DeleteConfigRequest) (*DeleteConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConfigFromDB not implemented")
}
func (UnimplementedContextBoxServiceServer) UpdateNodepool(context.Context, *UpdateNodepoolRequest) (*UpdateNodepoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNodepool not implemented")
}
func (UnimplementedContextBoxServiceServer) mustEmbedUnimplementedContextBoxServiceServer() {}

// UnsafeContextBoxServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContextBoxServiceServer will
// result in compilation errors.
type UnsafeContextBoxServiceServer interface {
	mustEmbedUnimplementedContextBoxServiceServer()
}

func RegisterContextBoxServiceServer(s grpc.ServiceRegistrar, srv ContextBoxServiceServer) {
	s.RegisterService(&ContextBoxService_ServiceDesc, srv)
}

func _ContextBoxService_SaveConfigFrontEnd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).SaveConfigFrontEnd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/SaveConfigFrontEnd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).SaveConfigFrontEnd(ctx, req.(*SaveConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_SaveConfigScheduler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).SaveConfigScheduler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/SaveConfigScheduler",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).SaveConfigScheduler(ctx, req.(*SaveConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_SaveConfigBuilder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).SaveConfigBuilder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/SaveConfigBuilder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).SaveConfigBuilder(ctx, req.(*SaveConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_GetConfigFromDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigFromDBRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).GetConfigFromDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/GetConfigFromDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).GetConfigFromDB(ctx, req.(*GetConfigFromDBRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_GetConfigScheduler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).GetConfigScheduler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/GetConfigScheduler",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).GetConfigScheduler(ctx, req.(*GetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_GetConfigBuilder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).GetConfigBuilder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/GetConfigBuilder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).GetConfigBuilder(ctx, req.(*GetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_GetAllConfigs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllConfigsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).GetAllConfigs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/GetAllConfigs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).GetAllConfigs(ctx, req.(*GetAllConfigsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_DeleteConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).DeleteConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/DeleteConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).DeleteConfig(ctx, req.(*DeleteConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_DeleteConfigFromDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).DeleteConfigFromDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/DeleteConfigFromDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).DeleteConfigFromDB(ctx, req.(*DeleteConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextBoxService_UpdateNodepool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNodepoolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextBoxServiceServer).UpdateNodepool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/claudie.ContextBoxService/UpdateNodepool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextBoxServiceServer).UpdateNodepool(ctx, req.(*UpdateNodepoolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ContextBoxService_ServiceDesc is the grpc.ServiceDesc for ContextBoxService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ContextBoxService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "claudie.ContextBoxService",
	HandlerType: (*ContextBoxServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveConfigFrontEnd",
			Handler:    _ContextBoxService_SaveConfigFrontEnd_Handler,
		},
		{
			MethodName: "SaveConfigScheduler",
			Handler:    _ContextBoxService_SaveConfigScheduler_Handler,
		},
		{
			MethodName: "SaveConfigBuilder",
			Handler:    _ContextBoxService_SaveConfigBuilder_Handler,
		},
		{
			MethodName: "GetConfigFromDB",
			Handler:    _ContextBoxService_GetConfigFromDB_Handler,
		},
		{
			MethodName: "GetConfigScheduler",
			Handler:    _ContextBoxService_GetConfigScheduler_Handler,
		},
		{
			MethodName: "GetConfigBuilder",
			Handler:    _ContextBoxService_GetConfigBuilder_Handler,
		},
		{
			MethodName: "GetAllConfigs",
			Handler:    _ContextBoxService_GetAllConfigs_Handler,
		},
		{
			MethodName: "DeleteConfig",
			Handler:    _ContextBoxService_DeleteConfig_Handler,
		},
		{
			MethodName: "DeleteConfigFromDB",
			Handler:    _ContextBoxService_DeleteConfigFromDB_Handler,
		},
		{
			MethodName: "UpdateNodepool",
			Handler:    _ContextBoxService_UpdateNodepool_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/contextBox.proto",
}

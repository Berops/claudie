// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: ansibler.proto

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
	AnsiblerService_InstallNodeRequirements_FullMethodName       = "/claudie.AnsiblerService/InstallNodeRequirements"
	AnsiblerService_InstallVPN_FullMethodName                    = "/claudie.AnsiblerService/InstallVPN"
	AnsiblerService_SetUpLoadbalancers_FullMethodName            = "/claudie.AnsiblerService/SetUpLoadbalancers"
	AnsiblerService_DetermineApiEndpointChange_FullMethodName    = "/claudie.AnsiblerService/DetermineApiEndpointChange"
	AnsiblerService_UpdateAPIEndpoint_FullMethodName             = "/claudie.AnsiblerService/UpdateAPIEndpoint"
	AnsiblerService_UpdateProxyEnvsOnNodes_FullMethodName        = "/claudie.AnsiblerService/UpdateProxyEnvsOnNodes"
	AnsiblerService_UpdateNoProxyEnvsInKubernetes_FullMethodName = "/claudie.AnsiblerService/UpdateNoProxyEnvsInKubernetes"
	AnsiblerService_RemoveClaudieUtilities_FullMethodName        = "/claudie.AnsiblerService/RemoveClaudieUtilities"
)

// AnsiblerServiceClient is the client API for AnsiblerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnsiblerServiceClient interface {
	// InstallNodeRequirements installs any requirements there are on all of the nodes.
	InstallNodeRequirements(ctx context.Context, in *InstallRequest, opts ...grpc.CallOption) (*InstallResponse, error)
	// InstallVPN sets up a VPN between the nodes in the k8s cluster and LB clusters.
	InstallVPN(ctx context.Context, in *InstallRequest, opts ...grpc.CallOption) (*InstallResponse, error)
	// SetUpLoadbalancers sets up the load balancers together with the DNS and verifies their configuration.
	SetUpLoadbalancers(ctx context.Context, in *SetUpLBRequest, opts ...grpc.CallOption) (*SetUpLBResponse, error)
	// DetermineApiEndpointChange determines if due to the changes of the loadbalancer infrastructure the api endpoint
	// needs to be moved.
	DetermineApiEndpointChange(ctx context.Context, in *DetermineApiEndpointChangeRequest, opts ...grpc.CallOption) (*DetermineApiEndpointChangeResponse, error)
	// UpdateAPIEndpoint handles changes of API endpoint between control nodes.
	// It will update the current stage based on the information from the desired state.
	UpdateAPIEndpoint(ctx context.Context, in *UpdateAPIEndpointRequest, opts ...grpc.CallOption) (*UpdateAPIEndpointResponse, error)
	// UpdateProxyEnvsOnNodes handles changes of HTTP_PROXY, HTTPS_PROXY, NO_PROXY, http_proxy, https_proxy and no_proxy envs in /etc/environment
	UpdateProxyEnvsOnNodes(ctx context.Context, in *UpdateProxyEnvsOnNodesRequest, opts ...grpc.CallOption) (*UpdateProxyEnvsOnNodesResponse, error)
	// UpdateNoProxyEnvsInKubernetes handles changes of NO_PROXY and no_proxy envs in kube-proxy and static pods
	UpdateNoProxyEnvsInKubernetes(ctx context.Context, in *UpdateNoProxyEnvsInKubernetesRequest, opts ...grpc.CallOption) (*UpdateNoProxyEnvsInKubernetesResponse, error)
	// Removes utilities installed by claudie via ansible playbooks.
	RemoveClaudieUtilities(ctx context.Context, in *RemoveClaudieUtilitiesRequest, opts ...grpc.CallOption) (*RemoveClaudieUtilitiesResponse, error)
}

type ansiblerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAnsiblerServiceClient(cc grpc.ClientConnInterface) AnsiblerServiceClient {
	return &ansiblerServiceClient{cc}
}

func (c *ansiblerServiceClient) InstallNodeRequirements(ctx context.Context, in *InstallRequest, opts ...grpc.CallOption) (*InstallResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InstallResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_InstallNodeRequirements_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) InstallVPN(ctx context.Context, in *InstallRequest, opts ...grpc.CallOption) (*InstallResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InstallResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_InstallVPN_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) SetUpLoadbalancers(ctx context.Context, in *SetUpLBRequest, opts ...grpc.CallOption) (*SetUpLBResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetUpLBResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_SetUpLoadbalancers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) DetermineApiEndpointChange(ctx context.Context, in *DetermineApiEndpointChangeRequest, opts ...grpc.CallOption) (*DetermineApiEndpointChangeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetermineApiEndpointChangeResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_DetermineApiEndpointChange_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) UpdateAPIEndpoint(ctx context.Context, in *UpdateAPIEndpointRequest, opts ...grpc.CallOption) (*UpdateAPIEndpointResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAPIEndpointResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_UpdateAPIEndpoint_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) UpdateProxyEnvsOnNodes(ctx context.Context, in *UpdateProxyEnvsOnNodesRequest, opts ...grpc.CallOption) (*UpdateProxyEnvsOnNodesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateProxyEnvsOnNodesResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_UpdateProxyEnvsOnNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) UpdateNoProxyEnvsInKubernetes(ctx context.Context, in *UpdateNoProxyEnvsInKubernetesRequest, opts ...grpc.CallOption) (*UpdateNoProxyEnvsInKubernetesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateNoProxyEnvsInKubernetesResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_UpdateNoProxyEnvsInKubernetes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ansiblerServiceClient) RemoveClaudieUtilities(ctx context.Context, in *RemoveClaudieUtilitiesRequest, opts ...grpc.CallOption) (*RemoveClaudieUtilitiesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveClaudieUtilitiesResponse)
	err := c.cc.Invoke(ctx, AnsiblerService_RemoveClaudieUtilities_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnsiblerServiceServer is the server API for AnsiblerService service.
// All implementations must embed UnimplementedAnsiblerServiceServer
// for forward compatibility.
type AnsiblerServiceServer interface {
	// InstallNodeRequirements installs any requirements there are on all of the nodes.
	InstallNodeRequirements(context.Context, *InstallRequest) (*InstallResponse, error)
	// InstallVPN sets up a VPN between the nodes in the k8s cluster and LB clusters.
	InstallVPN(context.Context, *InstallRequest) (*InstallResponse, error)
	// SetUpLoadbalancers sets up the load balancers together with the DNS and verifies their configuration.
	SetUpLoadbalancers(context.Context, *SetUpLBRequest) (*SetUpLBResponse, error)
	// DetermineApiEndpointChange determines if due to the changes of the loadbalancer infrastructure the api endpoint
	// needs to be moved.
	DetermineApiEndpointChange(context.Context, *DetermineApiEndpointChangeRequest) (*DetermineApiEndpointChangeResponse, error)
	// UpdateAPIEndpoint handles changes of API endpoint between control nodes.
	// It will update the current stage based on the information from the desired state.
	UpdateAPIEndpoint(context.Context, *UpdateAPIEndpointRequest) (*UpdateAPIEndpointResponse, error)
	// UpdateProxyEnvsOnNodes handles changes of HTTP_PROXY, HTTPS_PROXY, NO_PROXY, http_proxy, https_proxy and no_proxy envs in /etc/environment
	UpdateProxyEnvsOnNodes(context.Context, *UpdateProxyEnvsOnNodesRequest) (*UpdateProxyEnvsOnNodesResponse, error)
	// UpdateNoProxyEnvsInKubernetes handles changes of NO_PROXY and no_proxy envs in kube-proxy and static pods
	UpdateNoProxyEnvsInKubernetes(context.Context, *UpdateNoProxyEnvsInKubernetesRequest) (*UpdateNoProxyEnvsInKubernetesResponse, error)
	// Removes utilities installed by claudie via ansible playbooks.
	RemoveClaudieUtilities(context.Context, *RemoveClaudieUtilitiesRequest) (*RemoveClaudieUtilitiesResponse, error)
	mustEmbedUnimplementedAnsiblerServiceServer()
}

// UnimplementedAnsiblerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAnsiblerServiceServer struct{}

func (UnimplementedAnsiblerServiceServer) InstallNodeRequirements(context.Context, *InstallRequest) (*InstallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallNodeRequirements not implemented")
}
func (UnimplementedAnsiblerServiceServer) InstallVPN(context.Context, *InstallRequest) (*InstallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallVPN not implemented")
}
func (UnimplementedAnsiblerServiceServer) SetUpLoadbalancers(context.Context, *SetUpLBRequest) (*SetUpLBResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUpLoadbalancers not implemented")
}
func (UnimplementedAnsiblerServiceServer) DetermineApiEndpointChange(context.Context, *DetermineApiEndpointChangeRequest) (*DetermineApiEndpointChangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetermineApiEndpointChange not implemented")
}
func (UnimplementedAnsiblerServiceServer) UpdateAPIEndpoint(context.Context, *UpdateAPIEndpointRequest) (*UpdateAPIEndpointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAPIEndpoint not implemented")
}
func (UnimplementedAnsiblerServiceServer) UpdateProxyEnvsOnNodes(context.Context, *UpdateProxyEnvsOnNodesRequest) (*UpdateProxyEnvsOnNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProxyEnvsOnNodes not implemented")
}
func (UnimplementedAnsiblerServiceServer) UpdateNoProxyEnvsInKubernetes(context.Context, *UpdateNoProxyEnvsInKubernetesRequest) (*UpdateNoProxyEnvsInKubernetesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNoProxyEnvsInKubernetes not implemented")
}
func (UnimplementedAnsiblerServiceServer) RemoveClaudieUtilities(context.Context, *RemoveClaudieUtilitiesRequest) (*RemoveClaudieUtilitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveClaudieUtilities not implemented")
}
func (UnimplementedAnsiblerServiceServer) mustEmbedUnimplementedAnsiblerServiceServer() {}
func (UnimplementedAnsiblerServiceServer) testEmbeddedByValue()                         {}

// UnsafeAnsiblerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnsiblerServiceServer will
// result in compilation errors.
type UnsafeAnsiblerServiceServer interface {
	mustEmbedUnimplementedAnsiblerServiceServer()
}

func RegisterAnsiblerServiceServer(s grpc.ServiceRegistrar, srv AnsiblerServiceServer) {
	// If the following call pancis, it indicates UnimplementedAnsiblerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AnsiblerService_ServiceDesc, srv)
}

func _AnsiblerService_InstallNodeRequirements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).InstallNodeRequirements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_InstallNodeRequirements_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).InstallNodeRequirements(ctx, req.(*InstallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_InstallVPN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).InstallVPN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_InstallVPN_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).InstallVPN(ctx, req.(*InstallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_SetUpLoadbalancers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUpLBRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).SetUpLoadbalancers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_SetUpLoadbalancers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).SetUpLoadbalancers(ctx, req.(*SetUpLBRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_DetermineApiEndpointChange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetermineApiEndpointChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).DetermineApiEndpointChange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_DetermineApiEndpointChange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).DetermineApiEndpointChange(ctx, req.(*DetermineApiEndpointChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_UpdateAPIEndpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAPIEndpointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).UpdateAPIEndpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_UpdateAPIEndpoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).UpdateAPIEndpoint(ctx, req.(*UpdateAPIEndpointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_UpdateProxyEnvsOnNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProxyEnvsOnNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).UpdateProxyEnvsOnNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_UpdateProxyEnvsOnNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).UpdateProxyEnvsOnNodes(ctx, req.(*UpdateProxyEnvsOnNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_UpdateNoProxyEnvsInKubernetes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNoProxyEnvsInKubernetesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).UpdateNoProxyEnvsInKubernetes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_UpdateNoProxyEnvsInKubernetes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).UpdateNoProxyEnvsInKubernetes(ctx, req.(*UpdateNoProxyEnvsInKubernetesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnsiblerService_RemoveClaudieUtilities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveClaudieUtilitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnsiblerServiceServer).RemoveClaudieUtilities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnsiblerService_RemoveClaudieUtilities_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnsiblerServiceServer).RemoveClaudieUtilities(ctx, req.(*RemoveClaudieUtilitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AnsiblerService_ServiceDesc is the grpc.ServiceDesc for AnsiblerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AnsiblerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "claudie.AnsiblerService",
	HandlerType: (*AnsiblerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InstallNodeRequirements",
			Handler:    _AnsiblerService_InstallNodeRequirements_Handler,
		},
		{
			MethodName: "InstallVPN",
			Handler:    _AnsiblerService_InstallVPN_Handler,
		},
		{
			MethodName: "SetUpLoadbalancers",
			Handler:    _AnsiblerService_SetUpLoadbalancers_Handler,
		},
		{
			MethodName: "DetermineApiEndpointChange",
			Handler:    _AnsiblerService_DetermineApiEndpointChange_Handler,
		},
		{
			MethodName: "UpdateAPIEndpoint",
			Handler:    _AnsiblerService_UpdateAPIEndpoint_Handler,
		},
		{
			MethodName: "UpdateProxyEnvsOnNodes",
			Handler:    _AnsiblerService_UpdateProxyEnvsOnNodes_Handler,
		},
		{
			MethodName: "UpdateNoProxyEnvsInKubernetes",
			Handler:    _AnsiblerService_UpdateNoProxyEnvsInKubernetes_Handler,
		},
		{
			MethodName: "RemoveClaudieUtilities",
			Handler:    _AnsiblerService_RemoveClaudieUtilities_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ansibler.proto",
}

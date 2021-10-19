// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0
// source: terraformer.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BuildInfrastructureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *BuildInfrastructureRequest) Reset() {
	*x = BuildInfrastructureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_terraformer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildInfrastructureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildInfrastructureRequest) ProtoMessage() {}

func (x *BuildInfrastructureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_terraformer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildInfrastructureRequest.ProtoReflect.Descriptor instead.
func (*BuildInfrastructureRequest) Descriptor() ([]byte, []int) {
	return file_terraformer_proto_rawDescGZIP(), []int{0}
}

func (x *BuildInfrastructureRequest) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type BuildInfrastructureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *BuildInfrastructureResponse) Reset() {
	*x = BuildInfrastructureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_terraformer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildInfrastructureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildInfrastructureResponse) ProtoMessage() {}

func (x *BuildInfrastructureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_terraformer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildInfrastructureResponse.ProtoReflect.Descriptor instead.
func (*BuildInfrastructureResponse) Descriptor() ([]byte, []int) {
	return file_terraformer_proto_rawDescGZIP(), []int{1}
}

func (x *BuildInfrastructureResponse) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type DestroyInfrastructureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *DestroyInfrastructureRequest) Reset() {
	*x = DestroyInfrastructureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_terraformer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DestroyInfrastructureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestroyInfrastructureRequest) ProtoMessage() {}

func (x *DestroyInfrastructureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_terraformer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestroyInfrastructureRequest.ProtoReflect.Descriptor instead.
func (*DestroyInfrastructureRequest) Descriptor() ([]byte, []int) {
	return file_terraformer_proto_rawDescGZIP(), []int{2}
}

func (x *DestroyInfrastructureRequest) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

type DestroyInfrastructureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config *Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *DestroyInfrastructureResponse) Reset() {
	*x = DestroyInfrastructureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_terraformer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DestroyInfrastructureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestroyInfrastructureResponse) ProtoMessage() {}

func (x *DestroyInfrastructureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_terraformer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestroyInfrastructureResponse.ProtoReflect.Descriptor instead.
func (*DestroyInfrastructureResponse) Descriptor() ([]byte, []int) {
	return file_terraformer_proto_rawDescGZIP(), []int{3}
}

func (x *DestroyInfrastructureResponse) GetConfig() *Config {
	if x != nil {
		return x.Config
	}
	return nil
}

var File_terraformer_proto protoreflect.FileDescriptor

var file_terraformer_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x1a, 0x0c, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x1a, 0x42,
	0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75,
	0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x06, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x22, 0x47, 0x0a, 0x1b, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x72,
	0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x48, 0x0a, 0x1c,
	0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x49, 0x0a, 0x1d, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f,
	0x79, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x32, 0xe2, 0x01, 0x0a, 0x12, 0x54, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x62, 0x0a, 0x13, 0x42, 0x75, 0x69, 0x6c,
	0x64, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x24, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x68, 0x0a, 0x15,
	0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x75, 0x72, 0x65, 0x12, 0x26, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79,
	0x49, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_terraformer_proto_rawDescOnce sync.Once
	file_terraformer_proto_rawDescData = file_terraformer_proto_rawDesc
)

func file_terraformer_proto_rawDescGZIP() []byte {
	file_terraformer_proto_rawDescOnce.Do(func() {
		file_terraformer_proto_rawDescData = protoimpl.X.CompressGZIP(file_terraformer_proto_rawDescData)
	})
	return file_terraformer_proto_rawDescData
}

var file_terraformer_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_terraformer_proto_goTypes = []interface{}{
	(*BuildInfrastructureRequest)(nil),    // 0: platform.BuildInfrastructureRequest
	(*BuildInfrastructureResponse)(nil),   // 1: platform.BuildInfrastructureResponse
	(*DestroyInfrastructureRequest)(nil),  // 2: platform.DestroyInfrastructureRequest
	(*DestroyInfrastructureResponse)(nil), // 3: platform.DestroyInfrastructureResponse
	(*Config)(nil),                        // 4: platform.Config
}
var file_terraformer_proto_depIdxs = []int32{
	4, // 0: platform.BuildInfrastructureRequest.config:type_name -> platform.Config
	4, // 1: platform.BuildInfrastructureResponse.config:type_name -> platform.Config
	4, // 2: platform.DestroyInfrastructureRequest.config:type_name -> platform.Config
	4, // 3: platform.DestroyInfrastructureResponse.config:type_name -> platform.Config
	0, // 4: platform.TerraformerService.BuildInfrastructure:input_type -> platform.BuildInfrastructureRequest
	2, // 5: platform.TerraformerService.DestroyInfrastructure:input_type -> platform.DestroyInfrastructureRequest
	1, // 6: platform.TerraformerService.BuildInfrastructure:output_type -> platform.BuildInfrastructureResponse
	3, // 7: platform.TerraformerService.DestroyInfrastructure:output_type -> platform.DestroyInfrastructureResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_terraformer_proto_init() }
func file_terraformer_proto_init() {
	if File_terraformer_proto != nil {
		return
	}
	file_config_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_terraformer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildInfrastructureRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_terraformer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildInfrastructureResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_terraformer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DestroyInfrastructureRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_terraformer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DestroyInfrastructureResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_terraformer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_terraformer_proto_goTypes,
		DependencyIndexes: file_terraformer_proto_depIdxs,
		MessageInfos:      file_terraformer_proto_msgTypes,
	}.Build()
	File_terraformer_proto = out.File
	file_terraformer_proto_rawDesc = nil
	file_terraformer_proto_goTypes = nil
	file_terraformer_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TerraformerServiceClient is the client API for TerraformerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TerraformerServiceClient interface {
	BuildInfrastructure(ctx context.Context, in *BuildInfrastructureRequest, opts ...grpc.CallOption) (*BuildInfrastructureResponse, error)
	DestroyInfrastructure(ctx context.Context, in *DestroyInfrastructureRequest, opts ...grpc.CallOption) (*DestroyInfrastructureResponse, error)
}

type terraformerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTerraformerServiceClient(cc grpc.ClientConnInterface) TerraformerServiceClient {
	return &terraformerServiceClient{cc}
}

func (c *terraformerServiceClient) BuildInfrastructure(ctx context.Context, in *BuildInfrastructureRequest, opts ...grpc.CallOption) (*BuildInfrastructureResponse, error) {
	out := new(BuildInfrastructureResponse)
	err := c.cc.Invoke(ctx, "/platform.TerraformerService/BuildInfrastructure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terraformerServiceClient) DestroyInfrastructure(ctx context.Context, in *DestroyInfrastructureRequest, opts ...grpc.CallOption) (*DestroyInfrastructureResponse, error) {
	out := new(DestroyInfrastructureResponse)
	err := c.cc.Invoke(ctx, "/platform.TerraformerService/DestroyInfrastructure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TerraformerServiceServer is the server API for TerraformerService service.
type TerraformerServiceServer interface {
	BuildInfrastructure(context.Context, *BuildInfrastructureRequest) (*BuildInfrastructureResponse, error)
	DestroyInfrastructure(context.Context, *DestroyInfrastructureRequest) (*DestroyInfrastructureResponse, error)
}

// UnimplementedTerraformerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTerraformerServiceServer struct {
}

func (*UnimplementedTerraformerServiceServer) BuildInfrastructure(context.Context, *BuildInfrastructureRequest) (*BuildInfrastructureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildInfrastructure not implemented")
}
func (*UnimplementedTerraformerServiceServer) DestroyInfrastructure(context.Context, *DestroyInfrastructureRequest) (*DestroyInfrastructureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroyInfrastructure not implemented")
}

func RegisterTerraformerServiceServer(s *grpc.Server, srv TerraformerServiceServer) {
	s.RegisterService(&_TerraformerService_serviceDesc, srv)
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
		FullMethod: "/platform.TerraformerService/BuildInfrastructure",
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
		FullMethod: "/platform.TerraformerService/DestroyInfrastructure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerraformerServiceServer).DestroyInfrastructure(ctx, req.(*DestroyInfrastructureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TerraformerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "platform.TerraformerService",
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

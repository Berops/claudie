// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.27.1
// source: kubeEleven.proto

package pb

import (
	spec "github.com/berops/claudie/proto/pb/spec"
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

type BuildClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Desired     *spec.K8Scluster `protobuf:"bytes,1,opt,name=desired,proto3" json:"desired,omitempty"`
	ProxyEnvs   *spec.ProxyEnvs  `protobuf:"bytes,3,opt,name=proxyEnvs,proto3" json:"proxyEnvs,omitempty"`
	ProjectName string           `protobuf:"bytes,4,opt,name=projectName,proto3" json:"projectName,omitempty"`
	// Endpoint specifies if the endpoint
	// is on a loadbalancer. If empty the
	// endpoint is one of the nodes supplied
	// in part of the desired state.
	LoadBalancerEndpoint string `protobuf:"bytes,5,opt,name=loadBalancerEndpoint,proto3" json:"loadBalancerEndpoint,omitempty"`
}

func (x *BuildClusterRequest) Reset() {
	*x = BuildClusterRequest{}
	mi := &file_kubeEleven_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BuildClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildClusterRequest) ProtoMessage() {}

func (x *BuildClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kubeEleven_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildClusterRequest.ProtoReflect.Descriptor instead.
func (*BuildClusterRequest) Descriptor() ([]byte, []int) {
	return file_kubeEleven_proto_rawDescGZIP(), []int{0}
}

func (x *BuildClusterRequest) GetDesired() *spec.K8Scluster {
	if x != nil {
		return x.Desired
	}
	return nil
}

func (x *BuildClusterRequest) GetProxyEnvs() *spec.ProxyEnvs {
	if x != nil {
		return x.ProxyEnvs
	}
	return nil
}

func (x *BuildClusterRequest) GetProjectName() string {
	if x != nil {
		return x.ProjectName
	}
	return ""
}

func (x *BuildClusterRequest) GetLoadBalancerEndpoint() string {
	if x != nil {
		return x.LoadBalancerEndpoint
	}
	return ""
}

type BuildClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Desired *spec.K8Scluster `protobuf:"bytes,1,opt,name=desired,proto3" json:"desired,omitempty"`
}

func (x *BuildClusterResponse) Reset() {
	*x = BuildClusterResponse{}
	mi := &file_kubeEleven_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BuildClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildClusterResponse) ProtoMessage() {}

func (x *BuildClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kubeEleven_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildClusterResponse.ProtoReflect.Descriptor instead.
func (*BuildClusterResponse) Descriptor() ([]byte, []int) {
	return file_kubeEleven_proto_rawDescGZIP(), []int{1}
}

func (x *BuildClusterResponse) GetDesired() *spec.K8Scluster {
	if x != nil {
		return x.Desired
	}
	return nil
}

type DestroyClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectName string           `protobuf:"bytes,1,opt,name=projectName,proto3" json:"projectName,omitempty"`
	Current     *spec.K8Scluster `protobuf:"bytes,2,opt,name=current,proto3" json:"current,omitempty"`
	// Endpoint specifies if the endpoint
	// is on a loadbalancer. If empty the
	// endpoint is one of the nodes supplied
	// in part of the desired state.
	LoadBalancerEndpoint string `protobuf:"bytes,3,opt,name=loadBalancerEndpoint,proto3" json:"loadBalancerEndpoint,omitempty"`
}

func (x *DestroyClusterRequest) Reset() {
	*x = DestroyClusterRequest{}
	mi := &file_kubeEleven_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DestroyClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestroyClusterRequest) ProtoMessage() {}

func (x *DestroyClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kubeEleven_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestroyClusterRequest.ProtoReflect.Descriptor instead.
func (*DestroyClusterRequest) Descriptor() ([]byte, []int) {
	return file_kubeEleven_proto_rawDescGZIP(), []int{2}
}

func (x *DestroyClusterRequest) GetProjectName() string {
	if x != nil {
		return x.ProjectName
	}
	return ""
}

func (x *DestroyClusterRequest) GetCurrent() *spec.K8Scluster {
	if x != nil {
		return x.Current
	}
	return nil
}

func (x *DestroyClusterRequest) GetLoadBalancerEndpoint() string {
	if x != nil {
		return x.LoadBalancerEndpoint
	}
	return ""
}

type DestroyClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Current *spec.K8Scluster `protobuf:"bytes,1,opt,name=current,proto3" json:"current,omitempty"`
}

func (x *DestroyClusterResponse) Reset() {
	*x = DestroyClusterResponse{}
	mi := &file_kubeEleven_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DestroyClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestroyClusterResponse) ProtoMessage() {}

func (x *DestroyClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kubeEleven_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestroyClusterResponse.ProtoReflect.Descriptor instead.
func (*DestroyClusterResponse) Descriptor() ([]byte, []int) {
	return file_kubeEleven_proto_rawDescGZIP(), []int{3}
}

func (x *DestroyClusterResponse) GetCurrent() *spec.K8Scluster {
	if x != nil {
		return x.Current
	}
	return nil
}

var File_kubeEleven_proto protoreflect.FileDescriptor

var file_kubeEleven_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6b, 0x75, 0x62, 0x65, 0x45, 0x6c, 0x65, 0x76, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x63, 0x6c, 0x61, 0x75, 0x64, 0x69, 0x65, 0x1a, 0x13, 0x73, 0x70, 0x65,
	0x63, 0x2f, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xc6, 0x01, 0x0a, 0x13, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x07, 0x64, 0x65, 0x73, 0x69,
	0x72, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x70, 0x65, 0x63,
	0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x64, 0x65, 0x73,
	0x69, 0x72, 0x65, 0x64, 0x12, 0x2d, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x45, 0x6e, 0x76,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x45, 0x6e, 0x76, 0x73, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x45,
	0x6e, 0x76, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x14, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x14, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x42, 0x0a, 0x14, 0x42, 0x75, 0x69,
	0x6c, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2a, 0x0a, 0x07, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x22, 0x99, 0x01,
	0x0a, 0x15, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x70, 0x65,
	0x63, 0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x14, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x14, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x44, 0x0a, 0x16, 0x44, 0x65, 0x73,
	0x74, 0x72, 0x6f, 0x79, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x4b, 0x38, 0x73, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x32,
	0xb3, 0x01, 0x0a, 0x11, 0x4b, 0x75, 0x62, 0x65, 0x45, 0x6c, 0x65, 0x76, 0x65, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x0c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1c, 0x2e, 0x63, 0x6c, 0x61, 0x75, 0x64, 0x69, 0x65, 0x2e,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x6c, 0x61, 0x75, 0x64, 0x69, 0x65, 0x2e, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x12, 0x1e, 0x2e, 0x63, 0x6c, 0x61, 0x75, 0x64, 0x69, 0x65, 0x2e, 0x44,
	0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6c, 0x61, 0x75, 0x64, 0x69, 0x65, 0x2e, 0x44,
	0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kubeEleven_proto_rawDescOnce sync.Once
	file_kubeEleven_proto_rawDescData = file_kubeEleven_proto_rawDesc
)

func file_kubeEleven_proto_rawDescGZIP() []byte {
	file_kubeEleven_proto_rawDescOnce.Do(func() {
		file_kubeEleven_proto_rawDescData = protoimpl.X.CompressGZIP(file_kubeEleven_proto_rawDescData)
	})
	return file_kubeEleven_proto_rawDescData
}

var file_kubeEleven_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_kubeEleven_proto_goTypes = []any{
	(*BuildClusterRequest)(nil),    // 0: claudie.BuildClusterRequest
	(*BuildClusterResponse)(nil),   // 1: claudie.BuildClusterResponse
	(*DestroyClusterRequest)(nil),  // 2: claudie.DestroyClusterRequest
	(*DestroyClusterResponse)(nil), // 3: claudie.DestroyClusterResponse
	(*spec.K8Scluster)(nil),        // 4: spec.K8scluster
	(*spec.ProxyEnvs)(nil),         // 5: spec.ProxyEnvs
}
var file_kubeEleven_proto_depIdxs = []int32{
	4, // 0: claudie.BuildClusterRequest.desired:type_name -> spec.K8scluster
	5, // 1: claudie.BuildClusterRequest.proxyEnvs:type_name -> spec.ProxyEnvs
	4, // 2: claudie.BuildClusterResponse.desired:type_name -> spec.K8scluster
	4, // 3: claudie.DestroyClusterRequest.current:type_name -> spec.K8scluster
	4, // 4: claudie.DestroyClusterResponse.current:type_name -> spec.K8scluster
	0, // 5: claudie.KubeElevenService.BuildCluster:input_type -> claudie.BuildClusterRequest
	2, // 6: claudie.KubeElevenService.DestroyCluster:input_type -> claudie.DestroyClusterRequest
	1, // 7: claudie.KubeElevenService.BuildCluster:output_type -> claudie.BuildClusterResponse
	3, // 8: claudie.KubeElevenService.DestroyCluster:output_type -> claudie.DestroyClusterResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_kubeEleven_proto_init() }
func file_kubeEleven_proto_init() {
	if File_kubeEleven_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kubeEleven_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kubeEleven_proto_goTypes,
		DependencyIndexes: file_kubeEleven_proto_depIdxs,
		MessageInfos:      file_kubeEleven_proto_msgTypes,
	}.Build()
	File_kubeEleven_proto = out.File
	file_kubeEleven_proto_rawDesc = nil
	file_kubeEleven_proto_goTypes = nil
	file_kubeEleven_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: proto/kuber.proto

package pb

import (
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

type SetUpStorageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredState *Project `protobuf:"bytes,1,opt,name=desiredState,proto3" json:"desiredState,omitempty"`
}

func (x *SetUpStorageRequest) Reset() {
	*x = SetUpStorageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUpStorageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUpStorageRequest) ProtoMessage() {}

func (x *SetUpStorageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetUpStorageRequest.ProtoReflect.Descriptor instead.
func (*SetUpStorageRequest) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{0}
}

func (x *SetUpStorageRequest) GetDesiredState() *Project {
	if x != nil {
		return x.DesiredState
	}
	return nil
}

type SetUpStorageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredState *Project `protobuf:"bytes,1,opt,name=desiredState,proto3" json:"desiredState,omitempty"`
	ErrorMessage string   `protobuf:"bytes,2,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
}

func (x *SetUpStorageResponse) Reset() {
	*x = SetUpStorageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUpStorageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUpStorageResponse) ProtoMessage() {}

func (x *SetUpStorageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetUpStorageResponse.ProtoReflect.Descriptor instead.
func (*SetUpStorageResponse) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{1}
}

func (x *SetUpStorageResponse) GetDesiredState() *Project {
	if x != nil {
		return x.DesiredState
	}
	return nil
}

func (x *SetUpStorageResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type StoreKubeconfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cluster *K8Scluster `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
}

func (x *StoreKubeconfigRequest) Reset() {
	*x = StoreKubeconfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreKubeconfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreKubeconfigRequest) ProtoMessage() {}

func (x *StoreKubeconfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreKubeconfigRequest.ProtoReflect.Descriptor instead.
func (*StoreKubeconfigRequest) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{2}
}

func (x *StoreKubeconfigRequest) GetCluster() *K8Scluster {
	if x != nil {
		return x.Cluster
	}
	return nil
}

type StoreKubeconfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage string `protobuf:"bytes,1,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
}

func (x *StoreKubeconfigResponse) Reset() {
	*x = StoreKubeconfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreKubeconfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreKubeconfigResponse) ProtoMessage() {}

func (x *StoreKubeconfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreKubeconfigResponse.ProtoReflect.Descriptor instead.
func (*StoreKubeconfigResponse) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{3}
}

func (x *StoreKubeconfigResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type DeleteKubeconfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cluster *K8Scluster `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
}

func (x *DeleteKubeconfigRequest) Reset() {
	*x = DeleteKubeconfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteKubeconfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKubeconfigRequest) ProtoMessage() {}

func (x *DeleteKubeconfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKubeconfigRequest.ProtoReflect.Descriptor instead.
func (*DeleteKubeconfigRequest) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteKubeconfigRequest) GetCluster() *K8Scluster {
	if x != nil {
		return x.Cluster
	}
	return nil
}

type DeleteKubeconfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage string `protobuf:"bytes,1,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
}

func (x *DeleteKubeconfigResponse) Reset() {
	*x = DeleteKubeconfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteKubeconfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKubeconfigResponse) ProtoMessage() {}

func (x *DeleteKubeconfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKubeconfigResponse.ProtoReflect.Descriptor instead.
func (*DeleteKubeconfigResponse) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteKubeconfigResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type DeleteNodesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cluster     *K8Scluster `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	WorkerNodes []string    `protobuf:"bytes,2,rep,name=workerNodes,proto3" json:"workerNodes,omitempty"`
	MasterNodes []string    `protobuf:"bytes,3,rep,name=masterNodes,proto3" json:"masterNodes,omitempty"`
}

func (x *DeleteNodesRequest) Reset() {
	*x = DeleteNodesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNodesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNodesRequest) ProtoMessage() {}

func (x *DeleteNodesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNodesRequest.ProtoReflect.Descriptor instead.
func (*DeleteNodesRequest) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteNodesRequest) GetCluster() *K8Scluster {
	if x != nil {
		return x.Cluster
	}
	return nil
}

func (x *DeleteNodesRequest) GetWorkerNodes() []string {
	if x != nil {
		return x.WorkerNodes
	}
	return nil
}

func (x *DeleteNodesRequest) GetMasterNodes() []string {
	if x != nil {
		return x.MasterNodes
	}
	return nil
}

type DeleteNodesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cluster      *K8Scluster `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	ErrorMessage string      `protobuf:"bytes,2,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
}

func (x *DeleteNodesResponse) Reset() {
	*x = DeleteNodesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_kuber_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNodesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNodesResponse) ProtoMessage() {}

func (x *DeleteNodesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_kuber_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNodesResponse.ProtoReflect.Descriptor instead.
func (*DeleteNodesResponse) Descriptor() ([]byte, []int) {
	return file_proto_kuber_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteNodesResponse) GetCluster() *K8Scluster {
	if x != nil {
		return x.Cluster
	}
	return nil
}

func (x *DeleteNodesResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

var File_proto_kuber_proto protoreflect.FileDescriptor

var file_proto_kuber_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x1a, 0x12, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x4c, 0x0a, 0x13, 0x53, 0x65, 0x74, 0x55, 0x70, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x69,
	0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x52, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22,
	0x71, 0x0a, 0x14, 0x53, 0x65, 0x74, 0x55, 0x70, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72,
	0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x52, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x22,
	0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x48, 0x0a, 0x16, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x22, 0x3d, 0x0a, 0x17,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x49, 0x0a, 0x17, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x22, 0x3e, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a,
	0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a,
	0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12,
	0x20, 0x0a, 0x0b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x22, 0x69, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x4b, 0x38, 0x73, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xdc, 0x02, 0x0a,
	0x0c, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a,
	0x0c, 0x53, 0x65, 0x74, 0x55, 0x70, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x70, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x70, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x0f,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x20, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x21, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x75,
	0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x21, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x75, 0x62,
	0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x4a, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x1c,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f,
	0x64, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_kuber_proto_rawDescOnce sync.Once
	file_proto_kuber_proto_rawDescData = file_proto_kuber_proto_rawDesc
)

func file_proto_kuber_proto_rawDescGZIP() []byte {
	file_proto_kuber_proto_rawDescOnce.Do(func() {
		file_proto_kuber_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_kuber_proto_rawDescData)
	})
	return file_proto_kuber_proto_rawDescData
}

var file_proto_kuber_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_kuber_proto_goTypes = []interface{}{
	(*SetUpStorageRequest)(nil),      // 0: platform.SetUpStorageRequest
	(*SetUpStorageResponse)(nil),     // 1: platform.SetUpStorageResponse
	(*StoreKubeconfigRequest)(nil),   // 2: platform.StoreKubeconfigRequest
	(*StoreKubeconfigResponse)(nil),  // 3: platform.StoreKubeconfigResponse
	(*DeleteKubeconfigRequest)(nil),  // 4: platform.DeleteKubeconfigRequest
	(*DeleteKubeconfigResponse)(nil), // 5: platform.DeleteKubeconfigResponse
	(*DeleteNodesRequest)(nil),       // 6: platform.DeleteNodesRequest
	(*DeleteNodesResponse)(nil),      // 7: platform.DeleteNodesResponse
	(*Project)(nil),                  // 8: platform.Project
	(*K8Scluster)(nil),               // 9: platform.K8scluster
}
var file_proto_kuber_proto_depIdxs = []int32{
	8,  // 0: platform.SetUpStorageRequest.desiredState:type_name -> platform.Project
	8,  // 1: platform.SetUpStorageResponse.desiredState:type_name -> platform.Project
	9,  // 2: platform.StoreKubeconfigRequest.cluster:type_name -> platform.K8scluster
	9,  // 3: platform.DeleteKubeconfigRequest.cluster:type_name -> platform.K8scluster
	9,  // 4: platform.DeleteNodesRequest.cluster:type_name -> platform.K8scluster
	9,  // 5: platform.DeleteNodesResponse.cluster:type_name -> platform.K8scluster
	0,  // 6: platform.KuberService.SetUpStorage:input_type -> platform.SetUpStorageRequest
	2,  // 7: platform.KuberService.StoreKubeconfig:input_type -> platform.StoreKubeconfigRequest
	4,  // 8: platform.KuberService.DeleteKubeconfig:input_type -> platform.DeleteKubeconfigRequest
	6,  // 9: platform.KuberService.DeleteNodes:input_type -> platform.DeleteNodesRequest
	1,  // 10: platform.KuberService.SetUpStorage:output_type -> platform.SetUpStorageResponse
	3,  // 11: platform.KuberService.StoreKubeconfig:output_type -> platform.StoreKubeconfigResponse
	5,  // 12: platform.KuberService.DeleteKubeconfig:output_type -> platform.DeleteKubeconfigResponse
	7,  // 13: platform.KuberService.DeleteNodes:output_type -> platform.DeleteNodesResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_proto_kuber_proto_init() }
func file_proto_kuber_proto_init() {
	if File_proto_kuber_proto != nil {
		return
	}
	file_proto_config_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_kuber_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetUpStorageRequest); i {
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
		file_proto_kuber_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetUpStorageResponse); i {
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
		file_proto_kuber_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreKubeconfigRequest); i {
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
		file_proto_kuber_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreKubeconfigResponse); i {
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
		file_proto_kuber_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteKubeconfigRequest); i {
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
		file_proto_kuber_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteKubeconfigResponse); i {
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
		file_proto_kuber_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNodesRequest); i {
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
		file_proto_kuber_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNodesResponse); i {
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
			RawDescriptor: file_proto_kuber_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_kuber_proto_goTypes,
		DependencyIndexes: file_proto_kuber_proto_depIdxs,
		MessageInfos:      file_proto_kuber_proto_msgTypes,
	}.Build()
	File_proto_kuber_proto = out.File
	file_proto_kuber_proto_rawDesc = nil
	file_proto_kuber_proto_goTypes = nil
	file_proto_kuber_proto_depIdxs = nil
}

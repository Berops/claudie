// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.4
// source: proto/ansibler.proto

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

type InstallRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredState *Project `protobuf:"bytes,1,opt,name=desiredState,proto3" json:"desiredState,omitempty"`
	CurrentState *Project `protobuf:"bytes,2,opt,name=currentState,proto3" json:"currentState,omitempty"`
}

func (x *InstallRequest) Reset() {
	*x = InstallRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_ansibler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallRequest) ProtoMessage() {}

func (x *InstallRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_ansibler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallRequest.ProtoReflect.Descriptor instead.
func (*InstallRequest) Descriptor() ([]byte, []int) {
	return file_proto_ansibler_proto_rawDescGZIP(), []int{0}
}

func (x *InstallRequest) GetDesiredState() *Project {
	if x != nil {
		return x.DesiredState
	}
	return nil
}

func (x *InstallRequest) GetCurrentState() *Project {
	if x != nil {
		return x.CurrentState
	}
	return nil
}

type InstallResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredState *Project `protobuf:"bytes,1,opt,name=desiredState,proto3" json:"desiredState,omitempty"`
}

func (x *InstallResponse) Reset() {
	*x = InstallResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_ansibler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallResponse) ProtoMessage() {}

func (x *InstallResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_ansibler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallResponse.ProtoReflect.Descriptor instead.
func (*InstallResponse) Descriptor() ([]byte, []int) {
	return file_proto_ansibler_proto_rawDescGZIP(), []int{1}
}

func (x *InstallResponse) GetDesiredState() *Project {
	if x != nil {
		return x.DesiredState
	}
	return nil
}

type SetUpLBRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredState *Project `protobuf:"bytes,1,opt,name=desiredState,proto3" json:"desiredState,omitempty"`
	CurrentState *Project `protobuf:"bytes,2,opt,name=currentState,proto3" json:"currentState,omitempty"`
}

func (x *SetUpLBRequest) Reset() {
	*x = SetUpLBRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_ansibler_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUpLBRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUpLBRequest) ProtoMessage() {}

func (x *SetUpLBRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_ansibler_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetUpLBRequest.ProtoReflect.Descriptor instead.
func (*SetUpLBRequest) Descriptor() ([]byte, []int) {
	return file_proto_ansibler_proto_rawDescGZIP(), []int{2}
}

func (x *SetUpLBRequest) GetDesiredState() *Project {
	if x != nil {
		return x.DesiredState
	}
	return nil
}

func (x *SetUpLBRequest) GetCurrentState() *Project {
	if x != nil {
		return x.CurrentState
	}
	return nil
}

type SetUpLBResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredState *Project `protobuf:"bytes,1,opt,name=desiredState,proto3" json:"desiredState,omitempty"`
}

func (x *SetUpLBResponse) Reset() {
	*x = SetUpLBResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_ansibler_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUpLBResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUpLBResponse) ProtoMessage() {}

func (x *SetUpLBResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_ansibler_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetUpLBResponse.ProtoReflect.Descriptor instead.
func (*SetUpLBResponse) Descriptor() ([]byte, []int) {
	return file_proto_ansibler_proto_rawDescGZIP(), []int{3}
}

func (x *SetUpLBResponse) GetDesiredState() *Project {
	if x != nil {
		return x.DesiredState
	}
	return nil
}

var File_proto_ansibler_proto protoreflect.FileDescriptor

var file_proto_ansibler_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x6e, 0x73, 0x69, 0x62, 0x6c, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x1a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7e, 0x0a, 0x0e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52,
	0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x35, 0x0a,
	0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x22, 0x48, 0x0a, 0x0f, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72,
	0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x52, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x7e,
	0x0a, 0x0e, 0x53, 0x65, 0x74, 0x55, 0x70, 0x4c, 0x42, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x35, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72,
	0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x48,
	0x0a, 0x0f, 0x53, 0x65, 0x74, 0x55, 0x70, 0x4c, 0x42, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x35, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0c, 0x64, 0x65, 0x73, 0x69,
	0x72, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x32, 0xef, 0x01, 0x0a, 0x0f, 0x41, 0x6e, 0x73,
	0x69, 0x62, 0x6c, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x17,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x18, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x19, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x56, 0x50, 0x4e, 0x12, 0x18, 0x2e, 0x70, 0x6c, 0x61,
	0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x49, 0x0a, 0x12, 0x53, 0x65, 0x74, 0x55, 0x70, 0x4c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x72, 0x73, 0x12, 0x18, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2e, 0x53, 0x65, 0x74, 0x55, 0x70, 0x4c, 0x42, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x19, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x70,
	0x4c, 0x42, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_ansibler_proto_rawDescOnce sync.Once
	file_proto_ansibler_proto_rawDescData = file_proto_ansibler_proto_rawDesc
)

func file_proto_ansibler_proto_rawDescGZIP() []byte {
	file_proto_ansibler_proto_rawDescOnce.Do(func() {
		file_proto_ansibler_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_ansibler_proto_rawDescData)
	})
	return file_proto_ansibler_proto_rawDescData
}

var file_proto_ansibler_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_ansibler_proto_goTypes = []interface{}{
	(*InstallRequest)(nil),  // 0: platform.InstallRequest
	(*InstallResponse)(nil), // 1: platform.InstallResponse
	(*SetUpLBRequest)(nil),  // 2: platform.SetUpLBRequest
	(*SetUpLBResponse)(nil), // 3: platform.SetUpLBResponse
	(*Project)(nil),         // 4: platform.Project
}
var file_proto_ansibler_proto_depIdxs = []int32{
	4, // 0: platform.InstallRequest.desiredState:type_name -> platform.Project
	4, // 1: platform.InstallRequest.currentState:type_name -> platform.Project
	4, // 2: platform.InstallResponse.desiredState:type_name -> platform.Project
	4, // 3: platform.SetUpLBRequest.desiredState:type_name -> platform.Project
	4, // 4: platform.SetUpLBRequest.currentState:type_name -> platform.Project
	4, // 5: platform.SetUpLBResponse.desiredState:type_name -> platform.Project
	0, // 6: platform.AnsiblerService.InstallNodeRequirements:input_type -> platform.InstallRequest
	0, // 7: platform.AnsiblerService.InstallVPN:input_type -> platform.InstallRequest
	2, // 8: platform.AnsiblerService.SetUpLoadbalancers:input_type -> platform.SetUpLBRequest
	1, // 9: platform.AnsiblerService.InstallNodeRequirements:output_type -> platform.InstallResponse
	1, // 10: platform.AnsiblerService.InstallVPN:output_type -> platform.InstallResponse
	3, // 11: platform.AnsiblerService.SetUpLoadbalancers:output_type -> platform.SetUpLBResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_ansibler_proto_init() }
func file_proto_ansibler_proto_init() {
	if File_proto_ansibler_proto != nil {
		return
	}
	file_proto_config_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_ansibler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallRequest); i {
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
		file_proto_ansibler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallResponse); i {
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
		file_proto_ansibler_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetUpLBRequest); i {
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
		file_proto_ansibler_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetUpLBResponse); i {
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
			RawDescriptor: file_proto_ansibler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_ansibler_proto_goTypes,
		DependencyIndexes: file_proto_ansibler_proto_depIdxs,
		MessageInfos:      file_proto_ansibler_proto_msgTypes,
	}.Build()
	File_proto_ansibler_proto = out.File
	file_proto_ansibler_proto_rawDesc = nil
	file_proto_ansibler_proto_goTypes = nil
	file_proto_ansibler_proto_depIdxs = nil
}

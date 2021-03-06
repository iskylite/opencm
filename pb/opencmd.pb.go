// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: pb/opencmd.proto

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

// OpenCMADRegistry 注册opencmad服务
type OpenCMADRegistry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host       string       `protobuf:"bytes,1,opt,name=Host,proto3" json:"Host,omitempty"`
	OS         *OS          `protobuf:"bytes,2,opt,name=OS,proto3" json:"OS,omitempty"`
	Interfaces []*Interface `protobuf:"bytes,3,rep,name=Interfaces,proto3" json:"Interfaces,omitempty"`
}

func (x *OpenCMADRegistry) Reset() {
	*x = OpenCMADRegistry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_opencmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenCMADRegistry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenCMADRegistry) ProtoMessage() {}

func (x *OpenCMADRegistry) ProtoReflect() protoreflect.Message {
	mi := &file_pb_opencmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenCMADRegistry.ProtoReflect.Descriptor instead.
func (*OpenCMADRegistry) Descriptor() ([]byte, []int) {
	return file_pb_opencmd_proto_rawDescGZIP(), []int{0}
}

func (x *OpenCMADRegistry) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *OpenCMADRegistry) GetOS() *OS {
	if x != nil {
		return x.OS
	}
	return nil
}

func (x *OpenCMADRegistry) GetInterfaces() []*Interface {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

// OS 系统基本
type OS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BootTime       uint64 `protobuf:"varint,1,opt,name=BootTime,proto3" json:"BootTime,omitempty"`
	MemTotal       uint64 `protobuf:"varint,2,opt,name=MemTotal,proto3" json:"MemTotal,omitempty"` // G
	CPUNum         int32  `protobuf:"varint,3,opt,name=CPUNum,proto3" json:"CPUNum,omitempty"`
	Arch           string `protobuf:"bytes,4,opt,name=Arch,proto3" json:"Arch,omitempty"`
	Kernel         string `protobuf:"bytes,5,opt,name=Kernel,proto3" json:"Kernel,omitempty"`
	Version        string `protobuf:"bytes,6,opt,name=Version,proto3" json:"Version,omitempty"`
	ImageBuildTime string `protobuf:"bytes,7,opt,name=ImageBuildTime,proto3" json:"ImageBuildTime,omitempty"` // .ts
}

func (x *OS) Reset() {
	*x = OS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_opencmd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OS) ProtoMessage() {}

func (x *OS) ProtoReflect() protoreflect.Message {
	mi := &file_pb_opencmd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OS.ProtoReflect.Descriptor instead.
func (*OS) Descriptor() ([]byte, []int) {
	return file_pb_opencmd_proto_rawDescGZIP(), []int{1}
}

func (x *OS) GetBootTime() uint64 {
	if x != nil {
		return x.BootTime
	}
	return 0
}

func (x *OS) GetMemTotal() uint64 {
	if x != nil {
		return x.MemTotal
	}
	return 0
}

func (x *OS) GetCPUNum() int32 {
	if x != nil {
		return x.CPUNum
	}
	return 0
}

func (x *OS) GetArch() string {
	if x != nil {
		return x.Arch
	}
	return ""
}

func (x *OS) GetKernel() string {
	if x != nil {
		return x.Kernel
	}
	return ""
}

func (x *OS) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *OS) GetImageBuildTime() string {
	if x != nil {
		return x.ImageBuildTime
	}
	return ""
}

// Interface 网络配置
type Interface struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dev          string `protobuf:"bytes,1,opt,name=Dev,proto3" json:"Dev,omitempty"`
	HardwareAddr string `protobuf:"bytes,2,opt,name=HardwareAddr,proto3" json:"HardwareAddr,omitempty"`
	Flags        string `protobuf:"bytes,3,opt,name=Flags,proto3" json:"Flags,omitempty"`
	IP           string `protobuf:"bytes,4,opt,name=IP,proto3" json:"IP,omitempty"`
	Mask         int32  `protobuf:"varint,5,opt,name=Mask,proto3" json:"Mask,omitempty"`
}

func (x *Interface) Reset() {
	*x = Interface{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_opencmd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Interface) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Interface) ProtoMessage() {}

func (x *Interface) ProtoReflect() protoreflect.Message {
	mi := &file_pb_opencmd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Interface.ProtoReflect.Descriptor instead.
func (*Interface) Descriptor() ([]byte, []int) {
	return file_pb_opencmd_proto_rawDescGZIP(), []int{2}
}

func (x *Interface) GetDev() string {
	if x != nil {
		return x.Dev
	}
	return ""
}

func (x *Interface) GetHardwareAddr() string {
	if x != nil {
		return x.HardwareAddr
	}
	return ""
}

func (x *Interface) GetFlags() string {
	if x != nil {
		return x.Flags
	}
	return ""
}

func (x *Interface) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

func (x *Interface) GetMask() int32 {
	if x != nil {
		return x.Mask
	}
	return 0
}

// OpenCMADConfig opencmad注册响应数据
type OpenCMADConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectorFlags int32  `protobuf:"varint,1,opt,name=CollectorFlags,proto3" json:"CollectorFlags,omitempty"`
	NodeType       string `protobuf:"bytes,2,opt,name=NodeType,proto3" json:"NodeType,omitempty"`
}

func (x *OpenCMADConfig) Reset() {
	*x = OpenCMADConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_opencmd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenCMADConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenCMADConfig) ProtoMessage() {}

func (x *OpenCMADConfig) ProtoReflect() protoreflect.Message {
	mi := &file_pb_opencmd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenCMADConfig.ProtoReflect.Descriptor instead.
func (*OpenCMADConfig) Descriptor() ([]byte, []int) {
	return file_pb_opencmd_proto_rawDescGZIP(), []int{3}
}

func (x *OpenCMADConfig) GetCollectorFlags() int32 {
	if x != nil {
		return x.CollectorFlags
	}
	return 0
}

func (x *OpenCMADConfig) GetNodeType() string {
	if x != nil {
		return x.NodeType
	}
	return ""
}

var File_pb_opencmd_proto protoreflect.FileDescriptor

var file_pb_opencmd_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x62, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x63, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x10, 0x70, 0x62, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x10, 0x4f, 0x70, 0x65, 0x6e,
	0x43, 0x4d, 0x41, 0x44, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x48, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x6f, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x02, 0x4f, 0x53, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x70,
	0x62, 0x2e, 0x4f, 0x53, 0x52, 0x02, 0x4f, 0x53, 0x12, 0x2d, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70,
	0x62, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x0a, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x22, 0xc2, 0x01, 0x0a, 0x02, 0x4f, 0x53, 0x12, 0x1a,
	0x0a, 0x08, 0x42, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x08, 0x42, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x65,
	0x6d, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x4d, 0x65,
	0x6d, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x50, 0x55, 0x4e, 0x75, 0x6d,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x43, 0x50, 0x55, 0x4e, 0x75, 0x6d, 0x12, 0x12,
	0x0a, 0x04, 0x41, 0x72, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x41, 0x72,
	0x63, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x4b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x4b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x75, 0x69,
	0x6c, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x7b, 0x0a, 0x09,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x44, 0x65, 0x76,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x44, 0x65, 0x76, 0x12, 0x22, 0x0a, 0x0c, 0x48,
	0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x41, 0x64, 0x64, 0x72, 0x12,
	0x14, 0x0a, 0x05, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x50, 0x12, 0x12, 0x0a, 0x04, 0x4d, 0x61, 0x73, 0x6b, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x4d, 0x61, 0x73, 0x6b, 0x22, 0x54, 0x0a, 0x0e, 0x4f, 0x70, 0x65,
	0x6e, 0x43, 0x4d, 0x41, 0x44, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x26, 0x0a, 0x0e, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x46, 0x6c,
	0x61, 0x67, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x32,
	0x84, 0x01, 0x0a, 0x0e, 0x4f, 0x70, 0x65, 0x6e, 0x43, 0x4d, 0x44, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x3c, 0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70,
	0x65, 0x6e, 0x43, 0x4d, 0x41, 0x44, 0x12, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x4f, 0x70, 0x65, 0x6e,
	0x43, 0x4d, 0x41, 0x44, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x1a, 0x12, 0x2e, 0x70,
	0x62, 0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x43, 0x4d, 0x41, 0x44, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x34, 0x0a, 0x12, 0x55, 0x6e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x70,
	0x65, 0x6e, 0x43, 0x4d, 0x41, 0x44, 0x12, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x69, 0x63, 0x4d, 0x73, 0x67, 0x1a, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x69, 0x63, 0x4d, 0x73, 0x67, 0x42, 0x05, 0x5a, 0x03, 0x70, 0x62, 0x2f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_opencmd_proto_rawDescOnce sync.Once
	file_pb_opencmd_proto_rawDescData = file_pb_opencmd_proto_rawDesc
)

func file_pb_opencmd_proto_rawDescGZIP() []byte {
	file_pb_opencmd_proto_rawDescOnce.Do(func() {
		file_pb_opencmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_opencmd_proto_rawDescData)
	})
	return file_pb_opencmd_proto_rawDescData
}

var file_pb_opencmd_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pb_opencmd_proto_goTypes = []interface{}{
	(*OpenCMADRegistry)(nil), // 0: pb.OpenCMADRegistry
	(*OS)(nil),               // 1: pb.OS
	(*Interface)(nil),        // 2: pb.Interface
	(*OpenCMADConfig)(nil),   // 3: pb.OpenCMADConfig
	(*GenericMsg)(nil),       // 4: pb.GenericMsg
}
var file_pb_opencmd_proto_depIdxs = []int32{
	1, // 0: pb.OpenCMADRegistry.OS:type_name -> pb.OS
	2, // 1: pb.OpenCMADRegistry.Interfaces:type_name -> pb.Interface
	0, // 2: pb.OpenCMDService.RegisterOpenCMAD:input_type -> pb.OpenCMADRegistry
	4, // 3: pb.OpenCMDService.UnRegisterOpenCMAD:input_type -> pb.GenericMsg
	3, // 4: pb.OpenCMDService.RegisterOpenCMAD:output_type -> pb.OpenCMADConfig
	4, // 5: pb.OpenCMDService.UnRegisterOpenCMAD:output_type -> pb.GenericMsg
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pb_opencmd_proto_init() }
func file_pb_opencmd_proto_init() {
	if File_pb_opencmd_proto != nil {
		return
	}
	file_pb_generic_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pb_opencmd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenCMADRegistry); i {
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
		file_pb_opencmd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OS); i {
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
		file_pb_opencmd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Interface); i {
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
		file_pb_opencmd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenCMADConfig); i {
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
			RawDescriptor: file_pb_opencmd_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_opencmd_proto_goTypes,
		DependencyIndexes: file_pb_opencmd_proto_depIdxs,
		MessageInfos:      file_pb_opencmd_proto_msgTypes,
	}.Build()
	File_pb_opencmd_proto = out.File
	file_pb_opencmd_proto_rawDesc = nil
	file_pb_opencmd_proto_goTypes = nil
	file_pb_opencmd_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// OpenCMDServiceClient is the client API for OpenCMDService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OpenCMDServiceClient interface {
	RegisterOpenCMAD(ctx context.Context, in *OpenCMADRegistry, opts ...grpc.CallOption) (*OpenCMADConfig, error)
	UnRegisterOpenCMAD(ctx context.Context, in *GenericMsg, opts ...grpc.CallOption) (*GenericMsg, error)
}

type openCMDServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOpenCMDServiceClient(cc grpc.ClientConnInterface) OpenCMDServiceClient {
	return &openCMDServiceClient{cc}
}

func (c *openCMDServiceClient) RegisterOpenCMAD(ctx context.Context, in *OpenCMADRegistry, opts ...grpc.CallOption) (*OpenCMADConfig, error) {
	out := new(OpenCMADConfig)
	err := c.cc.Invoke(ctx, "/pb.OpenCMDService/RegisterOpenCMAD", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *openCMDServiceClient) UnRegisterOpenCMAD(ctx context.Context, in *GenericMsg, opts ...grpc.CallOption) (*GenericMsg, error) {
	out := new(GenericMsg)
	err := c.cc.Invoke(ctx, "/pb.OpenCMDService/UnRegisterOpenCMAD", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OpenCMDServiceServer is the server API for OpenCMDService service.
type OpenCMDServiceServer interface {
	RegisterOpenCMAD(context.Context, *OpenCMADRegistry) (*OpenCMADConfig, error)
	UnRegisterOpenCMAD(context.Context, *GenericMsg) (*GenericMsg, error)
}

// UnimplementedOpenCMDServiceServer can be embedded to have forward compatible implementations.
type UnimplementedOpenCMDServiceServer struct {
}

func (*UnimplementedOpenCMDServiceServer) RegisterOpenCMAD(context.Context, *OpenCMADRegistry) (*OpenCMADConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterOpenCMAD not implemented")
}
func (*UnimplementedOpenCMDServiceServer) UnRegisterOpenCMAD(context.Context, *GenericMsg) (*GenericMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnRegisterOpenCMAD not implemented")
}

func RegisterOpenCMDServiceServer(s *grpc.Server, srv OpenCMDServiceServer) {
	s.RegisterService(&_OpenCMDService_serviceDesc, srv)
}

func _OpenCMDService_RegisterOpenCMAD_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenCMADRegistry)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpenCMDServiceServer).RegisterOpenCMAD(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.OpenCMDService/RegisterOpenCMAD",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpenCMDServiceServer).RegisterOpenCMAD(ctx, req.(*OpenCMADRegistry))
	}
	return interceptor(ctx, in, info, handler)
}

func _OpenCMDService_UnRegisterOpenCMAD_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenericMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpenCMDServiceServer).UnRegisterOpenCMAD(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.OpenCMDService/UnRegisterOpenCMAD",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpenCMDServiceServer).UnRegisterOpenCMAD(ctx, req.(*GenericMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _OpenCMDService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.OpenCMDService",
	HandlerType: (*OpenCMDServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterOpenCMAD",
			Handler:    _OpenCMDService_RegisterOpenCMAD_Handler,
		},
		{
			MethodName: "UnRegisterOpenCMAD",
			Handler:    _OpenCMDService_UnRegisterOpenCMAD_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/opencmd.proto",
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.24.4
// source: accounts.proto

package accounts

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

type GetAccountIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	AccountName string `protobuf:"bytes,2,opt,name=accountName,proto3" json:"accountName,omitempty"`
}

func (x *GetAccountIn) Reset() {
	*x = GetAccountIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccountIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountIn) ProtoMessage() {}

func (x *GetAccountIn) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountIn.ProtoReflect.Descriptor instead.
func (*GetAccountIn) Descriptor() ([]byte, []int) {
	return file_accounts_proto_rawDescGZIP(), []int{0}
}

func (x *GetAccountIn) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetAccountIn) GetAccountName() string {
	if x != nil {
		return x.AccountName
	}
	return ""
}

type GetAccountOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsActive               bool   `protobuf:"varint,1,opt,name=isActive,proto3" json:"isActive,omitempty"`
	TargetNamespace        string `protobuf:"bytes,2,opt,name=targetNamespace,proto3" json:"targetNamespace,omitempty"`
	AccountId              string `protobuf:"bytes,3,opt,name=accountId,proto3" json:"accountId,omitempty"`
	KloudliteGatewayRegion string `protobuf:"bytes,4,opt,name=kloudliteGatewayRegion,proto3" json:"kloudliteGatewayRegion,omitempty"`
}

func (x *GetAccountOut) Reset() {
	*x = GetAccountOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccountOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountOut) ProtoMessage() {}

func (x *GetAccountOut) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountOut.ProtoReflect.Descriptor instead.
func (*GetAccountOut) Descriptor() ([]byte, []int) {
	return file_accounts_proto_rawDescGZIP(), []int{1}
}

func (x *GetAccountOut) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

func (x *GetAccountOut) GetTargetNamespace() string {
	if x != nil {
		return x.TargetNamespace
	}
	return ""
}

func (x *GetAccountOut) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *GetAccountOut) GetKloudliteGatewayRegion() string {
	if x != nil {
		return x.KloudliteGatewayRegion
	}
	return ""
}

var File_accounts_proto protoreflect.FileDescriptor

var file_accounts_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x48, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x0d, 0x47,
	0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4f, 0x75, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x69, 0x73, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x69, 0x73, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x36, 0x0a, 0x16, 0x6b, 0x6c, 0x6f, 0x75, 0x64, 0x6c, 0x69, 0x74, 0x65, 0x47, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x16, 0x6b, 0x6c, 0x6f, 0x75, 0x64, 0x6c, 0x69, 0x74, 0x65, 0x47, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x32, 0x37, 0x0a, 0x08, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x0d, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x6e, 0x1a, 0x0e, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4f, 0x75,
	0x74, 0x42, 0x1b, 0x5a, 0x19, 0x6b, 0x6c, 0x6f, 0x75, 0x64, 0x6c, 0x69, 0x74, 0x65, 0x2e, 0x69,
	0x6f, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_accounts_proto_rawDescOnce sync.Once
	file_accounts_proto_rawDescData = file_accounts_proto_rawDesc
)

func file_accounts_proto_rawDescGZIP() []byte {
	file_accounts_proto_rawDescOnce.Do(func() {
		file_accounts_proto_rawDescData = protoimpl.X.CompressGZIP(file_accounts_proto_rawDescData)
	})
	return file_accounts_proto_rawDescData
}

var file_accounts_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_accounts_proto_goTypes = []interface{}{
	(*GetAccountIn)(nil),  // 0: GetAccountIn
	(*GetAccountOut)(nil), // 1: GetAccountOut
}
var file_accounts_proto_depIdxs = []int32{
	0, // 0: Accounts.GetAccount:input_type -> GetAccountIn
	1, // 1: Accounts.GetAccount:output_type -> GetAccountOut
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_accounts_proto_init() }
func file_accounts_proto_init() {
	if File_accounts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_accounts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccountIn); i {
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
		file_accounts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccountOut); i {
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
			RawDescriptor: file_accounts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_accounts_proto_goTypes,
		DependencyIndexes: file_accounts_proto_depIdxs,
		MessageInfos:      file_accounts_proto_msgTypes,
	}.Build()
	File_accounts_proto = out.File
	file_accounts_proto_rawDesc = nil
	file_accounts_proto_goTypes = nil
	file_accounts_proto_depIdxs = nil
}

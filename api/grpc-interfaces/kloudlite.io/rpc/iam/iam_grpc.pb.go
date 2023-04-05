// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: iam.proto

package iam

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

// IAMClient is the client API for IAM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IAMClient interface {
	// Query
	Ping(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	Can(ctx context.Context, in *CanIn, opts ...grpc.CallOption) (*CanOut, error)
	ListUserMemberships(ctx context.Context, in *UserMembershipsIn, opts ...grpc.CallOption) (*ListMembershipsOut, error)
	GetMembership(ctx context.Context, in *GetMembershipIn, opts ...grpc.CallOption) (*GetMembershipOut, error)
	ListResourceMemberships(ctx context.Context, in *ResourceMembershipsIn, opts ...grpc.CallOption) (*ListMembershipsOut, error)
	ListMembershipsByResource(ctx context.Context, in *MembershipsByResourceIn, opts ...grpc.CallOption) (*ListMembershipsOut, error)
	ListMembershipsForUser(ctx context.Context, in *MembershipsForUserIn, opts ...grpc.CallOption) (*ListMembershipsOut, error)
	// Mutation
	AddMembership(ctx context.Context, in *AddMembershipIn, opts ...grpc.CallOption) (*AddMembershipOut, error)
	InviteMembership(ctx context.Context, in *AddMembershipIn, opts ...grpc.CallOption) (*AddMembershipOut, error)
	ConfirmMembership(ctx context.Context, in *ConfirmMembershipIn, opts ...grpc.CallOption) (*ConfirmMembershipOut, error)
	RemoveMembership(ctx context.Context, in *RemoveMembershipIn, opts ...grpc.CallOption) (*RemoveMembershipOut, error)
	RemoveResource(ctx context.Context, in *RemoveResourceIn, opts ...grpc.CallOption) (*RemoveResourceOut, error)
}

type iAMClient struct {
	cc grpc.ClientConnInterface
}

func NewIAMClient(cc grpc.ClientConnInterface) IAMClient {
	return &iAMClient{cc}
}

func (c *iAMClient) Ping(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/IAM/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) Can(ctx context.Context, in *CanIn, opts ...grpc.CallOption) (*CanOut, error) {
	out := new(CanOut)
	err := c.cc.Invoke(ctx, "/IAM/Can", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListUserMemberships(ctx context.Context, in *UserMembershipsIn, opts ...grpc.CallOption) (*ListMembershipsOut, error) {
	out := new(ListMembershipsOut)
	err := c.cc.Invoke(ctx, "/IAM/ListUserMemberships", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) GetMembership(ctx context.Context, in *GetMembershipIn, opts ...grpc.CallOption) (*GetMembershipOut, error) {
	out := new(GetMembershipOut)
	err := c.cc.Invoke(ctx, "/IAM/GetMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListResourceMemberships(ctx context.Context, in *ResourceMembershipsIn, opts ...grpc.CallOption) (*ListMembershipsOut, error) {
	out := new(ListMembershipsOut)
	err := c.cc.Invoke(ctx, "/IAM/ListResourceMemberships", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListMembershipsByResource(ctx context.Context, in *MembershipsByResourceIn, opts ...grpc.CallOption) (*ListMembershipsOut, error) {
	out := new(ListMembershipsOut)
	err := c.cc.Invoke(ctx, "/IAM/ListMembershipsByResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListMembershipsForUser(ctx context.Context, in *MembershipsForUserIn, opts ...grpc.CallOption) (*ListMembershipsOut, error) {
	out := new(ListMembershipsOut)
	err := c.cc.Invoke(ctx, "/IAM/ListMembershipsForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) AddMembership(ctx context.Context, in *AddMembershipIn, opts ...grpc.CallOption) (*AddMembershipOut, error) {
	out := new(AddMembershipOut)
	err := c.cc.Invoke(ctx, "/IAM/AddMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) InviteMembership(ctx context.Context, in *AddMembershipIn, opts ...grpc.CallOption) (*AddMembershipOut, error) {
	out := new(AddMembershipOut)
	err := c.cc.Invoke(ctx, "/IAM/InviteMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ConfirmMembership(ctx context.Context, in *ConfirmMembershipIn, opts ...grpc.CallOption) (*ConfirmMembershipOut, error) {
	out := new(ConfirmMembershipOut)
	err := c.cc.Invoke(ctx, "/IAM/ConfirmMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) RemoveMembership(ctx context.Context, in *RemoveMembershipIn, opts ...grpc.CallOption) (*RemoveMembershipOut, error) {
	out := new(RemoveMembershipOut)
	err := c.cc.Invoke(ctx, "/IAM/RemoveMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) RemoveResource(ctx context.Context, in *RemoveResourceIn, opts ...grpc.CallOption) (*RemoveResourceOut, error) {
	out := new(RemoveResourceOut)
	err := c.cc.Invoke(ctx, "/IAM/RemoveResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IAMServer is the server API for IAM service.
// All implementations must embed UnimplementedIAMServer
// for forward compatibility
type IAMServer interface {
	// Query
	Ping(context.Context, *Message) (*Message, error)
	Can(context.Context, *CanIn) (*CanOut, error)
	ListUserMemberships(context.Context, *UserMembershipsIn) (*ListMembershipsOut, error)
	GetMembership(context.Context, *GetMembershipIn) (*GetMembershipOut, error)
	ListResourceMemberships(context.Context, *ResourceMembershipsIn) (*ListMembershipsOut, error)
	ListMembershipsByResource(context.Context, *MembershipsByResourceIn) (*ListMembershipsOut, error)
	ListMembershipsForUser(context.Context, *MembershipsForUserIn) (*ListMembershipsOut, error)
	// Mutation
	AddMembership(context.Context, *AddMembershipIn) (*AddMembershipOut, error)
	InviteMembership(context.Context, *AddMembershipIn) (*AddMembershipOut, error)
	ConfirmMembership(context.Context, *ConfirmMembershipIn) (*ConfirmMembershipOut, error)
	RemoveMembership(context.Context, *RemoveMembershipIn) (*RemoveMembershipOut, error)
	RemoveResource(context.Context, *RemoveResourceIn) (*RemoveResourceOut, error)
	mustEmbedUnimplementedIAMServer()
}

// UnimplementedIAMServer must be embedded to have forward compatible implementations.
type UnimplementedIAMServer struct {
}

func (UnimplementedIAMServer) Ping(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedIAMServer) Can(context.Context, *CanIn) (*CanOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Can not implemented")
}
func (UnimplementedIAMServer) ListUserMemberships(context.Context, *UserMembershipsIn) (*ListMembershipsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserMemberships not implemented")
}
func (UnimplementedIAMServer) GetMembership(context.Context, *GetMembershipIn) (*GetMembershipOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMembership not implemented")
}
func (UnimplementedIAMServer) ListResourceMemberships(context.Context, *ResourceMembershipsIn) (*ListMembershipsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListResourceMemberships not implemented")
}
func (UnimplementedIAMServer) ListMembershipsByResource(context.Context, *MembershipsByResourceIn) (*ListMembershipsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMembershipsByResource not implemented")
}
func (UnimplementedIAMServer) ListMembershipsForUser(context.Context, *MembershipsForUserIn) (*ListMembershipsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMembershipsForUser not implemented")
}
func (UnimplementedIAMServer) AddMembership(context.Context, *AddMembershipIn) (*AddMembershipOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMembership not implemented")
}
func (UnimplementedIAMServer) InviteMembership(context.Context, *AddMembershipIn) (*AddMembershipOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteMembership not implemented")
}
func (UnimplementedIAMServer) ConfirmMembership(context.Context, *ConfirmMembershipIn) (*ConfirmMembershipOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmMembership not implemented")
}
func (UnimplementedIAMServer) RemoveMembership(context.Context, *RemoveMembershipIn) (*RemoveMembershipOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMembership not implemented")
}
func (UnimplementedIAMServer) RemoveResource(context.Context, *RemoveResourceIn) (*RemoveResourceOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveResource not implemented")
}
func (UnimplementedIAMServer) mustEmbedUnimplementedIAMServer() {}

// UnsafeIAMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IAMServer will
// result in compilation errors.
type UnsafeIAMServer interface {
	mustEmbedUnimplementedIAMServer()
}

func RegisterIAMServer(s grpc.ServiceRegistrar, srv IAMServer) {
	s.RegisterService(&IAM_ServiceDesc, srv)
}

func _IAM_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).Ping(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_Can_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CanIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).Can(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/Can",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).Can(ctx, req.(*CanIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListUserMemberships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMembershipsIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).ListUserMemberships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/ListUserMemberships",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).ListUserMemberships(ctx, req.(*UserMembershipsIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_GetMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMembershipIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).GetMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/GetMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).GetMembership(ctx, req.(*GetMembershipIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListResourceMemberships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceMembershipsIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).ListResourceMemberships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/ListResourceMemberships",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).ListResourceMemberships(ctx, req.(*ResourceMembershipsIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListMembershipsByResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MembershipsByResourceIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).ListMembershipsByResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/ListMembershipsByResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).ListMembershipsByResource(ctx, req.(*MembershipsByResourceIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListMembershipsForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MembershipsForUserIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).ListMembershipsForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/ListMembershipsForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).ListMembershipsForUser(ctx, req.(*MembershipsForUserIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_AddMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMembershipIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).AddMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/AddMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).AddMembership(ctx, req.(*AddMembershipIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_InviteMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMembershipIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).InviteMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/InviteMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).InviteMembership(ctx, req.(*AddMembershipIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ConfirmMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmMembershipIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).ConfirmMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/ConfirmMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).ConfirmMembership(ctx, req.(*ConfirmMembershipIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_RemoveMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveMembershipIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).RemoveMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/RemoveMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).RemoveMembership(ctx, req.(*RemoveMembershipIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_RemoveResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveResourceIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).RemoveResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IAM/RemoveResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).RemoveResource(ctx, req.(*RemoveResourceIn))
	}
	return interceptor(ctx, in, info, handler)
}

// IAM_ServiceDesc is the grpc.ServiceDesc for IAM service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IAM_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IAM",
	HandlerType: (*IAMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _IAM_Ping_Handler,
		},
		{
			MethodName: "Can",
			Handler:    _IAM_Can_Handler,
		},
		{
			MethodName: "ListUserMemberships",
			Handler:    _IAM_ListUserMemberships_Handler,
		},
		{
			MethodName: "GetMembership",
			Handler:    _IAM_GetMembership_Handler,
		},
		{
			MethodName: "ListResourceMemberships",
			Handler:    _IAM_ListResourceMemberships_Handler,
		},
		{
			MethodName: "ListMembershipsByResource",
			Handler:    _IAM_ListMembershipsByResource_Handler,
		},
		{
			MethodName: "ListMembershipsForUser",
			Handler:    _IAM_ListMembershipsForUser_Handler,
		},
		{
			MethodName: "AddMembership",
			Handler:    _IAM_AddMembership_Handler,
		},
		{
			MethodName: "InviteMembership",
			Handler:    _IAM_InviteMembership_Handler,
		},
		{
			MethodName: "ConfirmMembership",
			Handler:    _IAM_ConfirmMembership_Handler,
		},
		{
			MethodName: "RemoveMembership",
			Handler:    _IAM_RemoveMembership_Handler,
		},
		{
			MethodName: "RemoveResource",
			Handler:    _IAM_RemoveResource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iam.proto",
}

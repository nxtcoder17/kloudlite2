// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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
	Can(ctx context.Context, in *InCan, opts ...grpc.CallOption) (*OutCan, error)
	ListUserMemberships(ctx context.Context, in *InUserMemberships, opts ...grpc.CallOption) (*OutListMemberships, error)
	GetMembership(ctx context.Context, in *InGetMembership, opts ...grpc.CallOption) (*OutGetMembership, error)
	ListResourceMemberships(ctx context.Context, in *InResourceMemberships, opts ...grpc.CallOption) (*OutListMemberships, error)
	// Mutation
	AddMembership(ctx context.Context, in *InAddMembership, opts ...grpc.CallOption) (*OutAddMembership, error)
	InviteMembership(ctx context.Context, in *InAddMembership, opts ...grpc.CallOption) (*OutAddMembership, error)
	RemoveMembership(ctx context.Context, in *InRemoveMembership, opts ...grpc.CallOption) (*OutRemoveMembership, error)
	RemoveResource(ctx context.Context, in *InRemoveResource, opts ...grpc.CallOption) (*OutRemoveResource, error)
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

func (c *iAMClient) Can(ctx context.Context, in *InCan, opts ...grpc.CallOption) (*OutCan, error) {
	out := new(OutCan)
	err := c.cc.Invoke(ctx, "/IAM/Can", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListUserMemberships(ctx context.Context, in *InUserMemberships, opts ...grpc.CallOption) (*OutListMemberships, error) {
	out := new(OutListMemberships)
	err := c.cc.Invoke(ctx, "/IAM/ListUserMemberships", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) GetMembership(ctx context.Context, in *InGetMembership, opts ...grpc.CallOption) (*OutGetMembership, error) {
	out := new(OutGetMembership)
	err := c.cc.Invoke(ctx, "/IAM/GetMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListResourceMemberships(ctx context.Context, in *InResourceMemberships, opts ...grpc.CallOption) (*OutListMemberships, error) {
	out := new(OutListMemberships)
	err := c.cc.Invoke(ctx, "/IAM/ListResourceMemberships", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) AddMembership(ctx context.Context, in *InAddMembership, opts ...grpc.CallOption) (*OutAddMembership, error) {
	out := new(OutAddMembership)
	err := c.cc.Invoke(ctx, "/IAM/AddMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) InviteMembership(ctx context.Context, in *InAddMembership, opts ...grpc.CallOption) (*OutAddMembership, error) {
	out := new(OutAddMembership)
	err := c.cc.Invoke(ctx, "/IAM/InviteMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) RemoveMembership(ctx context.Context, in *InRemoveMembership, opts ...grpc.CallOption) (*OutRemoveMembership, error) {
	out := new(OutRemoveMembership)
	err := c.cc.Invoke(ctx, "/IAM/RemoveMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) RemoveResource(ctx context.Context, in *InRemoveResource, opts ...grpc.CallOption) (*OutRemoveResource, error) {
	out := new(OutRemoveResource)
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
	Can(context.Context, *InCan) (*OutCan, error)
	ListUserMemberships(context.Context, *InUserMemberships) (*OutListMemberships, error)
	GetMembership(context.Context, *InGetMembership) (*OutGetMembership, error)
	ListResourceMemberships(context.Context, *InResourceMemberships) (*OutListMemberships, error)
	// Mutation
	AddMembership(context.Context, *InAddMembership) (*OutAddMembership, error)
	InviteMembership(context.Context, *InAddMembership) (*OutAddMembership, error)
	RemoveMembership(context.Context, *InRemoveMembership) (*OutRemoveMembership, error)
	RemoveResource(context.Context, *InRemoveResource) (*OutRemoveResource, error)
	mustEmbedUnimplementedIAMServer()
}

// UnimplementedIAMServer must be embedded to have forward compatible implementations.
type UnimplementedIAMServer struct {
}

func (UnimplementedIAMServer) Ping(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedIAMServer) Can(context.Context, *InCan) (*OutCan, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Can not implemented")
}
func (UnimplementedIAMServer) ListUserMemberships(context.Context, *InUserMemberships) (*OutListMemberships, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserMemberships not implemented")
}
func (UnimplementedIAMServer) GetMembership(context.Context, *InGetMembership) (*OutGetMembership, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMembership not implemented")
}
func (UnimplementedIAMServer) ListResourceMemberships(context.Context, *InResourceMemberships) (*OutListMemberships, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListResourceMemberships not implemented")
}
func (UnimplementedIAMServer) AddMembership(context.Context, *InAddMembership) (*OutAddMembership, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMembership not implemented")
}
func (UnimplementedIAMServer) InviteMembership(context.Context, *InAddMembership) (*OutAddMembership, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteMembership not implemented")
}
func (UnimplementedIAMServer) RemoveMembership(context.Context, *InRemoveMembership) (*OutRemoveMembership, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMembership not implemented")
}
func (UnimplementedIAMServer) RemoveResource(context.Context, *InRemoveResource) (*OutRemoveResource, error) {
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
	in := new(InCan)
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
		return srv.(IAMServer).Can(ctx, req.(*InCan))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListUserMemberships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InUserMemberships)
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
		return srv.(IAMServer).ListUserMemberships(ctx, req.(*InUserMemberships))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_GetMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InGetMembership)
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
		return srv.(IAMServer).GetMembership(ctx, req.(*InGetMembership))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListResourceMemberships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InResourceMemberships)
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
		return srv.(IAMServer).ListResourceMemberships(ctx, req.(*InResourceMemberships))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_AddMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InAddMembership)
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
		return srv.(IAMServer).AddMembership(ctx, req.(*InAddMembership))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_InviteMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InAddMembership)
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
		return srv.(IAMServer).InviteMembership(ctx, req.(*InAddMembership))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_RemoveMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InRemoveMembership)
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
		return srv.(IAMServer).RemoveMembership(ctx, req.(*InRemoveMembership))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_RemoveResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InRemoveResource)
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
		return srv.(IAMServer).RemoveResource(ctx, req.(*InRemoveResource))
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
			MethodName: "AddMembership",
			Handler:    _IAM_AddMembership_Handler,
		},
		{
			MethodName: "InviteMembership",
			Handler:    _IAM_InviteMembership_Handler,
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

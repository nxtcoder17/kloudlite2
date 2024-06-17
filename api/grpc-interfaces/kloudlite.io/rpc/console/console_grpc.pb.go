// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: console.proto

package console

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

const (
	Console_GetProjectName_FullMethodName = "/Console/GetProjectName"
	Console_GetApp_FullMethodName         = "/Console/GetApp"
	Console_GetManagedSvc_FullMethodName  = "/Console/GetManagedSvc"
	Console_SetupAccount_FullMethodName   = "/Console/SetupAccount"
)

// ConsoleClient is the client API for Console service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsoleClient interface {
	GetProjectName(ctx context.Context, in *ProjectIn, opts ...grpc.CallOption) (*ProjectOut, error)
	GetApp(ctx context.Context, in *AppIn, opts ...grpc.CallOption) (*AppOut, error)
	GetManagedSvc(ctx context.Context, in *MSvcIn, opts ...grpc.CallOption) (*MSvcOut, error)
	SetupAccount(ctx context.Context, in *AccountSetupIn, opts ...grpc.CallOption) (*AccountSetupVoid, error)
}

type consoleClient struct {
	cc grpc.ClientConnInterface
}

func NewConsoleClient(cc grpc.ClientConnInterface) ConsoleClient {
	return &consoleClient{cc}
}

func (c *consoleClient) GetProjectName(ctx context.Context, in *ProjectIn, opts ...grpc.CallOption) (*ProjectOut, error) {
	out := new(ProjectOut)
	err := c.cc.Invoke(ctx, Console_GetProjectName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consoleClient) GetApp(ctx context.Context, in *AppIn, opts ...grpc.CallOption) (*AppOut, error) {
	out := new(AppOut)
	err := c.cc.Invoke(ctx, Console_GetApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consoleClient) GetManagedSvc(ctx context.Context, in *MSvcIn, opts ...grpc.CallOption) (*MSvcOut, error) {
	out := new(MSvcOut)
	err := c.cc.Invoke(ctx, Console_GetManagedSvc_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consoleClient) SetupAccount(ctx context.Context, in *AccountSetupIn, opts ...grpc.CallOption) (*AccountSetupVoid, error) {
	out := new(AccountSetupVoid)
	err := c.cc.Invoke(ctx, Console_SetupAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsoleServer is the server API for Console service.
// All implementations must embed UnimplementedConsoleServer
// for forward compatibility
type ConsoleServer interface {
	GetProjectName(context.Context, *ProjectIn) (*ProjectOut, error)
	GetApp(context.Context, *AppIn) (*AppOut, error)
	GetManagedSvc(context.Context, *MSvcIn) (*MSvcOut, error)
	SetupAccount(context.Context, *AccountSetupIn) (*AccountSetupVoid, error)
	mustEmbedUnimplementedConsoleServer()
}

// UnimplementedConsoleServer must be embedded to have forward compatible implementations.
type UnimplementedConsoleServer struct {
}

func (UnimplementedConsoleServer) GetProjectName(context.Context, *ProjectIn) (*ProjectOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectName not implemented")
}
func (UnimplementedConsoleServer) GetApp(context.Context, *AppIn) (*AppOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApp not implemented")
}
func (UnimplementedConsoleServer) GetManagedSvc(context.Context, *MSvcIn) (*MSvcOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManagedSvc not implemented")
}
func (UnimplementedConsoleServer) SetupAccount(context.Context, *AccountSetupIn) (*AccountSetupVoid, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetupAccount not implemented")
}
func (UnimplementedConsoleServer) mustEmbedUnimplementedConsoleServer() {}

// UnsafeConsoleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsoleServer will
// result in compilation errors.
type UnsafeConsoleServer interface {
	mustEmbedUnimplementedConsoleServer()
}

func RegisterConsoleServer(s grpc.ServiceRegistrar, srv ConsoleServer) {
	s.RegisterService(&Console_ServiceDesc, srv)
}

func _Console_GetProjectName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsoleServer).GetProjectName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Console_GetProjectName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsoleServer).GetProjectName(ctx, req.(*ProjectIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Console_GetApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsoleServer).GetApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Console_GetApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsoleServer).GetApp(ctx, req.(*AppIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Console_GetManagedSvc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MSvcIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsoleServer).GetManagedSvc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Console_GetManagedSvc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsoleServer).GetManagedSvc(ctx, req.(*MSvcIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Console_SetupAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountSetupIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsoleServer).SetupAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Console_SetupAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsoleServer).SetupAccount(ctx, req.(*AccountSetupIn))
	}
	return interceptor(ctx, in, info, handler)
}

// Console_ServiceDesc is the grpc.ServiceDesc for Console service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Console_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Console",
	HandlerType: (*ConsoleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProjectName",
			Handler:    _Console_GetProjectName_Handler,
		},
		{
			MethodName: "GetApp",
			Handler:    _Console_GetApp_Handler,
		},
		{
			MethodName: "GetManagedSvc",
			Handler:    _Console_GetManagedSvc_Handler,
		},
		{
			MethodName: "SetupAccount",
			Handler:    _Console_SetupAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "console.proto",
}

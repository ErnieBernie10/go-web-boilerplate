// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/framer.proto

package proto

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
	AppUserService_Register_FullMethodName = "/public.AppUserService/Register"
	AppUserService_Login_FullMethodName    = "/public.AppUserService/Login"
)

// AppUserServiceClient is the client API for AppUserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// gRPC Services
type AppUserServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*AppUser, error)
}

type appUserServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppUserServiceClient(cc grpc.ClientConnInterface) AppUserServiceClient {
	return &appUserServiceClient{cc}
}

func (c *appUserServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, AppUserService_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appUserServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*AppUser, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AppUser)
	err := c.cc.Invoke(ctx, AppUserService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppUserServiceServer is the server API for AppUserService service.
// All implementations must embed UnimplementedAppUserServiceServer
// for forward compatibility.
//
// gRPC Services
type AppUserServiceServer interface {
	Register(context.Context, *RegisterRequest) (*EmptyResponse, error)
	Login(context.Context, *LoginRequest) (*AppUser, error)
	mustEmbedUnimplementedAppUserServiceServer()
}

// UnimplementedAppUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAppUserServiceServer struct{}

func (UnimplementedAppUserServiceServer) Register(context.Context, *RegisterRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAppUserServiceServer) Login(context.Context, *LoginRequest) (*AppUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAppUserServiceServer) mustEmbedUnimplementedAppUserServiceServer() {}
func (UnimplementedAppUserServiceServer) testEmbeddedByValue()                        {}

// UnsafeAppUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppUserServiceServer will
// result in compilation errors.
type UnsafeAppUserServiceServer interface {
	mustEmbedUnimplementedAppUserServiceServer()
}

func RegisterAppUserServiceServer(s grpc.ServiceRegistrar, srv AppUserServiceServer) {
	// If the following call pancis, it indicates UnimplementedAppUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AppUserService_ServiceDesc, srv)
}

func _AppUserService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppUserServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AppUserService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppUserServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppUserService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppUserServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AppUserService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppUserServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AppUserService_ServiceDesc is the grpc.ServiceDesc for AppUserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppUserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "public.AppUserService",
	HandlerType: (*AppUserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _AppUserService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AppUserService_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/framer.proto",
}

const (
	FrameService_CreateFrame_FullMethodName = "/public.FrameService/CreateFrame"
	FrameService_GetFrame_FullMethodName    = "/public.FrameService/GetFrame"
	FrameService_UpdateFrame_FullMethodName = "/public.FrameService/UpdateFrame"
	FrameService_DeleteFrame_FullMethodName = "/public.FrameService/DeleteFrame"
	FrameService_ListFrames_FullMethodName  = "/public.FrameService/ListFrames"
)

// FrameServiceClient is the client API for FrameService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrameServiceClient interface {
	CreateFrame(ctx context.Context, in *CreateFrameRequest, opts ...grpc.CallOption) (*Frame, error)
	GetFrame(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*Frame, error)
	UpdateFrame(ctx context.Context, in *UpdateFrameRequest, opts ...grpc.CallOption) (*Frame, error)
	DeleteFrame(ctx context.Context, in *DeleteByIdRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	ListFrames(ctx context.Context, in *EmptyResponse, opts ...grpc.CallOption) (*ListFramesResponse, error)
}

type frameServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFrameServiceClient(cc grpc.ClientConnInterface) FrameServiceClient {
	return &frameServiceClient{cc}
}

func (c *frameServiceClient) CreateFrame(ctx context.Context, in *CreateFrameRequest, opts ...grpc.CallOption) (*Frame, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Frame)
	err := c.cc.Invoke(ctx, FrameService_CreateFrame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frameServiceClient) GetFrame(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*Frame, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Frame)
	err := c.cc.Invoke(ctx, FrameService_GetFrame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frameServiceClient) UpdateFrame(ctx context.Context, in *UpdateFrameRequest, opts ...grpc.CallOption) (*Frame, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Frame)
	err := c.cc.Invoke(ctx, FrameService_UpdateFrame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frameServiceClient) DeleteFrame(ctx context.Context, in *DeleteByIdRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, FrameService_DeleteFrame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frameServiceClient) ListFrames(ctx context.Context, in *EmptyResponse, opts ...grpc.CallOption) (*ListFramesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListFramesResponse)
	err := c.cc.Invoke(ctx, FrameService_ListFrames_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FrameServiceServer is the server API for FrameService service.
// All implementations must embed UnimplementedFrameServiceServer
// for forward compatibility.
type FrameServiceServer interface {
	CreateFrame(context.Context, *CreateFrameRequest) (*Frame, error)
	GetFrame(context.Context, *GetByIdRequest) (*Frame, error)
	UpdateFrame(context.Context, *UpdateFrameRequest) (*Frame, error)
	DeleteFrame(context.Context, *DeleteByIdRequest) (*EmptyResponse, error)
	ListFrames(context.Context, *EmptyResponse) (*ListFramesResponse, error)
	mustEmbedUnimplementedFrameServiceServer()
}

// UnimplementedFrameServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFrameServiceServer struct{}

func (UnimplementedFrameServiceServer) CreateFrame(context.Context, *CreateFrameRequest) (*Frame, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFrame not implemented")
}
func (UnimplementedFrameServiceServer) GetFrame(context.Context, *GetByIdRequest) (*Frame, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFrame not implemented")
}
func (UnimplementedFrameServiceServer) UpdateFrame(context.Context, *UpdateFrameRequest) (*Frame, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFrame not implemented")
}
func (UnimplementedFrameServiceServer) DeleteFrame(context.Context, *DeleteByIdRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFrame not implemented")
}
func (UnimplementedFrameServiceServer) ListFrames(context.Context, *EmptyResponse) (*ListFramesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFrames not implemented")
}
func (UnimplementedFrameServiceServer) mustEmbedUnimplementedFrameServiceServer() {}
func (UnimplementedFrameServiceServer) testEmbeddedByValue()                      {}

// UnsafeFrameServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrameServiceServer will
// result in compilation errors.
type UnsafeFrameServiceServer interface {
	mustEmbedUnimplementedFrameServiceServer()
}

func RegisterFrameServiceServer(s grpc.ServiceRegistrar, srv FrameServiceServer) {
	// If the following call pancis, it indicates UnimplementedFrameServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FrameService_ServiceDesc, srv)
}

func _FrameService_CreateFrame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFrameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrameServiceServer).CreateFrame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FrameService_CreateFrame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrameServiceServer).CreateFrame(ctx, req.(*CreateFrameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FrameService_GetFrame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrameServiceServer).GetFrame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FrameService_GetFrame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrameServiceServer).GetFrame(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FrameService_UpdateFrame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFrameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrameServiceServer).UpdateFrame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FrameService_UpdateFrame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrameServiceServer).UpdateFrame(ctx, req.(*UpdateFrameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FrameService_DeleteFrame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrameServiceServer).DeleteFrame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FrameService_DeleteFrame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrameServiceServer).DeleteFrame(ctx, req.(*DeleteByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FrameService_ListFrames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrameServiceServer).ListFrames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FrameService_ListFrames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrameServiceServer).ListFrames(ctx, req.(*EmptyResponse))
	}
	return interceptor(ctx, in, info, handler)
}

// FrameService_ServiceDesc is the grpc.ServiceDesc for FrameService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FrameService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "public.FrameService",
	HandlerType: (*FrameServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFrame",
			Handler:    _FrameService_CreateFrame_Handler,
		},
		{
			MethodName: "GetFrame",
			Handler:    _FrameService_GetFrame_Handler,
		},
		{
			MethodName: "UpdateFrame",
			Handler:    _FrameService_UpdateFrame_Handler,
		},
		{
			MethodName: "DeleteFrame",
			Handler:    _FrameService_DeleteFrame_Handler,
		},
		{
			MethodName: "ListFrames",
			Handler:    _FrameService_ListFrames_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/framer.proto",
}

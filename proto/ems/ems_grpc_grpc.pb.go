// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.6
// source: proto/ems/ems_grpc.proto

package ems

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
	GRPCConfigOper_GetConfig_FullMethodName            = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/GetConfig"
	GRPCConfigOper_MergeConfig_FullMethodName          = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/MergeConfig"
	GRPCConfigOper_DeleteConfig_FullMethodName         = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/DeleteConfig"
	GRPCConfigOper_RemoveConfig_FullMethodName         = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/RemoveConfig"
	GRPCConfigOper_ReplaceConfig_FullMethodName        = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/ReplaceConfig"
	GRPCConfigOper_CliConfig_FullMethodName            = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/CliConfig"
	GRPCConfigOper_CommitReplace_FullMethodName        = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/CommitReplace"
	GRPCConfigOper_CommitConfig_FullMethodName         = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/CommitConfig"
	GRPCConfigOper_ConfigDiscardChanges_FullMethodName = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/ConfigDiscardChanges"
	GRPCConfigOper_GetOper_FullMethodName              = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/GetOper"
	GRPCConfigOper_CreateSubs_FullMethodName           = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/CreateSubs"
	GRPCConfigOper_GetProtoFile_FullMethodName         = "/IOSXRExtensibleManagabilityService.gRPCConfigOper/GetProtoFile"
)

// GRPCConfigOperClient is the client API for GRPCConfigOper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GRPCConfigOperClient interface {
	GetConfig(ctx context.Context, in *ConfigGetArgs, opts ...grpc.CallOption) (GRPCConfigOper_GetConfigClient, error)
	MergeConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error)
	DeleteConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error)
	RemoveConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error)
	ReplaceConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error)
	CliConfig(ctx context.Context, in *CliConfigArgs, opts ...grpc.CallOption) (*CliConfigReply, error)
	CommitReplace(ctx context.Context, in *CommitReplaceArgs, opts ...grpc.CallOption) (*CommitReplaceReply, error)
	// Do we need implicit or explicit commit
	CommitConfig(ctx context.Context, in *CommitArgs, opts ...grpc.CallOption) (*CommitReply, error)
	ConfigDiscardChanges(ctx context.Context, in *DiscardChangesArgs, opts ...grpc.CallOption) (*DiscardChangesReply, error)
	// Get only returns oper data
	GetOper(ctx context.Context, in *GetOperArgs, opts ...grpc.CallOption) (GRPCConfigOper_GetOperClient, error)
	// Get Telemetry Data
	CreateSubs(ctx context.Context, in *CreateSubsArgs, opts ...grpc.CallOption) (GRPCConfigOper_CreateSubsClient, error)
	// Get Proto File
	GetProtoFile(ctx context.Context, in *GetProtoFileArgs, opts ...grpc.CallOption) (GRPCConfigOper_GetProtoFileClient, error)
}

type gRPCConfigOperClient struct {
	cc grpc.ClientConnInterface
}

func NewGRPCConfigOperClient(cc grpc.ClientConnInterface) GRPCConfigOperClient {
	return &gRPCConfigOperClient{cc}
}

func (c *gRPCConfigOperClient) GetConfig(ctx context.Context, in *ConfigGetArgs, opts ...grpc.CallOption) (GRPCConfigOper_GetConfigClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCConfigOper_ServiceDesc.Streams[0], GRPCConfigOper_GetConfig_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCConfigOperGetConfigClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCConfigOper_GetConfigClient interface {
	Recv() (*ConfigGetReply, error)
	grpc.ClientStream
}

type gRPCConfigOperGetConfigClient struct {
	grpc.ClientStream
}

func (x *gRPCConfigOperGetConfigClient) Recv() (*ConfigGetReply, error) {
	m := new(ConfigGetReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gRPCConfigOperClient) MergeConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error) {
	out := new(ConfigReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_MergeConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) DeleteConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error) {
	out := new(ConfigReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_DeleteConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) RemoveConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error) {
	out := new(ConfigReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_RemoveConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) ReplaceConfig(ctx context.Context, in *ConfigArgs, opts ...grpc.CallOption) (*ConfigReply, error) {
	out := new(ConfigReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_ReplaceConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) CliConfig(ctx context.Context, in *CliConfigArgs, opts ...grpc.CallOption) (*CliConfigReply, error) {
	out := new(CliConfigReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_CliConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) CommitReplace(ctx context.Context, in *CommitReplaceArgs, opts ...grpc.CallOption) (*CommitReplaceReply, error) {
	out := new(CommitReplaceReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_CommitReplace_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) CommitConfig(ctx context.Context, in *CommitArgs, opts ...grpc.CallOption) (*CommitReply, error) {
	out := new(CommitReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_CommitConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) ConfigDiscardChanges(ctx context.Context, in *DiscardChangesArgs, opts ...grpc.CallOption) (*DiscardChangesReply, error) {
	out := new(DiscardChangesReply)
	err := c.cc.Invoke(ctx, GRPCConfigOper_ConfigDiscardChanges_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCConfigOperClient) GetOper(ctx context.Context, in *GetOperArgs, opts ...grpc.CallOption) (GRPCConfigOper_GetOperClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCConfigOper_ServiceDesc.Streams[1], GRPCConfigOper_GetOper_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCConfigOperGetOperClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCConfigOper_GetOperClient interface {
	Recv() (*GetOperReply, error)
	grpc.ClientStream
}

type gRPCConfigOperGetOperClient struct {
	grpc.ClientStream
}

func (x *gRPCConfigOperGetOperClient) Recv() (*GetOperReply, error) {
	m := new(GetOperReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gRPCConfigOperClient) CreateSubs(ctx context.Context, in *CreateSubsArgs, opts ...grpc.CallOption) (GRPCConfigOper_CreateSubsClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCConfigOper_ServiceDesc.Streams[2], GRPCConfigOper_CreateSubs_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCConfigOperCreateSubsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCConfigOper_CreateSubsClient interface {
	Recv() (*CreateSubsReply, error)
	grpc.ClientStream
}

type gRPCConfigOperCreateSubsClient struct {
	grpc.ClientStream
}

func (x *gRPCConfigOperCreateSubsClient) Recv() (*CreateSubsReply, error) {
	m := new(CreateSubsReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gRPCConfigOperClient) GetProtoFile(ctx context.Context, in *GetProtoFileArgs, opts ...grpc.CallOption) (GRPCConfigOper_GetProtoFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCConfigOper_ServiceDesc.Streams[3], GRPCConfigOper_GetProtoFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCConfigOperGetProtoFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCConfigOper_GetProtoFileClient interface {
	Recv() (*GetProtoFileReply, error)
	grpc.ClientStream
}

type gRPCConfigOperGetProtoFileClient struct {
	grpc.ClientStream
}

func (x *gRPCConfigOperGetProtoFileClient) Recv() (*GetProtoFileReply, error) {
	m := new(GetProtoFileReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GRPCConfigOperServer is the server API for GRPCConfigOper service.
// All implementations must embed UnimplementedGRPCConfigOperServer
// for forward compatibility
type GRPCConfigOperServer interface {
	GetConfig(*ConfigGetArgs, GRPCConfigOper_GetConfigServer) error
	MergeConfig(context.Context, *ConfigArgs) (*ConfigReply, error)
	DeleteConfig(context.Context, *ConfigArgs) (*ConfigReply, error)
	RemoveConfig(context.Context, *ConfigArgs) (*ConfigReply, error)
	ReplaceConfig(context.Context, *ConfigArgs) (*ConfigReply, error)
	CliConfig(context.Context, *CliConfigArgs) (*CliConfigReply, error)
	CommitReplace(context.Context, *CommitReplaceArgs) (*CommitReplaceReply, error)
	// Do we need implicit or explicit commit
	CommitConfig(context.Context, *CommitArgs) (*CommitReply, error)
	ConfigDiscardChanges(context.Context, *DiscardChangesArgs) (*DiscardChangesReply, error)
	// Get only returns oper data
	GetOper(*GetOperArgs, GRPCConfigOper_GetOperServer) error
	// Get Telemetry Data
	CreateSubs(*CreateSubsArgs, GRPCConfigOper_CreateSubsServer) error
	// Get Proto File
	GetProtoFile(*GetProtoFileArgs, GRPCConfigOper_GetProtoFileServer) error
	mustEmbedUnimplementedGRPCConfigOperServer()
}

// UnimplementedGRPCConfigOperServer must be embedded to have forward compatible implementations.
type UnimplementedGRPCConfigOperServer struct {
}

func (UnimplementedGRPCConfigOperServer) GetConfig(*ConfigGetArgs, GRPCConfigOper_GetConfigServer) error {
	return status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) MergeConfig(context.Context, *ConfigArgs) (*ConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) DeleteConfig(context.Context, *ConfigArgs) (*ConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) RemoveConfig(context.Context, *ConfigArgs) (*ConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) ReplaceConfig(context.Context, *ConfigArgs) (*ConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) CliConfig(context.Context, *CliConfigArgs) (*CliConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CliConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) CommitReplace(context.Context, *CommitReplaceArgs) (*CommitReplaceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitReplace not implemented")
}
func (UnimplementedGRPCConfigOperServer) CommitConfig(context.Context, *CommitArgs) (*CommitReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitConfig not implemented")
}
func (UnimplementedGRPCConfigOperServer) ConfigDiscardChanges(context.Context, *DiscardChangesArgs) (*DiscardChangesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigDiscardChanges not implemented")
}
func (UnimplementedGRPCConfigOperServer) GetOper(*GetOperArgs, GRPCConfigOper_GetOperServer) error {
	return status.Errorf(codes.Unimplemented, "method GetOper not implemented")
}
func (UnimplementedGRPCConfigOperServer) CreateSubs(*CreateSubsArgs, GRPCConfigOper_CreateSubsServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateSubs not implemented")
}
func (UnimplementedGRPCConfigOperServer) GetProtoFile(*GetProtoFileArgs, GRPCConfigOper_GetProtoFileServer) error {
	return status.Errorf(codes.Unimplemented, "method GetProtoFile not implemented")
}
func (UnimplementedGRPCConfigOperServer) mustEmbedUnimplementedGRPCConfigOperServer() {}

// UnsafeGRPCConfigOperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GRPCConfigOperServer will
// result in compilation errors.
type UnsafeGRPCConfigOperServer interface {
	mustEmbedUnimplementedGRPCConfigOperServer()
}

func RegisterGRPCConfigOperServer(s grpc.ServiceRegistrar, srv GRPCConfigOperServer) {
	s.RegisterService(&GRPCConfigOper_ServiceDesc, srv)
}

func _GRPCConfigOper_GetConfig_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConfigGetArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCConfigOperServer).GetConfig(m, &gRPCConfigOperGetConfigServer{stream})
}

type GRPCConfigOper_GetConfigServer interface {
	Send(*ConfigGetReply) error
	grpc.ServerStream
}

type gRPCConfigOperGetConfigServer struct {
	grpc.ServerStream
}

func (x *gRPCConfigOperGetConfigServer) Send(m *ConfigGetReply) error {
	return x.ServerStream.SendMsg(m)
}

func _GRPCConfigOper_MergeConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).MergeConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_MergeConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).MergeConfig(ctx, req.(*ConfigArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_DeleteConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).DeleteConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_DeleteConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).DeleteConfig(ctx, req.(*ConfigArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_RemoveConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).RemoveConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_RemoveConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).RemoveConfig(ctx, req.(*ConfigArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_ReplaceConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).ReplaceConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_ReplaceConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).ReplaceConfig(ctx, req.(*ConfigArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_CliConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CliConfigArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).CliConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_CliConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).CliConfig(ctx, req.(*CliConfigArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_CommitReplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitReplaceArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).CommitReplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_CommitReplace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).CommitReplace(ctx, req.(*CommitReplaceArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_CommitConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).CommitConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_CommitConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).CommitConfig(ctx, req.(*CommitArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_ConfigDiscardChanges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscardChangesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCConfigOperServer).ConfigDiscardChanges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCConfigOper_ConfigDiscardChanges_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCConfigOperServer).ConfigDiscardChanges(ctx, req.(*DiscardChangesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCConfigOper_GetOper_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetOperArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCConfigOperServer).GetOper(m, &gRPCConfigOperGetOperServer{stream})
}

type GRPCConfigOper_GetOperServer interface {
	Send(*GetOperReply) error
	grpc.ServerStream
}

type gRPCConfigOperGetOperServer struct {
	grpc.ServerStream
}

func (x *gRPCConfigOperGetOperServer) Send(m *GetOperReply) error {
	return x.ServerStream.SendMsg(m)
}

func _GRPCConfigOper_CreateSubs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CreateSubsArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCConfigOperServer).CreateSubs(m, &gRPCConfigOperCreateSubsServer{stream})
}

type GRPCConfigOper_CreateSubsServer interface {
	Send(*CreateSubsReply) error
	grpc.ServerStream
}

type gRPCConfigOperCreateSubsServer struct {
	grpc.ServerStream
}

func (x *gRPCConfigOperCreateSubsServer) Send(m *CreateSubsReply) error {
	return x.ServerStream.SendMsg(m)
}

func _GRPCConfigOper_GetProtoFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetProtoFileArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCConfigOperServer).GetProtoFile(m, &gRPCConfigOperGetProtoFileServer{stream})
}

type GRPCConfigOper_GetProtoFileServer interface {
	Send(*GetProtoFileReply) error
	grpc.ServerStream
}

type gRPCConfigOperGetProtoFileServer struct {
	grpc.ServerStream
}

func (x *gRPCConfigOperGetProtoFileServer) Send(m *GetProtoFileReply) error {
	return x.ServerStream.SendMsg(m)
}

// GRPCConfigOper_ServiceDesc is the grpc.ServiceDesc for GRPCConfigOper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GRPCConfigOper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IOSXRExtensibleManagabilityService.gRPCConfigOper",
	HandlerType: (*GRPCConfigOperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MergeConfig",
			Handler:    _GRPCConfigOper_MergeConfig_Handler,
		},
		{
			MethodName: "DeleteConfig",
			Handler:    _GRPCConfigOper_DeleteConfig_Handler,
		},
		{
			MethodName: "RemoveConfig",
			Handler:    _GRPCConfigOper_RemoveConfig_Handler,
		},
		{
			MethodName: "ReplaceConfig",
			Handler:    _GRPCConfigOper_ReplaceConfig_Handler,
		},
		{
			MethodName: "CliConfig",
			Handler:    _GRPCConfigOper_CliConfig_Handler,
		},
		{
			MethodName: "CommitReplace",
			Handler:    _GRPCConfigOper_CommitReplace_Handler,
		},
		{
			MethodName: "CommitConfig",
			Handler:    _GRPCConfigOper_CommitConfig_Handler,
		},
		{
			MethodName: "ConfigDiscardChanges",
			Handler:    _GRPCConfigOper_ConfigDiscardChanges_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetConfig",
			Handler:       _GRPCConfigOper_GetConfig_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetOper",
			Handler:       _GRPCConfigOper_GetOper_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CreateSubs",
			Handler:       _GRPCConfigOper_CreateSubs_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetProtoFile",
			Handler:       _GRPCConfigOper_GetProtoFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/ems/ems_grpc.proto",
}

const (
	GRPCExec_ShowCmdTextOutput_FullMethodName = "/IOSXRExtensibleManagabilityService.gRPCExec/ShowCmdTextOutput"
	GRPCExec_ShowCmdJSONOutput_FullMethodName = "/IOSXRExtensibleManagabilityService.gRPCExec/ShowCmdJSONOutput"
	GRPCExec_ActionJSON_FullMethodName        = "/IOSXRExtensibleManagabilityService.gRPCExec/ActionJSON"
)

// GRPCExecClient is the client API for GRPCExec service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GRPCExecClient interface {
	// Exec commands
	ShowCmdTextOutput(ctx context.Context, in *ShowCmdArgs, opts ...grpc.CallOption) (GRPCExec_ShowCmdTextOutputClient, error)
	ShowCmdJSONOutput(ctx context.Context, in *ShowCmdArgs, opts ...grpc.CallOption) (GRPCExec_ShowCmdJSONOutputClient, error)
	ActionJSON(ctx context.Context, in *ActionJSONArgs, opts ...grpc.CallOption) (GRPCExec_ActionJSONClient, error)
}

type gRPCExecClient struct {
	cc grpc.ClientConnInterface
}

func NewGRPCExecClient(cc grpc.ClientConnInterface) GRPCExecClient {
	return &gRPCExecClient{cc}
}

func (c *gRPCExecClient) ShowCmdTextOutput(ctx context.Context, in *ShowCmdArgs, opts ...grpc.CallOption) (GRPCExec_ShowCmdTextOutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCExec_ServiceDesc.Streams[0], GRPCExec_ShowCmdTextOutput_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCExecShowCmdTextOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCExec_ShowCmdTextOutputClient interface {
	Recv() (*ShowCmdTextReply, error)
	grpc.ClientStream
}

type gRPCExecShowCmdTextOutputClient struct {
	grpc.ClientStream
}

func (x *gRPCExecShowCmdTextOutputClient) Recv() (*ShowCmdTextReply, error) {
	m := new(ShowCmdTextReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gRPCExecClient) ShowCmdJSONOutput(ctx context.Context, in *ShowCmdArgs, opts ...grpc.CallOption) (GRPCExec_ShowCmdJSONOutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCExec_ServiceDesc.Streams[1], GRPCExec_ShowCmdJSONOutput_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCExecShowCmdJSONOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCExec_ShowCmdJSONOutputClient interface {
	Recv() (*ShowCmdJSONReply, error)
	grpc.ClientStream
}

type gRPCExecShowCmdJSONOutputClient struct {
	grpc.ClientStream
}

func (x *gRPCExecShowCmdJSONOutputClient) Recv() (*ShowCmdJSONReply, error) {
	m := new(ShowCmdJSONReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gRPCExecClient) ActionJSON(ctx context.Context, in *ActionJSONArgs, opts ...grpc.CallOption) (GRPCExec_ActionJSONClient, error) {
	stream, err := c.cc.NewStream(ctx, &GRPCExec_ServiceDesc.Streams[2], GRPCExec_ActionJSON_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCExecActionJSONClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCExec_ActionJSONClient interface {
	Recv() (*ActionJSONReply, error)
	grpc.ClientStream
}

type gRPCExecActionJSONClient struct {
	grpc.ClientStream
}

func (x *gRPCExecActionJSONClient) Recv() (*ActionJSONReply, error) {
	m := new(ActionJSONReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GRPCExecServer is the server API for GRPCExec service.
// All implementations must embed UnimplementedGRPCExecServer
// for forward compatibility
type GRPCExecServer interface {
	// Exec commands
	ShowCmdTextOutput(*ShowCmdArgs, GRPCExec_ShowCmdTextOutputServer) error
	ShowCmdJSONOutput(*ShowCmdArgs, GRPCExec_ShowCmdJSONOutputServer) error
	ActionJSON(*ActionJSONArgs, GRPCExec_ActionJSONServer) error
	mustEmbedUnimplementedGRPCExecServer()
}

// UnimplementedGRPCExecServer must be embedded to have forward compatible implementations.
type UnimplementedGRPCExecServer struct {
}

func (UnimplementedGRPCExecServer) ShowCmdTextOutput(*ShowCmdArgs, GRPCExec_ShowCmdTextOutputServer) error {
	return status.Errorf(codes.Unimplemented, "method ShowCmdTextOutput not implemented")
}
func (UnimplementedGRPCExecServer) ShowCmdJSONOutput(*ShowCmdArgs, GRPCExec_ShowCmdJSONOutputServer) error {
	return status.Errorf(codes.Unimplemented, "method ShowCmdJSONOutput not implemented")
}
func (UnimplementedGRPCExecServer) ActionJSON(*ActionJSONArgs, GRPCExec_ActionJSONServer) error {
	return status.Errorf(codes.Unimplemented, "method ActionJSON not implemented")
}
func (UnimplementedGRPCExecServer) mustEmbedUnimplementedGRPCExecServer() {}

// UnsafeGRPCExecServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GRPCExecServer will
// result in compilation errors.
type UnsafeGRPCExecServer interface {
	mustEmbedUnimplementedGRPCExecServer()
}

func RegisterGRPCExecServer(s grpc.ServiceRegistrar, srv GRPCExecServer) {
	s.RegisterService(&GRPCExec_ServiceDesc, srv)
}

func _GRPCExec_ShowCmdTextOutput_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ShowCmdArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCExecServer).ShowCmdTextOutput(m, &gRPCExecShowCmdTextOutputServer{stream})
}

type GRPCExec_ShowCmdTextOutputServer interface {
	Send(*ShowCmdTextReply) error
	grpc.ServerStream
}

type gRPCExecShowCmdTextOutputServer struct {
	grpc.ServerStream
}

func (x *gRPCExecShowCmdTextOutputServer) Send(m *ShowCmdTextReply) error {
	return x.ServerStream.SendMsg(m)
}

func _GRPCExec_ShowCmdJSONOutput_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ShowCmdArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCExecServer).ShowCmdJSONOutput(m, &gRPCExecShowCmdJSONOutputServer{stream})
}

type GRPCExec_ShowCmdJSONOutputServer interface {
	Send(*ShowCmdJSONReply) error
	grpc.ServerStream
}

type gRPCExecShowCmdJSONOutputServer struct {
	grpc.ServerStream
}

func (x *gRPCExecShowCmdJSONOutputServer) Send(m *ShowCmdJSONReply) error {
	return x.ServerStream.SendMsg(m)
}

func _GRPCExec_ActionJSON_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ActionJSONArgs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCExecServer).ActionJSON(m, &gRPCExecActionJSONServer{stream})
}

type GRPCExec_ActionJSONServer interface {
	Send(*ActionJSONReply) error
	grpc.ServerStream
}

type gRPCExecActionJSONServer struct {
	grpc.ServerStream
}

func (x *gRPCExecActionJSONServer) Send(m *ActionJSONReply) error {
	return x.ServerStream.SendMsg(m)
}

// GRPCExec_ServiceDesc is the grpc.ServiceDesc for GRPCExec service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GRPCExec_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IOSXRExtensibleManagabilityService.gRPCExec",
	HandlerType: (*GRPCExecServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ShowCmdTextOutput",
			Handler:       _GRPCExec_ShowCmdTextOutput_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ShowCmdJSONOutput",
			Handler:       _GRPCExec_ShowCmdJSONOutput_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ActionJSON",
			Handler:       _GRPCExec_ActionJSON_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/ems/ems_grpc.proto",
}

const (
	OpenConfiggRPC_GetModels_FullMethodName = "/IOSXRExtensibleManagabilityService.OpenConfiggRPC/GetModels"
)

// OpenConfiggRPCClient is the client API for OpenConfiggRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OpenConfiggRPCClient interface {
	// get-models rpc implementation per
	// github.com/openconfig/public/blob/master/release/models/rpc/openconfig-rpc.yang
	GetModels(ctx context.Context, in *GetModelsInput, opts ...grpc.CallOption) (*GetModelsOutput, error)
}

type openConfiggRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewOpenConfiggRPCClient(cc grpc.ClientConnInterface) OpenConfiggRPCClient {
	return &openConfiggRPCClient{cc}
}

func (c *openConfiggRPCClient) GetModels(ctx context.Context, in *GetModelsInput, opts ...grpc.CallOption) (*GetModelsOutput, error) {
	out := new(GetModelsOutput)
	err := c.cc.Invoke(ctx, OpenConfiggRPC_GetModels_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OpenConfiggRPCServer is the server API for OpenConfiggRPC service.
// All implementations must embed UnimplementedOpenConfiggRPCServer
// for forward compatibility
type OpenConfiggRPCServer interface {
	// get-models rpc implementation per
	// github.com/openconfig/public/blob/master/release/models/rpc/openconfig-rpc.yang
	GetModels(context.Context, *GetModelsInput) (*GetModelsOutput, error)
	mustEmbedUnimplementedOpenConfiggRPCServer()
}

// UnimplementedOpenConfiggRPCServer must be embedded to have forward compatible implementations.
type UnimplementedOpenConfiggRPCServer struct {
}

func (UnimplementedOpenConfiggRPCServer) GetModels(context.Context, *GetModelsInput) (*GetModelsOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetModels not implemented")
}
func (UnimplementedOpenConfiggRPCServer) mustEmbedUnimplementedOpenConfiggRPCServer() {}

// UnsafeOpenConfiggRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OpenConfiggRPCServer will
// result in compilation errors.
type UnsafeOpenConfiggRPCServer interface {
	mustEmbedUnimplementedOpenConfiggRPCServer()
}

func RegisterOpenConfiggRPCServer(s grpc.ServiceRegistrar, srv OpenConfiggRPCServer) {
	s.RegisterService(&OpenConfiggRPC_ServiceDesc, srv)
}

func _OpenConfiggRPC_GetModels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetModelsInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpenConfiggRPCServer).GetModels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OpenConfiggRPC_GetModels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpenConfiggRPCServer).GetModels(ctx, req.(*GetModelsInput))
	}
	return interceptor(ctx, in, info, handler)
}

// OpenConfiggRPC_ServiceDesc is the grpc.ServiceDesc for OpenConfiggRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OpenConfiggRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IOSXRExtensibleManagabilityService.OpenConfiggRPC",
	HandlerType: (*OpenConfiggRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetModels",
			Handler:    _OpenConfiggRPC_GetModels_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ems/ems_grpc.proto",
}

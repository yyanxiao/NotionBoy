// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: server.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	model "notionboy/api/pb/model"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Service_Status_FullMethodName             = "/servicev1.Service/Status"
	Service_GenrateToken_FullMethodName       = "/servicev1.Service/GenrateToken"
	Service_OAuthURL_FullMethodName           = "/servicev1.Service/OAuthURL"
	Service_OAuthCallback_FullMethodName      = "/servicev1.Service/OAuthCallback"
	Service_GenerateApiKey_FullMethodName     = "/servicev1.Service/GenerateApiKey"
	Service_DeleteApiKey_FullMethodName       = "/servicev1.Service/DeleteApiKey"
	Service_CreateConversation_FullMethodName = "/servicev1.Service/CreateConversation"
	Service_UpdateConversation_FullMethodName = "/servicev1.Service/UpdateConversation"
	Service_GetConversation_FullMethodName    = "/servicev1.Service/GetConversation"
	Service_ListConversations_FullMethodName  = "/servicev1.Service/ListConversations"
	Service_DeleteConversation_FullMethodName = "/servicev1.Service/DeleteConversation"
	Service_CreateMessage_FullMethodName      = "/servicev1.Service/CreateMessage"
	Service_GetMessage_FullMethodName         = "/servicev1.Service/GetMessage"
	Service_ListMessages_FullMethodName       = "/servicev1.Service/ListMessages"
	Service_DeleteMessage_FullMethodName      = "/servicev1.Service/DeleteMessage"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CheckStatusResponse, error)
	// GenrateToken generates a token for the user. using api key in the header.
	GenrateToken(ctx context.Context, in *model.GenrateTokenRequest, opts ...grpc.CallOption) (*model.GenrateTokenResponse, error)
	// get Oauth url
	OAuthURL(ctx context.Context, in *model.OAuthURLRequest, opts ...grpc.CallOption) (*model.OAuthURLResponse, error)
	// AuthCallback callback for oauth, will generate a token for the user
	OAuthCallback(ctx context.Context, in *model.OAuthCallbackRequest, opts ...grpc.CallOption) (*model.GenrateTokenResponse, error)
	// GenerateApiKey generate a new api key for the user
	GenerateApiKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*model.GenerateApiKeyResponse, error)
	// DeleteApiKey delete the api key for the user
	DeleteApiKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateConversation(ctx context.Context, in *model.CreateConversationRequest, opts ...grpc.CallOption) (*model.Conversation, error)
	// UpdateConversation update the conversation
	UpdateConversation(ctx context.Context, in *model.UpdateConversationRequest, opts ...grpc.CallOption) (*model.Conversation, error)
	GetConversation(ctx context.Context, in *model.GetConversationRequest, opts ...grpc.CallOption) (*model.Conversation, error)
	ListConversations(ctx context.Context, in *model.ListConversationsRequest, opts ...grpc.CallOption) (*model.ListConversationsResponse, error)
	DeleteConversation(ctx context.Context, in *model.DeleteConversationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateMessage(ctx context.Context, in *model.CreateMessageRequest, opts ...grpc.CallOption) (Service_CreateMessageClient, error)
	GetMessage(ctx context.Context, in *model.GetMessageRequest, opts ...grpc.CallOption) (*model.Message, error)
	ListMessages(ctx context.Context, in *model.ListMessagesRequest, opts ...grpc.CallOption) (*model.ListMessagesResponse, error)
	DeleteMessage(ctx context.Context, in *model.DeleteMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CheckStatusResponse, error) {
	out := new(CheckStatusResponse)
	err := c.cc.Invoke(ctx, Service_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GenrateToken(ctx context.Context, in *model.GenrateTokenRequest, opts ...grpc.CallOption) (*model.GenrateTokenResponse, error) {
	out := new(model.GenrateTokenResponse)
	err := c.cc.Invoke(ctx, Service_GenrateToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) OAuthURL(ctx context.Context, in *model.OAuthURLRequest, opts ...grpc.CallOption) (*model.OAuthURLResponse, error) {
	out := new(model.OAuthURLResponse)
	err := c.cc.Invoke(ctx, Service_OAuthURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) OAuthCallback(ctx context.Context, in *model.OAuthCallbackRequest, opts ...grpc.CallOption) (*model.GenrateTokenResponse, error) {
	out := new(model.GenrateTokenResponse)
	err := c.cc.Invoke(ctx, Service_OAuthCallback_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GenerateApiKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*model.GenerateApiKeyResponse, error) {
	out := new(model.GenerateApiKeyResponse)
	err := c.cc.Invoke(ctx, Service_GenerateApiKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteApiKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Service_DeleteApiKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CreateConversation(ctx context.Context, in *model.CreateConversationRequest, opts ...grpc.CallOption) (*model.Conversation, error) {
	out := new(model.Conversation)
	err := c.cc.Invoke(ctx, Service_CreateConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateConversation(ctx context.Context, in *model.UpdateConversationRequest, opts ...grpc.CallOption) (*model.Conversation, error) {
	out := new(model.Conversation)
	err := c.cc.Invoke(ctx, Service_UpdateConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetConversation(ctx context.Context, in *model.GetConversationRequest, opts ...grpc.CallOption) (*model.Conversation, error) {
	out := new(model.Conversation)
	err := c.cc.Invoke(ctx, Service_GetConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ListConversations(ctx context.Context, in *model.ListConversationsRequest, opts ...grpc.CallOption) (*model.ListConversationsResponse, error) {
	out := new(model.ListConversationsResponse)
	err := c.cc.Invoke(ctx, Service_ListConversations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteConversation(ctx context.Context, in *model.DeleteConversationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Service_DeleteConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CreateMessage(ctx context.Context, in *model.CreateMessageRequest, opts ...grpc.CallOption) (Service_CreateMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &Service_ServiceDesc.Streams[0], Service_CreateMessage_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceCreateMessageClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Service_CreateMessageClient interface {
	Recv() (*model.Message, error)
	grpc.ClientStream
}

type serviceCreateMessageClient struct {
	grpc.ClientStream
}

func (x *serviceCreateMessageClient) Recv() (*model.Message, error) {
	m := new(model.Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serviceClient) GetMessage(ctx context.Context, in *model.GetMessageRequest, opts ...grpc.CallOption) (*model.Message, error) {
	out := new(model.Message)
	err := c.cc.Invoke(ctx, Service_GetMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ListMessages(ctx context.Context, in *model.ListMessagesRequest, opts ...grpc.CallOption) (*model.ListMessagesResponse, error) {
	out := new(model.ListMessagesResponse)
	err := c.cc.Invoke(ctx, Service_ListMessages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteMessage(ctx context.Context, in *model.DeleteMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Service_DeleteMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	Status(context.Context, *emptypb.Empty) (*CheckStatusResponse, error)
	// GenrateToken generates a token for the user. using api key in the header.
	GenrateToken(context.Context, *model.GenrateTokenRequest) (*model.GenrateTokenResponse, error)
	// get Oauth url
	OAuthURL(context.Context, *model.OAuthURLRequest) (*model.OAuthURLResponse, error)
	// AuthCallback callback for oauth, will generate a token for the user
	OAuthCallback(context.Context, *model.OAuthCallbackRequest) (*model.GenrateTokenResponse, error)
	// GenerateApiKey generate a new api key for the user
	GenerateApiKey(context.Context, *emptypb.Empty) (*model.GenerateApiKeyResponse, error)
	// DeleteApiKey delete the api key for the user
	DeleteApiKey(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	CreateConversation(context.Context, *model.CreateConversationRequest) (*model.Conversation, error)
	// UpdateConversation update the conversation
	UpdateConversation(context.Context, *model.UpdateConversationRequest) (*model.Conversation, error)
	GetConversation(context.Context, *model.GetConversationRequest) (*model.Conversation, error)
	ListConversations(context.Context, *model.ListConversationsRequest) (*model.ListConversationsResponse, error)
	DeleteConversation(context.Context, *model.DeleteConversationRequest) (*emptypb.Empty, error)
	CreateMessage(*model.CreateMessageRequest, Service_CreateMessageServer) error
	GetMessage(context.Context, *model.GetMessageRequest) (*model.Message, error)
	ListMessages(context.Context, *model.ListMessagesRequest) (*model.ListMessagesResponse, error)
	DeleteMessage(context.Context, *model.DeleteMessageRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) Status(context.Context, *emptypb.Empty) (*CheckStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedServiceServer) GenrateToken(context.Context, *model.GenrateTokenRequest) (*model.GenrateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenrateToken not implemented")
}
func (UnimplementedServiceServer) OAuthURL(context.Context, *model.OAuthURLRequest) (*model.OAuthURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthURL not implemented")
}
func (UnimplementedServiceServer) OAuthCallback(context.Context, *model.OAuthCallbackRequest) (*model.GenrateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthCallback not implemented")
}
func (UnimplementedServiceServer) GenerateApiKey(context.Context, *emptypb.Empty) (*model.GenerateApiKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateApiKey not implemented")
}
func (UnimplementedServiceServer) DeleteApiKey(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteApiKey not implemented")
}
func (UnimplementedServiceServer) CreateConversation(context.Context, *model.CreateConversationRequest) (*model.Conversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateConversation not implemented")
}
func (UnimplementedServiceServer) UpdateConversation(context.Context, *model.UpdateConversationRequest) (*model.Conversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConversation not implemented")
}
func (UnimplementedServiceServer) GetConversation(context.Context, *model.GetConversationRequest) (*model.Conversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversation not implemented")
}
func (UnimplementedServiceServer) ListConversations(context.Context, *model.ListConversationsRequest) (*model.ListConversationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListConversations not implemented")
}
func (UnimplementedServiceServer) DeleteConversation(context.Context, *model.DeleteConversationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConversation not implemented")
}
func (UnimplementedServiceServer) CreateMessage(*model.CreateMessageRequest, Service_CreateMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}
func (UnimplementedServiceServer) GetMessage(context.Context, *model.GetMessageRequest) (*model.Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedServiceServer) ListMessages(context.Context, *model.ListMessagesRequest) (*model.ListMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMessages not implemented")
}
func (UnimplementedServiceServer) DeleteMessage(context.Context, *model.DeleteMessageRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Status(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GenrateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.GenrateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GenrateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GenrateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GenrateToken(ctx, req.(*model.GenrateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_OAuthURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.OAuthURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).OAuthURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_OAuthURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).OAuthURL(ctx, req.(*model.OAuthURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_OAuthCallback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.OAuthCallbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).OAuthCallback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_OAuthCallback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).OAuthCallback(ctx, req.(*model.OAuthCallbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GenerateApiKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GenerateApiKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GenerateApiKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GenerateApiKey(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteApiKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteApiKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_DeleteApiKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteApiKey(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CreateConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.CreateConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_CreateConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateConversation(ctx, req.(*model.CreateConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.UpdateConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_UpdateConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateConversation(ctx, req.(*model.UpdateConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.GetConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetConversation(ctx, req.(*model.GetConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ListConversations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.ListConversationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ListConversations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ListConversations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ListConversations(ctx, req.(*model.ListConversationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.DeleteConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_DeleteConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteConversation(ctx, req.(*model.DeleteConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CreateMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(model.CreateMessageRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServiceServer).CreateMessage(m, &serviceCreateMessageServer{stream})
}

type Service_CreateMessageServer interface {
	Send(*model.Message) error
	grpc.ServerStream
}

type serviceCreateMessageServer struct {
	grpc.ServerStream
}

func (x *serviceCreateMessageServer) Send(m *model.Message) error {
	return x.ServerStream.SendMsg(m)
}

func _Service_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.GetMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetMessage(ctx, req.(*model.GetMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ListMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.ListMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ListMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ListMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ListMessages(ctx, req.(*model.ListMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.DeleteMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_DeleteMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteMessage(ctx, req.(*model.DeleteMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "servicev1.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _Service_Status_Handler,
		},
		{
			MethodName: "GenrateToken",
			Handler:    _Service_GenrateToken_Handler,
		},
		{
			MethodName: "OAuthURL",
			Handler:    _Service_OAuthURL_Handler,
		},
		{
			MethodName: "OAuthCallback",
			Handler:    _Service_OAuthCallback_Handler,
		},
		{
			MethodName: "GenerateApiKey",
			Handler:    _Service_GenerateApiKey_Handler,
		},
		{
			MethodName: "DeleteApiKey",
			Handler:    _Service_DeleteApiKey_Handler,
		},
		{
			MethodName: "CreateConversation",
			Handler:    _Service_CreateConversation_Handler,
		},
		{
			MethodName: "UpdateConversation",
			Handler:    _Service_UpdateConversation_Handler,
		},
		{
			MethodName: "GetConversation",
			Handler:    _Service_GetConversation_Handler,
		},
		{
			MethodName: "ListConversations",
			Handler:    _Service_ListConversations_Handler,
		},
		{
			MethodName: "DeleteConversation",
			Handler:    _Service_DeleteConversation_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _Service_GetMessage_Handler,
		},
		{
			MethodName: "ListMessages",
			Handler:    _Service_ListMessages_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _Service_DeleteMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateMessage",
			Handler:       _Service_CreateMessage_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server.proto",
}

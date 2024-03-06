// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.0
// source: service_account.proto

package pb

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
	AccountService_CreateUser_FullMethodName      = "/pb.AccountService/CreateUser"
	AccountService_GetUser_FullMethodName         = "/pb.AccountService/GetUser"
	AccountService_GetUsers_FullMethodName        = "/pb.AccountService/GetUsers"
	AccountService_GetUserByEmail_FullMethodName  = "/pb.AccountService/GetUserByEmail"
	AccountService_UpdateUser_FullMethodName      = "/pb.AccountService/UpdateUser"
	AccountService_UpdateUserPas_FullMethodName   = "/pb.AccountService/UpdateUserPas"
	AccountService_Login_FullMethodName           = "/pb.AccountService/Login"
	AccountService_Logout_FullMethodName          = "/pb.AccountService/Logout"
	AccountService_RefreshToken_FullMethodName    = "/pb.AccountService/RefreshToken"
	AccountService_SendVertifyEmai_FullMethodName = "/pb.AccountService/SendVertifyEmai"
	AccountService_VertifyEmail_FullMethodName    = "/pb.AccountService/VertifyEmail"
	AccountService_SSOGoogleLogin_FullMethodName  = "/pb.AccountService/SSOGoogleLogin"
	AccountService_ValidateToken_FullMethodName   = "/pb.AccountService/ValidateToken"
)

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserDTOResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserDTOResponse, error)
	GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*UserDTOsResponse, error)
	GetUserByEmail(ctx context.Context, in *GetUserByEmailRequest, opts ...grpc.CallOption) (*UserDTOResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserDTOResponse, error)
	UpdateUserPas(ctx context.Context, in *UpdateUserPasRequest, opts ...grpc.CallOption) (*UserDTOResponse, error)
	Login(ctx context.Context, in *LoginRequset, opts ...grpc.CallOption) (*AuthDTOResponse, error)
	Logout(ctx context.Context, in *LogoutRequset, opts ...grpc.CallOption) (*AuthDTOResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenRequset, opts ...grpc.CallOption) (*AuthDTOResponse, error)
	SendVertifyEmai(ctx context.Context, in *SendVertifyEmailRequset, opts ...grpc.CallOption) (*VertifyEmailResponse, error)
	VertifyEmail(ctx context.Context, in *VertifyEmailRequset, opts ...grpc.CallOption) (*VertifyEmailResponse, error)
	SSOGoogleLogin(ctx context.Context, in *GoogleIDTokenRequest, opts ...grpc.CallOption) (*AuthDTOResponse, error)
	ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*AuthDTOResponse, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserDTOResponse, error) {
	out := new(UserDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserDTOResponse, error) {
	out := new(UserDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_GetUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*UserDTOsResponse, error) {
	out := new(UserDTOsResponse)
	err := c.cc.Invoke(ctx, AccountService_GetUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUserByEmail(ctx context.Context, in *GetUserByEmailRequest, opts ...grpc.CallOption) (*UserDTOResponse, error) {
	out := new(UserDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_GetUserByEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserDTOResponse, error) {
	out := new(UserDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_UpdateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateUserPas(ctx context.Context, in *UpdateUserPasRequest, opts ...grpc.CallOption) (*UserDTOResponse, error) {
	out := new(UserDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_UpdateUserPas_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Login(ctx context.Context, in *LoginRequset, opts ...grpc.CallOption) (*AuthDTOResponse, error) {
	out := new(AuthDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Logout(ctx context.Context, in *LogoutRequset, opts ...grpc.CallOption) (*AuthDTOResponse, error) {
	out := new(AuthDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_Logout_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequset, opts ...grpc.CallOption) (*AuthDTOResponse, error) {
	out := new(AuthDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_RefreshToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) SendVertifyEmai(ctx context.Context, in *SendVertifyEmailRequset, opts ...grpc.CallOption) (*VertifyEmailResponse, error) {
	out := new(VertifyEmailResponse)
	err := c.cc.Invoke(ctx, AccountService_SendVertifyEmai_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) VertifyEmail(ctx context.Context, in *VertifyEmailRequset, opts ...grpc.CallOption) (*VertifyEmailResponse, error) {
	out := new(VertifyEmailResponse)
	err := c.cc.Invoke(ctx, AccountService_VertifyEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) SSOGoogleLogin(ctx context.Context, in *GoogleIDTokenRequest, opts ...grpc.CallOption) (*AuthDTOResponse, error) {
	out := new(AuthDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_SSOGoogleLogin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*AuthDTOResponse, error) {
	out := new(AuthDTOResponse)
	err := c.cc.Invoke(ctx, AccountService_ValidateToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility
type AccountServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*UserDTOResponse, error)
	GetUser(context.Context, *GetUserRequest) (*UserDTOResponse, error)
	GetUsers(context.Context, *GetUsersRequest) (*UserDTOsResponse, error)
	GetUserByEmail(context.Context, *GetUserByEmailRequest) (*UserDTOResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UserDTOResponse, error)
	UpdateUserPas(context.Context, *UpdateUserPasRequest) (*UserDTOResponse, error)
	Login(context.Context, *LoginRequset) (*AuthDTOResponse, error)
	Logout(context.Context, *LogoutRequset) (*AuthDTOResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequset) (*AuthDTOResponse, error)
	SendVertifyEmai(context.Context, *SendVertifyEmailRequset) (*VertifyEmailResponse, error)
	VertifyEmail(context.Context, *VertifyEmailRequset) (*VertifyEmailResponse, error)
	SSOGoogleLogin(context.Context, *GoogleIDTokenRequest) (*AuthDTOResponse, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*AuthDTOResponse, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccountServiceServer struct {
}

func (UnimplementedAccountServiceServer) CreateUser(context.Context, *CreateUserRequest) (*UserDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedAccountServiceServer) GetUser(context.Context, *GetUserRequest) (*UserDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedAccountServiceServer) GetUsers(context.Context, *GetUsersRequest) (*UserDTOsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedAccountServiceServer) GetUserByEmail(context.Context, *GetUserByEmailRequest) (*UserDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByEmail not implemented")
}
func (UnimplementedAccountServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UserDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedAccountServiceServer) UpdateUserPas(context.Context, *UpdateUserPasRequest) (*UserDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserPas not implemented")
}
func (UnimplementedAccountServiceServer) Login(context.Context, *LoginRequset) (*AuthDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccountServiceServer) Logout(context.Context, *LogoutRequset) (*AuthDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAccountServiceServer) RefreshToken(context.Context, *RefreshTokenRequset) (*AuthDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedAccountServiceServer) SendVertifyEmai(context.Context, *SendVertifyEmailRequset) (*VertifyEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendVertifyEmai not implemented")
}
func (UnimplementedAccountServiceServer) VertifyEmail(context.Context, *VertifyEmailRequset) (*VertifyEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VertifyEmail not implemented")
}
func (UnimplementedAccountServiceServer) SSOGoogleLogin(context.Context, *GoogleIDTokenRequest) (*AuthDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SSOGoogleLogin not implemented")
}
func (UnimplementedAccountServiceServer) ValidateToken(context.Context, *ValidateTokenRequest) (*AuthDTOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUsers(ctx, req.(*GetUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUserByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetUserByEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserByEmail(ctx, req.(*GetUserByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdateUserPas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserPasRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdateUserPas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_UpdateUserPas_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdateUserPas(ctx, req.(*UpdateUserPasRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Login(ctx, req.(*LoginRequset))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Logout(ctx, req.(*LogoutRequset))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_RefreshToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequset))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_SendVertifyEmai_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendVertifyEmailRequset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).SendVertifyEmai(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_SendVertifyEmai_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).SendVertifyEmai(ctx, req.(*SendVertifyEmailRequset))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_VertifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VertifyEmailRequset)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).VertifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_VertifyEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).VertifyEmail(ctx, req.(*VertifyEmailRequset))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_SSOGoogleLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoogleIDTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).SSOGoogleLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_SSOGoogleLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).SSOGoogleLogin(ctx, req.(*GoogleIDTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_ValidateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).ValidateToken(ctx, req.(*ValidateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _AccountService_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _AccountService_GetUser_Handler,
		},
		{
			MethodName: "GetUsers",
			Handler:    _AccountService_GetUsers_Handler,
		},
		{
			MethodName: "GetUserByEmail",
			Handler:    _AccountService_GetUserByEmail_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _AccountService_UpdateUser_Handler,
		},
		{
			MethodName: "UpdateUserPas",
			Handler:    _AccountService_UpdateUserPas_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AccountService_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AccountService_Logout_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _AccountService_RefreshToken_Handler,
		},
		{
			MethodName: "SendVertifyEmai",
			Handler:    _AccountService_SendVertifyEmai_Handler,
		},
		{
			MethodName: "VertifyEmail",
			Handler:    _AccountService_VertifyEmail_Handler,
		},
		{
			MethodName: "SSOGoogleLogin",
			Handler:    _AccountService_SSOGoogleLogin_Handler,
		},
		{
			MethodName: "ValidateToken",
			Handler:    _AccountService_ValidateToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_account.proto",
}

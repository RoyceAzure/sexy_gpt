package gapi

import (
	"context"
	"fmt"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util/validate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func processResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.UserDTOResponse, error) {
	res := pb.UserDTOResponse{Message: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		headerMD := metadata.Pairs(
			util.DBMSGKey, err.Error(),
		)

		grpc.SendHeader(ctx, headerMD)
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func processResponses(ctx context.Context, code codes.Code, msg string, err error) (*pb.UserDTOsResponse, error) {
	res := pb.UserDTOsResponse{Message: msg}
	if code != codes.OK {
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		headerMD := metadata.Pairs(
			util.DBMSGKey, err.Error(),
		)

		grpc.SendHeader(ctx, headerMD)
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func processAuthResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.AuthDTOResponse, error) {
	res := pb.AuthDTOResponse{Message: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		// util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		headerMD := metadata.Pairs(
			util.DBMSGKey, err.Error(),
		)
		grpc.SendHeader(ctx, headerMD)
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func processVertifyEmailResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.VertifyEmailResponse, error) {
	res := pb.VertifyEmailResponse{Message: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		headerMD := metadata.Pairs(
			util.DBMSGKey, err.Error(),
		)

		grpc.SendHeader(ctx, headerMD)
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func (server *AccountServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserDTOResponse, error) {
	if violations := validateCreateUserReq(req); violations != nil {
		return nil, gpt_error.InvalidArgumentError(violations)
	}

	return server.accountServiceDao.CreateUser(ctx, req)
}

func (server *AccountServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserDTOResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.GetUser(ctx, req, token)
}

func (server *AccountServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.UserDTOsResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processResponses(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.GetUsers(ctx, req, token)
}

func (server *AccountServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserDTOResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.GetUserByEmail(ctx, req, token)
}

func (server *AccountServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDTOResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.UpdateUser(ctx, req, token)
}

func (server *AccountServer) UpdateUserPas(ctx context.Context, req *pb.UpdateUserPasRequest) (*pb.UserDTOResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.UpdateUserPas(ctx, req, token)
}

func (server *AccountServer) Login(ctx context.Context, req *pb.LoginRequset) (*pb.AuthDTOResponse, error) {
	return server.accountServiceDao.Login(ctx, req)
}

func (server *AccountServer) Logout(ctx context.Context, req *pb.LogoutRequset) (*pb.AuthDTOResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.Logout(ctx, req, token)
}

func (server *AccountServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequset) (*pb.AuthDTOResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processAuthResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	return server.accountServiceDao.RefreshToken(ctx, req, token)
}

func (server *AccountServer) SendVertifyEmai(ctx context.Context, req *pb.SendVertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	return server.accountServiceDao.SendVertifyEmai(ctx, req)
}

func (server *AccountServer) VertifyEmail(ctx context.Context, req *pb.VertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	return server.accountServiceDao.VertifyEmail(ctx, req)

}

func (server *AccountServer) SSOGoogleLogin(ctx context.Context, req *pb.GoogleIDTokenRequest) (*pb.AuthDTOResponse, error) {
	return server.accountServiceDao.SSOGoogleLogin(ctx, req)
}

func validateCreateUserReq(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmailFormat(req.GetEmail()); err != nil {
		violations = append(violations, validate.FieldViolation("email", err))
	}

	if err := validate.ValidateStrongPas(req.GetPassword()); err != nil {
		violations = append(violations, validate.FieldViolations("password", err))
	}

	if err := validate.ValidateEmptyString(req.GetUserName()); err != nil {
		violations = append(violations, validate.FieldViolation("user_name", err))
	}
	return violations
}

// 如果有userId, 則驗證是否為uuid格式
func validateGetUserReq(req *pb.GetUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetUserId() != "" {
		if err := validate.ValidateUUID(req.GetUserId()); err != nil {
			violations = append(violations, validate.FieldViolation("users_id", err))
		}
	}
	return violations
}

func validateUpdateUserReq(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}
	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}
	if err := validate.ValidateEmptyString(req.GetUserName()); err != nil {
		violations = append(violations, validate.FieldViolation("user_name", err))
	}
	return violations
}

func validateUpdateUserPasReq(req *pb.UpdateUserPasRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmptyString(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}
	if err := validate.ValidateUUID(req.GetUserId()); err != nil {
		violations = append(violations, validate.FieldViolation("users_id", err))
	}
	if err := validate.ValidateStrongPas(req.GetPassword()); err != nil {
		violations = append(violations, validate.FieldViolations("password", err))
	}
	return violations
}

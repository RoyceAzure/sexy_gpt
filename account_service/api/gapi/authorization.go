package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/model"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authorizUser(ctx context.Context) (*model.AuthUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, gpt_error.ErrInValidatePreConditionOp.ErrStr("misssing metadata")
	}

	values := md.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, gpt_error.ErrInValidatePreConditionOp.ErrStr("misssing authorization header")
	}
	authHeader := strings.Fields(values[0])
	if len(authHeader) != 2 {

		return nil, gpt_error.ErrInValidatePreConditionOp.ErrStr("invalid auth format")
	}
	authType := strings.ToLower(authHeader[0])
	if authType != authorizationTypeBearer {
		return nil, gpt_error.ErrInValidatePreConditionOp.ErrStr(fmt.Sprintf("unsportted authorization type : %s", authType))
	}

	accessToken := authHeader[1]
	payload, err := server.tokenMaker.VertifyToken(accessToken)
	if err != nil {
		return nil, gpt_error.ErrInValidatePreConditionOp.ErrStr("invalid token")
	}

	user, err := server.Service.IsValidateUser(ctx, payload.Subject.UPN)
	if err != nil {
		return nil, err
	}

	_, err = server.Service.IsUserLogin(ctx, payload.Subject.UserId)
	if err != nil {
		return nil, err
	}

	res := model.AuthUser{
		UserId:     user.UserID.Bytes,
		UserName:   user.UserName,
		Email:      user.Email,
		RoleId:     user.RoleID.Bytes,
		RoleName:   user.RoleName.String,
		IsInternal: user.IsInternal,
	}

	return &res, nil
}

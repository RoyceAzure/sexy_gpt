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

/*
des:

	驗證使用者token，錯誤直接回傳給使用者

parm:

	ctx

errors:

	ErrUnauthicated : token或其內容有誤
	ErrInvalidSession : session不存在, session過期, session被block, 需要重新登入
	ErrNotFound: user 不存在
	ErrInternal : 未預期錯誤
*/
func (server *Server) authorizUser(ctx context.Context) (*model.AuthUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("misssing metadata, %w", gpt_error.ErrUnauthicated)
	}

	values := md.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("misssing authorization header, %w", gpt_error.ErrUnauthicated)
	}
	authHeader := strings.Fields(values[0])
	if len(authHeader) != 2 {
		return nil, fmt.Errorf("invalid auth format, %w", gpt_error.ErrUnauthicated)
	}
	authType := strings.ToLower(authHeader[0])
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("unsportted authorization type : %s, %w", authType, gpt_error.ErrUnauthicated)
	}

	accessToken := authHeader[1]
	payload, err := server.tokenMaker.VertifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token, %w", gpt_error.ErrUnauthicated)
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

package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authorizUser(ctx context.Context) (*token.TokenPayload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("misssing metadata")
	}

	values := md.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("misssing authorization header")
	}
	authHeader := strings.Fields(values[0])
	if len(authHeader) != 2 {
		return nil, fmt.Errorf("invalid auth format")
	}
	authType := strings.ToLower(authHeader[0])
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("unsportted authorization type : %s", authType)
	}

	accessToken := authHeader[1]
	payload, err := server.tokenMaker.VertifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invaliad access token : %s", err)
	}

	err = server.userService.IsValidateUser(ctx, payload.Subject.UPN)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

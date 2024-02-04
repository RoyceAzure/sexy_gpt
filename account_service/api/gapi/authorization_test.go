package gapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	mock_service "github.com/RoyceAzure/sexy_gpt/account_service/service/mock"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/model"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthorizUser(t *testing.T) {

	key := random.RandomString(32)
	email := random.RandomEmailString(5)
	subject := &token.TokenSubject{
		UPN: email,
	}
	testCases := []struct {
		name          string
		buildStub     func(service *mock_service.MockIService)
		buildContext  func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context
		checkResponse func(t *testing.T, res *model.AuthUser, err error)
	}{
		{
			name: "ok",
			buildStub: func(service *mock_service.MockIService) {
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *model.AuthUser, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, s.Code())
				require.NoError(t, err)
			},
		},
		{
			name: "misssing metadata",
			buildStub: func(service *mock_service.MockIService) {
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *model.AuthUser, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.True(t, errors.Is(err, gpt_error.ErrUnauthicated))
			},
		},
		{
			name: "invaliad access token",
			buildStub: func(service *mock_service.MockIService) {
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				otherKey := random.RandomString(32)
				otherMaker, _ := token.NewPasetoMaker(otherKey)
				return newContextWithBearerToken(t, otherMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *model.AuthUser, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.True(t, errors.Is(err, gpt_error.ErrUnauthicated))
			},
		},
		{
			name: "user not found",
			buildStub: func(service *mock_service.MockIService) {
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(nil, gpt_error.ErrNotFound)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *model.AuthUser, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.True(t, errors.Is(err, gpt_error.ErrNotFound))
			},
		},
		{
			name: "user not validated email",
			buildStub: func(service *mock_service.MockIService) {
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(nil, gpt_error.ErrUnauthicated)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *model.AuthUser, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.True(t, errors.Is(err, gpt_error.ErrUnauthicated))
			},
		},
		{
			name: "user session invalidated",
			buildStub: func(service *mock_service.MockIService) {
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(nil, gpt_error.ErrInvalidSession)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *model.AuthUser, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.True(t, errors.Is(err, gpt_error.ErrInvalidSession))
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mock_service.NewMockIService(ctrl)

		tc.buildStub(service)

		tokenMaker, err := token.NewPasetoMaker(key)

		require.NoError(t, err)

		server, err := NewServer(config.Config{}, nil, tokenMaker, service)

		require.NoError(t, err)

		ctx := tc.buildContext(t, tokenMaker, "audi", "isur", time.Hour*1)

		res, err := server.authorizUser(ctx)

		tc.checkResponse(t, res, err)
	}
}

package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	mock_db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/mock"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	mock_service "github.com/RoyceAzure/sexy_gpt/account_service/service/mock"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	userName := random.RandomString(10)
	email := random.RandomString(5) + "@gmail.com"
	password := "a123A68789"
	roleName := random.RandomString(5)
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStub     func(dao *mock_db.MockDao, service *mock_service.MockIService)
		checkResponse func(t *testing.T, res *pb.UserDTOResponse, err error)
	}{
		{
			name: "ok",
			req: &pb.CreateUserRequest{
				UserName: userName,
				Email:    email,
				Password: password,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.CreateUserTxResults{
					User: db.User{
						UserName: userName,
						Email:    email,
					},
					Role: db.Role{
						RoleName: roleName,
					},
				}, nil)

				service.EXPECT().
					SendVertifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.NotNil(t, res)

				require.Equal(t, email, res.Data.Email)
				require.Equal(t, userName, res.Data.UserName)
				require.Equal(t, roleName, res.Data.RoleName)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, s.Code())
				require.NoError(t, err)
			},
		},
		{
			name: "invalidate par",
			req: &pb.CreateUserRequest{
				UserName: "",
				Email:    "",
				Password: "'",
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0).Return(db.CreateUserTxResults{}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Nil(t, res)
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, s.Code())
			},
		},
		{
			name: "invalidate par email",
			req: &pb.CreateUserRequest{
				UserName: userName,
				Email:    "fsafs3215",
				Password: "dsafsAA123",
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0).Return(db.CreateUserTxResults{}, nil)

				service.EXPECT().
					SendVertifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, s.Code())
			},
		},
		{
			name: "invalidate par email",
			req: &pb.CreateUserRequest{
				UserName: userName,
				Email:    email,
				Password: "dsafsAA",
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0).Return(db.CreateUserTxResults{}, nil)

				service.EXPECT().
					SendVertifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, s.Code())
			},
		},
		{
			name: "uniqueViolation",
			req: &pb.CreateUserRequest{
				UserName: userName,
				Email:    email,
				Password: password,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.CreateUserTxResults{}, &pgconn.PgError{
					Code: "23505",
				})

				service.EXPECT().
					SendVertifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, s.Code())
			},
		},
		{
			name: "Internal error",
			req: &pb.CreateUserRequest{
				UserName: userName,
				Email:    email,
				Password: password,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.CreateUserTxResults{}, fmt.Errorf("other err"))

				service.EXPECT().
					SendVertifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock_dao := mock_db.NewMockDao(ctrl)

		ctrl2 := gomock.NewController(t)
		defer ctrl2.Finish()
		mock_service := mock_service.NewMockIService(ctrl2)

		tc.buildStub(mock_dao, mock_service)
		server, err := NewServer(config.Config{}, mock_dao, nil, mock_service)

		require.NoError(t, err)

		res, err := server.CreateUser(context.Background(), tc.req)

		tc.checkResponse(t, res, err)
	}

}

/*
subject只要解的出來就行  裡面資料不用正確
mock userService isVaildate 回傳要測試的資料就行
*/
func TestGetUsers(t *testing.T) {
	key := random.RandomString(32)
	email := random.RandomEmailString(5)
	subject := &token.TokenSubject{
		UPN: email,
	}

	testCases := []struct {
		name          string
		req           *pb.GetUsersRequest
		buildStub     func(dao *mock_db.MockDao, service *mock_service.MockIService)
		buildContext  func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context
		checkResponse func(t *testing.T, res *pb.UserDTOsResponse, err error)
	}{
		{
			name: "ok",
			req:  &pb.GetUsersRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUsersDTO(gomock.Any(), gomock.Eq(db.GetUsersDTOParams{
						Limit:  10,
						Offset: 0,
					})).
					Times(1).Return([]db.UserRoleView{}, nil)

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)

				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{IsInternal: true}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOsResponse, err error) {
				require.NotNil(t, res)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, s.Code())
				require.NoError(t, err)
			},
		},
		{
			name: "Internal err",
			req:  &pb.GetUsersRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUsersDTO(gomock.Any(), gomock.Eq(db.GetUsersDTOParams{
						Limit:  10,
						Offset: 0,
					})).
					Times(1).Return([]db.UserRoleView{}, fmt.Errorf("any err"))

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)

				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{IsInternal: true}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOsResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
		{
			name: "validated failed",
			req:  &pb.GetUsersRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUsersDTO(gomock.Any(), gomock.Any()).
					Times(0)
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{IsInternal: false}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOsResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, s.Code())
			},
		},
		{
			name: "validated failed user not vertified",
			req:  &pb.GetUsersRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUsersDTO(gomock.Any(), gomock.Any()).
					Times(0)

				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(nil, gpt_error.ErrUnauthicated)

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOsResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, s.Code())
			},
		},
		{
			name: "validated failed user not found",
			req:  &pb.GetUsersRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUsersDTO(gomock.Any(), gomock.Any()).
					Times(0)

				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(nil, gpt_error.ErrNotFound)
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOsResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, s.Code())
			},
		},
		{
			name: "validated failed internal ",
			req:  &pb.GetUsersRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUsersDTO(gomock.Any(), gomock.Any()).
					Times(0)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(nil, gpt_error.ErrInternal)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOsResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock_dao := mock_db.NewMockDao(ctrl)

		ctrl2 := gomock.NewController(t)
		defer ctrl2.Finish()

		service := mock_service.NewMockIService(ctrl2)

		tc.buildStub(mock_dao, service)

		tokenMaker, err := token.NewPasetoMaker(key)

		require.NoError(t, err)

		server, err := NewServer(config.Config{}, mock_dao, tokenMaker, service)

		require.NoError(t, err)

		ctx := tc.buildContext(t, tokenMaker, "audi", "isur", time.Hour*1)

		res, err := server.GetUsers(ctx, tc.req)

		tc.checkResponse(t, res, err)
	}
}

func TestGetUser(t *testing.T) {
	key := random.RandomString(32)
	email := random.RandomEmailString(5)
	subject := &token.TokenSubject{
		UPN: email,
	}
	testCases := []struct {
		name          string
		req           *pb.GetUserRequest
		buildStub     func(dao *mock_db.MockDao, service *mock_service.MockIService)
		buildContext  func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context
		checkResponse func(t *testing.T, res *pb.UserDTOResponse, err error)
	}{
		{
			name: "ok",
			req: &pb.GetUserRequest{
				UserId: uuid.New().String(),
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUserDTO(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserRoleView{}, nil)
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Data)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, s.Code())
				require.NoError(t, err)
			},
		},
		{
			name: "user not exists",
			req: &pb.GetUserRequest{
				UserId: uuid.New().String(),
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUserDTO(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserRoleView{}, gpt_error.DB_ERR_NOT_FOUND)
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, s.Code())
			},
		},
		{
			name: "Internal",
			req: &pb.GetUserRequest{
				UserId: uuid.New().String(),
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUserDTO(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserRoleView{}, fmt.Errorf("other err"))

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock_dao := mock_db.NewMockDao(ctrl)

		ctrl2 := gomock.NewController(t)
		defer ctrl2.Finish()

		service := mock_service.NewMockIService(ctrl2)

		tc.buildStub(mock_dao, service)

		tokenMaker, err := token.NewPasetoMaker(key)

		require.NoError(t, err)

		server, err := NewServer(config.Config{}, mock_dao, tokenMaker, service)

		require.NoError(t, err)

		ctx := tc.buildContext(t, tokenMaker, "audi", "isur", time.Hour*1)

		res, err := server.GetUser(ctx, tc.req)

		tc.checkResponse(t, res, err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	key := random.RandomString(32)
	email := random.RandomEmailString(5)
	subject := &token.TokenSubject{
		UPN: email,
	}
	testCases := []struct {
		name          string
		req           *pb.GetUserByEmailRequest
		buildStub     func(dao *mock_db.MockDao, service *mock_service.MockIService)
		buildContext  func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context
		checkResponse func(t *testing.T, res *pb.UserDTOResponse, err error)
	}{
		{
			name: "ok",
			req: &pb.GetUserByEmailRequest{
				Email: email,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUserDTOByEmail(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserRoleView{}, nil)
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.NoError(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, s.Code())
				require.NoError(t, err)
			},
		},
		// {
		// 	name: "bad req",
		// 	req:  &pb.GetUserByEmailRequest{},
		// 	buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
		// 		dao.EXPECT().
		// 			GetUserDTOByEmail(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 		service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
		// 			Return(true, nil)
		// 		service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
		// 			Return(&db.UserRoleView{}, nil)
		// 	},
		// 	buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
		// 		return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
		// 		require.Error(t, err)
		// 		s, ok := status.FromError(err)
		// 		require.True(t, ok)
		// 		require.Equal(t, codes.InvalidArgument, s.Code())
		// 	},
		// },
		{
			name: "user not exists",
			req: &pb.GetUserByEmailRequest{
				Email: email,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUserDTOByEmail(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserRoleView{}, gpt_error.DB_ERR_NOT_FOUND)
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, s.Code())
			},
		},
		{
			name: "Internal",
			req: &pb.GetUserByEmailRequest{
				Email: email,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					GetUserDTOByEmail(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserRoleView{}, fmt.Errorf("other err"))
				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock_dao := mock_db.NewMockDao(ctrl)

		ctrl2 := gomock.NewController(t)
		defer ctrl2.Finish()

		service := mock_service.NewMockIService(ctrl2)

		tc.buildStub(mock_dao, service)

		tokenMaker, err := token.NewPasetoMaker(key)

		require.NoError(t, err)

		server, err := NewServer(config.Config{}, mock_dao, tokenMaker, service)

		require.NoError(t, err)

		ctx := tc.buildContext(t, tokenMaker, "audi", "isur", time.Hour*1)

		res, err := server.GetUserByEmail(ctx, tc.req)

		tc.checkResponse(t, res, err)
	}
}
func TestUpdateUser(t *testing.T) {
	key := random.RandomString(32)
	email := random.RandomEmailString(5)
	userName := random.RandomEmailString(5)
	userId := uuid.New()
	subject := &token.TokenSubject{
		UPN: email,
	}
	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStub     func(dao *mock_db.MockDao, service *mock_service.MockIService)
		buildContext  func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context
		checkResponse func(t *testing.T, res *pb.UserDTOResponse, err error)
	}{
		{
			name: "ok",
			req: &pb.UpdateUserRequest{
				UserId:   userId.String(),
				UserName: userName,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(db.User{}, nil)

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Data)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, s.Code())
				require.NoError(t, err)
			},
		},
		{
			name: "bad req",
			req:  &pb.UpdateUserRequest{},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, s.Code())
			},
		},
		{
			name: "user not exists",
			req: &pb.UpdateUserRequest{
				UserId:   userId.String(),
				UserName: userName,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(db.User{}, gpt_error.DB_ERR_NOT_FOUND)

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, s.Code())
			},
		},
		{
			name: "Internal",
			req: &pb.UpdateUserRequest{
				UserId:   userId.String(),
				UserName: userName,
			},
			buildStub: func(dao *mock_db.MockDao, service *mock_service.MockIService) {
				dao.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(db.User{}, fmt.Errorf("other err"))

				service.EXPECT().IsUserLogin(gomock.Any(), gomock.Any()).Times(1).
					Return(true, nil)
				service.EXPECT().IsValidateUser(gomock.Any(), gomock.Eq(email)).Times(1).
					Return(&db.UserRoleView{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker, audience string, issuer string, duration time.Duration) context.Context {
				return newContextWithBearerToken(t, tokenMaker, subject, audience, issuer, duration)
			},
			checkResponse: func(t *testing.T, res *pb.UserDTOResponse, err error) {
				require.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock_dao := mock_db.NewMockDao(ctrl)

		ctrl2 := gomock.NewController(t)
		defer ctrl2.Finish()

		service := mock_service.NewMockIService(ctrl2)

		tc.buildStub(mock_dao, service)

		tokenMaker, err := token.NewPasetoMaker(key)

		require.NoError(t, err)

		server, err := NewServer(config.Config{}, mock_dao, tokenMaker, service)

		require.NoError(t, err)

		ctx := tc.buildContext(t, tokenMaker, "audi", "isur", time.Hour*1)

		res, err := server.UpdateUser(ctx, tc.req)

		tc.checkResponse(t, res, err)
	}
}

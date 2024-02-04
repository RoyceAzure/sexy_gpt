package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock_db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/mock"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestIsValidateUser(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		buildStub     func(dao *mock_db.MockDao)
		checkResponse func(t *testing.T, actual error)
	}{
		{
			name:  "ok",
			email: "royce",
			buildStub: func(dao *mock_db.MockDao) {
				dao.EXPECT().GetUserDTOByEmail(gomock.Any(), gomock.Eq("royce")).
					Times(1).Return(db.UserRoleView{
					IsEmailVerified: true,
				}, nil)
			},
			checkResponse: func(t *testing.T, actual error) {
				require.Equal(t, nil, actual)
			},
		},
		{
			name:  "user is not validated",
			email: "royce",
			buildStub: func(dao *mock_db.MockDao) {
				dao.EXPECT().GetUserDTOByEmail(gomock.Any(), gomock.Eq("royce")).
					Times(1).Return(db.UserRoleView{
					IsEmailVerified: false,
				}, nil)
			},
			checkResponse: func(t *testing.T, actual error) {
				require.True(t, errors.Is(actual, gpt_error.ErrInValidatePreConditionOp))
			},
		},
		{
			name:  "user is not exists",
			email: "royce",
			buildStub: func(dao *mock_db.MockDao) {
				dao.EXPECT().GetUserDTOByEmail(gomock.Any(), gomock.Eq("royce")).
					Times(1).Return(db.UserRoleView{}, gpt_error.DB_ERR_NOT_FOUND)
			},
			checkResponse: func(t *testing.T, actual error) {
				require.True(t, errors.Is(actual, gpt_error.ErrNotFound))
			},
		},
		{
			name:  "Internal error",
			email: "royce",
			buildStub: func(dao *mock_db.MockDao) {
				dao.EXPECT().GetUserDTOByEmail(gomock.Any(), gomock.Eq("royce")).
					Times(1).Return(db.UserRoleView{}, fmt.Errorf("something else"))
			},
			checkResponse: func(t *testing.T, actual error) {
				require.True(t, errors.Is(actual, gpt_error.ErrInternal))
			},
		},
	}

	for _, tc := range testCases {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock_dao := mock_db.NewMockDao(ctrl)

		tc.buildStub(mock_dao)

		userService := NewService(mock_dao, nil)

		_, err := userService.IsValidateUser(context.Background(), tc.email)

		tc.checkResponse(t, err)
	}
}

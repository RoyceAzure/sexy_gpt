package service

import (
	"context"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
)

type IUserService interface {
	IsValidateUser(ctx context.Context, email string) (*db.UserRoleView, error)
}

func (userService *Service) IsValidateUser(ctx context.Context, email string) (*db.UserRoleView, error) {
	user, err := userService.dao.GetUserDTOByEmail(ctx, email)
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, gpt_error.ErrNotFound.ErrStr("user not exists")
		}
		return nil, gpt_error.ErrInternal
	}

	if !user.IsEmailVerified {
		return nil, gpt_error.ErrInValidatePreConditionOp.ErrStr("user is not validated")
	}

	return &user, nil
}

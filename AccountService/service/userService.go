package service

import (
	"context"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"

)

type IUserService interface {
	IsValidateUser(ctx context.Context, email string) error
}

type UserService struct {
	dao db.Dao
}

func NewUserService(dao db.Dao) *UserService {
	return &UserService{
		dao: dao,
	}
}

func (userService *UserService) IsValidateUser(ctx context.Context, email string) error {
	user, err := userService.dao.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == gpt_error.ERR_NOT_FOUND.Error() {
			return gpt_error.ErrNotFound.ErrStr("user not exists")
		}
		return gpt_error.ErrInternal
	}

	if !user.IsEmailVerified {
		return gpt_error.ErrInValidatePreConditionOp.ErrStr("user is not validated")
	}

	return nil
}

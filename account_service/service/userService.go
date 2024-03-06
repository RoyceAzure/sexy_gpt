package service

import (
	"context"
	"fmt"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
)

type IUserService interface {
	/*
		des:
			驗證user email是否驗證過
		parm:
			email:  user email
		errors:
			ErrNotFound : user 不存在
			ErrUnauthicated : user email沒有驗證
			ErrInternal : 未預期錯誤
	*/
	IsValidateUser(ctx context.Context, email string) (*db.UserRoleView, error)
}

func (userService *Service) IsValidateUser(ctx context.Context, email string) (*db.UserRoleView, error) {
	user, err := userService.dao.GetUserDTOByEmail(ctx, email)
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return nil, fmt.Errorf("user not exists, %w", gpt_error.ErrNotFound)
		}
		return nil, gpt_error.ErrInternal
	}

	if !user.IsEmailVerified {
		return nil, fmt.Errorf("user is not validated, %w", gpt_error.ErrUnauthicated)
	}

	return &user, nil
}

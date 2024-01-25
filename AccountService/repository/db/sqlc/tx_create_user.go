package db

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/constants"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/jackc/pgx/v5"
)

type CreateUserTxResults struct {
	User User
	Role Role
}

func (dao *SQLDao) CreateUserTx(ctx context.Context, parm *CreateUserParams) (CreateUserTxResults, error) {
	var result CreateUserTxResults
	err := dao.execTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, *parm)
		if err != nil {
			return err
		}

		role, err := q.GetRoleByRoleName(ctx, constants.DEFAULT_USER_ROLE)
		if err != nil {
			if err.Error() == gpt_error.ERR_NOT_FOUND.Error() {
				return gpt_error.ErrInternal
			}
			return err
		}

		_, err = q.CreateUserRole(ctx, CreateUserRoleParams{
			UserID: result.User.UserID,
			RoleID: role.RoleID,
			CrUser: constants.SYSTEM_USER,
		})
		if err != nil {
			return err
		}
		result.Role = role

		return nil
	})
	return result, err
}

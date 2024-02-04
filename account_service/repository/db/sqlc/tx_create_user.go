package db

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/constants"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/jackc/pgx/v5"
)

type CreateUserTxParms struct {
	Arg         *CreateUserParams
	AfterCreate func(User) error
}

type CreateUserTxResults struct {
	User User
	Role Role
}

func (dao *SQLDao) CreateUserTx(ctx context.Context, parm *CreateUserTxParms) (CreateUserTxResults, error) {
	var result CreateUserTxResults
	//建立user 未commit之前  需要滿足table 外件約束 所以使用ReadUncommitted
	err := dao.execTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadUncommitted,
		AccessMode: pgx.ReadWrite,
	}, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, *parm.Arg)
		if err != nil {
			return err
		}

		role, err := q.GetRoleByRoleName(ctx, constants.DEFAULT_USER_ROLE)
		if err != nil {
			if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
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

		//after create
		if parm.AfterCreate != nil {
			return parm.AfterCreate(result.User)
		}
		return nil
	})
	return result, err
}

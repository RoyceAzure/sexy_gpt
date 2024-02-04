package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParams struct {
	ID          int64
	SecretCode  string
	IsUsed      bool
	IsValidated bool
}

type VerifyEmailTxResults struct {
	User        User
	VerifyEmail VertifyEmail
}

/*
verify email isUsed改為true，
isValidate 改為false，
User isVerfified 改為true，
*/
func (dao *SQLDao) UpdateVerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResults, error) {
	var result VerifyEmailTxResults
	err := dao.execTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}, func(q *Queries) error {
		var err error
		varifyEmail, err := q.UpdateVertifyEmail(ctx, UpdateVertifyEmailParams{
			ID: arg.ID,
			IsUsed: pgtype.Bool{
				Bool:  arg.IsUsed,
				Valid: true,
			},
			IsValidated: pgtype.Bool{
				Bool:  arg.IsValidated,
				Valid: true,
			},
			UsedDate: pgtype.Timestamptz{
				Time:  time.Now().UTC(),
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			UserID: varifyEmail.UserID,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		return err
	})

	return result, err
}

package db

import (
	"context"
	"fmt"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dao interface {
	Querier
	CreateUserTx(context.Context, *CreateUserTxParms) (CreateUserTxResults, error)
	UpdateVerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResults, error)
}

type SQLDao struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewSQLDao(connPool *pgxpool.Pool) Dao {
	return &SQLDao{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

/*
Transation
*/
func (dao *SQLDao) execTx(ctx context.Context, options pgx.TxOptions, fn func(*Queries) error) error {
	tx, err := dao.connPool.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("%s, %w", rbErr.Error(), gpt_error.ErrInternal)
		}
		return err
	}

	return tx.Commit(ctx)
}

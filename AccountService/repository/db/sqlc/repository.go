package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dao interface {
	Querier
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

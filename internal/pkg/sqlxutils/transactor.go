package sqlxutils

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	WithTx(ctx context.Context, txName string, tFunc func(ctx context.Context) error) error
	WithConfiguredTx(ctx context.Context, txName string, tFunc func(ctx context.Context) error, opts *sql.TxOptions) error
}

type SqlxTransactor struct {
	db *sqlx.DB
}

func NewSqlxTransactor(db *sqlx.DB) *SqlxTransactor {
	return &SqlxTransactor{db: db}
}

func (r *SqlxTransactor) WithTx(ctx context.Context, txName string, tFunc func(ctx context.Context) error) error {
	return WithTxx(ctx, txName, r.db, tFunc, nil)
}

func (r *SqlxTransactor) WithConfiguredTx(ctx context.Context, txName string, tFunc func(ctx context.Context) error, opts *sql.TxOptions) error {
	return WithTxx(ctx, txName, r.db, tFunc, opts)
}

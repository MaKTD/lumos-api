package sqlxutils

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Querying interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SqlxQuerying interface {
	Querying
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

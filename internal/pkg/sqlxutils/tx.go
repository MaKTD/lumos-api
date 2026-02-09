package sqlxutils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ContextTrxKey struct{}

func InjectTxx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, ContextTrxKey{}, tx)
}

func ExtractTxx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(ContextTrxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

func WithTxx(
	ctx context.Context,
	txName string,
	db *sqlx.DB,
	tFunc func(ctx context.Context) error,
	opts *sql.TxOptions,
) (err error) {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to start transaction with name = %s: %w", txName, err)
	}

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			if err != nil {
				err = fmt.Errorf(
					"failed to rollback transaction with name = %s: trx execution error %w, rollback error %w",
					txName,
					err,
					rollbackErr,
				)
			} else {
				err = fmt.Errorf("failed to rollback transaction with name = %s: %w", txName, rollbackErr)
			}
		}
	}()

	err = tFunc(InjectTxx(ctx, tx))
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction with name = %s: %w", txName, err)
	}

	return nil
}

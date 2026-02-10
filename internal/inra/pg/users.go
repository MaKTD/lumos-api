package pg

import (
	"context"
	"database/sql"
	"doctormakarhina/lumos/internal/core/domain"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
  SELECT id, email, name, tariff, expires_at, subscription_id, subscription_status, last_sub_price
  FROM lumos.users
  WHERE email = $1
  LIMIT 1
 `

	var user domain.User
	err := r.db.GetContext(ctx, &user, q, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	const q = `
  INSERT INTO lumos.users (id, email, name, tariff, expires_at, subscription_id, subscription_status, last_sub_price)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  RETURNING *
 `

	var created domain.User
	err := r.db.GetContext(
		ctx,
		&created,
		q,
		user.ID,
		user.Email,
		user.Name,
		user.Tariff,
		user.ExpiresAt,
		user.SubscriptionID,
		user.SubscriptionStatus,
		user.LastSubPrice,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *UserRepo) FindByEmailOrCreate(ctx context.Context, user domain.User) (*domain.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	// Concurrency safety:
	// serialize "find-or-create" operations per email within this transaction.
	// This avoids the race where two concurrent requests both don't find a row and both try to insert.
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, user.Email); err != nil {
		return nil, err
	}

	const qFind = `
  SELECT id, email, name, tariff, expires_at, subscription_id, subscription_status, last_sub_price
  FROM lumos.users
  WHERE email = $1
  LIMIT 1
 `
	var existing domain.User
	err = tx.GetContext(ctx, &existing, qFind, user.Email)
	if err == nil {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &existing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	const qInsert = `
  INSERT INTO lumos.users (id, email, name, tariff, expires_at, subscription_id, subscription_status, last_sub_price)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  RETURNING *
 `
	var created domain.User
	err = tx.GetContext(
		ctx,
		&created,
		qInsert,
		user.ID,
		user.Email,
		user.Name,
		user.Tariff,
		user.ExpiresAt,
		user.SubscriptionID,
		user.SubscriptionStatus,
		user.LastSubPrice,
	)
	if err != nil {
		// Conflicts (e.g., UNIQUE(email)) are treated as errors and returned as-is.
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *UserRepo) UpdateSub(ctx context.Context, user domain.User) (*domain.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	// Concurrency safety:
	// serialize subscription updates per user within this transaction to avoid concurrent updates racing.
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, user.Email); err != nil {
		return nil, err
	}

	const q = `
		UPDATE lumos.users
		SET
      tariff = $1,
			last_sub_price = $2,
			expires_at = $3,
			subscription_id = $4,
			subscription_status = $5
		WHERE id = $6
		RETURNING *
	`

	var updated domain.User
	err = tx.GetContext(
		ctx,
		&updated,
		q,
		user.Tariff,
		user.LastSubPrice,
		user.ExpiresAt,
		user.SubscriptionID,
		user.SubscriptionStatus,
		user.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %s, do not found for sub update", user.ID)
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *UserRepo) UpdateSubStatusBySubID(ctx context.Context, subscriptionID string, status string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Concurrency safety:
	// serialize status updates per subscription within this transaction to avoid concurrent updates racing.
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, subscriptionID); err != nil {
		return err
	}

	const q = `
		UPDATE lumos.users
		SET subscription_status = $1
		WHERE subscription_id = $2
	`

	res, err := tx.ExecContext(ctx, q, status, subscriptionID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user with provided subscription_id do not found for status update")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

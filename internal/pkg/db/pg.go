package db

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPG(
	ctx context.Context,
	dsn string,
	maxConns int,
	maxIdleConns int,
	maxConnIdleTime time.Duration,
) (*sqlx.DB, error) {

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxConnIdleTime)

	return db, nil
}

package postgresql

import (
	"context"
	"database/sql"
	"doctormakarhina/lumos/internal/pkg/errs"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"time"
)

type Config struct {
	Url               string
	MaxConns          int32
	MinConns          int32
	HealthCheckPeriod time.Duration
	MaxConnIdleTime   time.Duration
	Debug             bool
	ConnectionTimeout time.Duration
}

func NewPgPool(ctx context.Context, config Config, logger *slog.Logger) (*sql.DB, error) {
	poolConf, err := pgxpool.ParseConfig(config.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pg connection string: %w", err)
	}
	poolLogger := logger.With(
		slog.String("context", "pgPool"),
		slog.String("host", poolConf.ConnConfig.Host),
		slog.Int("port", int(poolConf.ConnConfig.Port)),
		slog.String("database", poolConf.ConnConfig.Database),
	)
	poolConf.MaxConns = config.MaxConns
	poolConf.MinConns = config.MinConns
	poolConf.HealthCheckPeriod = config.HealthCheckPeriod
	poolConf.MaxConnIdleTime = config.MaxConnIdleTime
	poolConf.ConnConfig.ConnectTimeout = config.ConnectionTimeout
	if config.Debug {
		poolConf.BeforeConnect = func(ctx context.Context, config *pgx.ConnConfig) error {
			poolLogger.Debug("new connection will be made")
			return nil
		}
		poolConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			poolLogger.Debug("new connection was made")
			return nil
		}
		poolConf.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
			poolLogger.Debug("connection will be acquired")
			return true
		}
		poolConf.AfterRelease = func(conn *pgx.Conn) bool {
			poolLogger.Debug("connection was released")
			return true
		}
	}

	pgxPool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.ErrCodeInternal, "failed to init pgx pool")
	}
	pool := stdlib.OpenDBFromPool(pgxPool)
	return pool, nil
}

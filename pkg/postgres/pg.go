package postgres

import (
	"context"
	"time"

	"github.com/VoevodinAnton/metrics/db"
	"github.com/VoevodinAnton/metrics/pkg/config"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	maxOpenConns      = 50
	connMaxLifetime   = 3
	healthCheckPeriod = 1 * time.Minute
	maxConnIdleTime   = 1 * time.Minute
	minConns          = 10
)

func NewPgxConn(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, error) {
	if err := db.RunMigrations(cfg.DatabaseDSN); err != nil {
		return nil, errors.Wrap(err, "runMigrations")
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ParseConfig")
	}

	poolCfg.MaxConns = int32(maxOpenConns)
	poolCfg.HealthCheckPeriod = healthCheckPeriod
	poolCfg.MaxConnIdleTime = maxConnIdleTime
	poolCfg.MaxConnLifetime = time.Duration(connMaxLifetime) * time.Minute
	poolCfg.MinConns = minConns

	connPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ConnectConfig")
	}

	if err := connPool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "ping db")
	}

	return connPool, nil
}

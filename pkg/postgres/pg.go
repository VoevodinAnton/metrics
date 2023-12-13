package postgres

import (
	"context"
	"time"

	"github.com/VoevodinAnton/metrics/pkg/config"
	"github.com/jackc/pgx/v4/pgxpool"
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
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ParseConfig")
	}

	poolCfg.MaxConns = int32(maxOpenConns)
	poolCfg.HealthCheckPeriod = healthCheckPeriod
	poolCfg.MaxConnIdleTime = maxConnIdleTime
	poolCfg.MaxConnLifetime = time.Duration(connMaxLifetime) * time.Minute
	poolCfg.MinConns = minConns

	connPoll, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ConnectConfig")
	}

	return connPoll, nil
}

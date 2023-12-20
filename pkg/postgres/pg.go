package postgres

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/VoevodinAnton/metrics/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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
	if err := runMigrations(cfg.DatabaseDSN); err != nil {
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

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

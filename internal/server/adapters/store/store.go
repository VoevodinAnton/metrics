package store

import (
	"context"

	"github.com/VoevodinAnton/metrics/internal/server/adapters/store/memory"
	pg_store "github.com/VoevodinAnton/metrics/internal/server/adapters/store/postgres"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/VoevodinAnton/metrics/pkg/postgres"
	"github.com/pkg/errors"
)

type Store interface {
	GetGauge(ctx context.Context, name string) (models.Metric, error)
	GetCounter(ctx context.Context, name string) (models.Metric, error)
	PutCounter(ctx context.Context, update models.Metric) error
	PutGauge(ctx context.Context, update models.Metric) error
	GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error)
	GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error)
	Ping(ctx context.Context) error
	Close()
}

func NewStore(cfg *config.Config) (Store, error) {
	if cfg.Postgres.DatabaseDSN != "" {
		db, err := postgres.NewPgxConn(context.Background(), cfg.Postgres)
		if err != nil {
			return nil, errors.Wrap(err, "postgres.NewPgxConn")
		}
		return pg_store.NewStore(db), nil
	} else {
		return memory.NewStorage(), nil
	}
}

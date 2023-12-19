package postgres

import (
	"context"
	"fmt"

	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetGauge(ctx context.Context, name string) (models.Metric, error) {
	const query = `SELECT name, value FROM gauge_metrics WHERE name = $1
		ORDER BY updated_at DESC LIMIT 1;`
	row := s.db.QueryRow(ctx, query, name)
	var metric = models.Metric{
		Type: models.Gauge,
	}
	err := row.Scan(&metric.Name, &metric.Value)
	if err != nil {
		return models.Metric{}, errors.Wrap(err, "row.Scan gauge")
	}
	fmt.Println(metric)
	return metric, nil
}

func (s *Store) GetCounter(ctx context.Context, name string) (models.Metric, error) {
	const query = `SELECT name, value FROM counter_metrics WHERE name = $1 
		ORDER BY updated_at DESC LIMIT 1;`
	row := s.db.QueryRow(ctx, query, name)
	var metric = models.Metric{
		Type: models.Counter,
	}
	err := row.Scan(&metric.Name, &metric.Value)
	if err != nil {
		return models.Metric{}, errors.Wrap(err, "row.Scan counter")
	}

	return metric, nil
}

func (s *Store) PutCounter(ctx context.Context, update models.Metric) error {
	zap.L().Debug("store.postgres.putCounter", zap.Reflect("counterMetricPut", update))
	const query = `INSERT INTO counter_metrics (name, value) VALUES ($1, $2);`
	_, err := s.db.Exec(ctx, query, update.Name, update.Value)
	if err != nil {
		return errors.Wrap(err, "db.Exec counter")
	}

	return nil
}

func (s *Store) PutGauge(ctx context.Context, update models.Metric) error {
	zap.L().Debug("store.postgres.putGauge", zap.Reflect("gaugeMetricPut", update))
	const query = `INSERT INTO gauge_metrics (name, value) VALUES ($1, $2);`
	_, err := s.db.Exec(ctx, query, update.Name, update.Value)
	if err != nil {
		return errors.Wrap(err, "db.Exec gauge")
	}

	return nil
}

func (s *Store) GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error) {
	const query = `SELECT name, value FROM counter_metrics cm1 WHERE updated_at  = (
		SELECT MAX(updated_at)
		FROM counter_metrics cm2
		WHERE cm2.name = cm1.name
	);`

	return s.getMetrics(ctx, query)
}

func (s *Store) GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error) {
	const query = `SELECT name, value FROM gauge_metrics gm1 WHERE updated_at  = (
		SELECT MAX(updated_at)
		FROM gauge_metrics gm2
		WHERE gm2.name = gm1.name
	);`

	return s.getMetrics(ctx, query)
}

func (s *Store) getMetrics(ctx context.Context, query string) (map[string]models.Metric, error) {
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "db.Query gauge")
	}
	defer rows.Close()

	metrics := make(map[string]models.Metric, 0)
	for rows.Next() {
		var metric models.Metric
		if err := rows.Scan(&metric.Name, &metric.Value); err != nil {
			return nil, errors.Wrap(err, "rows.Scan geuge")
		}
		metrics[metric.Name] = metric
	}

	return metrics, nil
}

func (s *Store) Ping(ctx context.Context) error {
	return errors.Wrap(s.db.Ping(ctx), "db.Ping")
}

func (s *Store) Close() {
	s.db.Close()
}

package postgres

import (
	"context"

	"github.com/VoevodinAnton/metrics/internal/server/models"

	"github.com/jackc/pgtype"
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

func (s *Store) GetGaugeMetric(ctx context.Context, name string) (models.Metric, error) {
	row := s.db.QueryRow(ctx, getGaugeMetricQuery, name)
	var metric = models.Metric{
		Type: models.Gauge,
	}
	err := row.Scan(&metric.Name, &metric.Value)
	if err != nil {
		return models.Metric{}, errors.Wrap(err, "row.Scan gauge")
	}

	return metric, nil
}

func (s *Store) GetCounterMetric(ctx context.Context, name string) (models.Metric, error) {
	row := s.db.QueryRow(ctx, getCounterMetricQuery, name)
	var metric = models.Metric{
		Type: models.Counter,
	}
	var value pgtype.Numeric
	err := row.Scan(&metric.Name, &value)
	if err != nil {
		return models.Metric{}, errors.Wrap(err, "row.Scan counter")
	}
	metric.Value = value.Int.Int64()

	return metric, nil
}

func (s *Store) PutCounterMetric(ctx context.Context, update models.Metric) error {
	zap.L().Debug("store.postgres.putCounterMetric", zap.Reflect("counterMetricPut", update))
	_, err := s.db.Exec(ctx, insertCounterMetricQuery, update.Name, update.Value)
	if err != nil {
		return errors.Wrap(err, "db.Exec counter")
	}

	return nil
}

func (s *Store) PutCounterMetrics(ctx context.Context, updates []models.Metric) error {
	zap.L().Debug("store.postgres.putCounterMetrics", zap.Reflect("counterMetricsPut", updates))

	return s.putMetrics(ctx, insertCounterMetricQueryName, insertCounterMetricQuery, updates)
}

func (s *Store) PutGaugeMetric(ctx context.Context, update models.Metric) error {
	zap.L().Debug("store.postgres.putGaugeMetric", zap.Reflect("gaugeMetricPut", update))
	_, err := s.db.Exec(ctx, insertGaugeMetricQuery, update.Name, update.Value)
	if err != nil {
		return errors.Wrap(err, "db.Exec gauge")
	}

	return nil
}

func (s *Store) PutGaugeMetrics(ctx context.Context, updates []models.Metric) error {
	zap.L().Debug("store.postgres.putGaugeMetrics", zap.Reflect("gaugeMetricsPut", updates))

	return s.putMetrics(ctx, insertGaugeMetricQueryName, insertGaugeMetricQuery, updates)
}

func (s *Store) putMetrics(ctx context.Context, queryName, query string, updates []models.Metric) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Begin")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	_, err = tx.Prepare(ctx, queryName, query)
	if err != nil {
		return errors.Wrap(err, "tx.Prepare")
	}
	for _, update := range updates {
		_, err := tx.Exec(ctx, queryName, update.Name, update.Value)
		if err != nil {
			return errors.Wrap(err, "tx.Exec")
		}
	}

	return tx.Commit(ctx) //nolint: wrapcheck //unnecessary
}

func (s *Store) GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error) {
	return s.getMetrics(ctx, getCounterMetricsQuery)
}

func (s *Store) GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error) {
	return s.getMetrics(ctx, getGaugeMetricsQuery)
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

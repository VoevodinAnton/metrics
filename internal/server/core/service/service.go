package service

import (
	"context"

	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/pkg/errors"
)

type Store interface {
	PutGauge(ctx context.Context, metric models.Metric) error
	GetGauge(ctx context.Context, name string) (models.Metric, error)
	PutCounter(ctx context.Context, metric models.Metric) error
	GetCounter(ctx context.Context, name string) (models.Metric, error)
	GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error)
	GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error)
	Ping(ctx context.Context) error
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetMetric(ctx context.Context, metric *domain.Metrics) (*domain.Metrics, error) {
	var metricResp models.Metric
	var err error
	switch metric.MType {
	case models.Gauge:
		metricResp, err = s.store.GetGauge(ctx, metric.ID)
		if err != nil {
			return nil, errors.Wrap(err, "getGauge")
		}
	case models.Counter:
		metricResp, err = s.store.GetCounter(ctx, metric.ID)
		if err != nil {
			return nil, errors.Wrap(err, "getCounter")
		}
	}

	return metricToResponse(metricResp), nil
}

func (s *Service) UpdateMetric(ctx context.Context, metric *domain.Metrics) error {
	metricUpdate := requestToMetric(metric)
	switch metric.MType {
	case models.Gauge:
		err := s.store.PutGauge(ctx, metricUpdate)
		if err != nil {
			return errors.Wrap(err, "updateGauge")
		}
	case models.Counter:
		err := s.store.PutCounter(ctx, metricUpdate)
		if err != nil {
			return errors.Wrap(err, "updateCounter")
		}
	}

	return nil
}

func (s *Service) GetMetrics(ctx context.Context) (*[]domain.Metrics, error) {
	counterMetrics, err := s.store.GetCounterMetrics(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getCounterMetrics")
	}
	gaugeMetrics, err := s.store.GetGaugeMetrics(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getGaugeMetrics")
	}
	resp := make([]domain.Metrics, 0, len(counterMetrics)+len(gaugeMetrics))
	for _, v := range counterMetrics {
		resp = append(resp, *metricToResponse(v))
	}
	for _, v := range gaugeMetrics {
		resp = append(resp, *metricToResponse(v))
	}

	return &resp, nil
}

func (s *Service) Ping(ctx context.Context) error {
	return errors.Wrap(s.store.Ping(ctx), "ping")
}

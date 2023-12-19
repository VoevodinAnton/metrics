package service

import (
	"context"

	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/pkg/errors"
)

type Store interface {
	GetCounterMetric(ctx context.Context, name string) (models.Metric, error)
	GetGaugeMetric(ctx context.Context, name string) (models.Metric, error)
	GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error)
	GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error)
	PutCounterMetric(ctx context.Context, metric models.Metric) error
	PutGaugeMetric(ctx context.Context, metric models.Metric) error
	PutCounterMetrics(ctx context.Context, updates []models.Metric) error
	PutGaugeMetrics(ctx context.Context, updates []models.Metric) error

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
		metricResp, err = s.store.GetGaugeMetric(ctx, metric.ID)
		if err != nil {
			return nil, errors.Wrap(err, "getGauge")
		}
	case models.Counter:
		metricResp, err = s.store.GetCounterMetric(ctx, metric.ID)
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
		err := s.store.PutGaugeMetric(ctx, metricUpdate)
		if err != nil {
			return errors.Wrap(err, "putGaugeMetric")
		}
	case models.Counter:
		err := s.store.PutCounterMetric(ctx, metricUpdate)
		if err != nil {
			return errors.Wrap(err, "putCounterMetric")
		}
	}
	return s.updateMetric(ctx, metric)
}

func (s *Service) UpdatesMetrics(ctx context.Context, metrics *[]domain.Metrics) error {
	metricsModel := requestToMetrics(metrics)
	gaugeMetrics := make([]models.Metric, 0)
	counterMetrics := make([]models.Metric, 0)
	for _, metric := range metricsModel {
		switch metric.Type {
		case models.Gauge:
			gaugeMetrics = append(gaugeMetrics, metric)
		case models.Counter:
			counterMetrics = append(counterMetrics, metric)
		}
	}
	if len(counterMetrics) != 0 {
		if err := s.store.PutCounterMetrics(ctx, counterMetrics); err != nil {
			return errors.Wrap(err, "store.PutCounterMetrics")
		}
	}
	if len(gaugeMetrics) != 0 {
		if err := s.store.PutGaugeMetrics(ctx, gaugeMetrics); err != nil {
			return errors.Wrap(err, "store.PutGaugeMetrics")
		}
	}

	return nil
}

func (s *Service) updateMetric(ctx context.Context, metric *domain.Metrics) error {
	metricUpdate := requestToMetric(metric)
	switch metric.MType {
	case models.Gauge:
		err := s.store.PutGaugeMetric(ctx, metricUpdate)
		if err != nil {
			return errors.Wrap(err, "updateGauge")
		}
	case models.Counter:
		err := s.store.PutCounterMetric(ctx, metricUpdate)
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

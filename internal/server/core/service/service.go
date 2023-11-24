package service

import (
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/pkg/errors"
)

type Store interface {
	UpdateGauge(metric *models.Metric) error
	GetGauge(name string) (*models.Metric, error)
	UpdateCounter(metric *models.Metric) error
	GetCounter(name string) (*models.Metric, error)
	GetCounterMetrics() (map[string]models.Metric, error)
	GetGaugeMetrics() (map[string]models.Metric, error)
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetMetric(metric *domain.Metrics) (*domain.Metrics, error) {
	var metricResp *models.Metric
	var err error
	switch metric.MType {
	case string(models.Gauge):
		metricResp, err = s.store.GetGauge(metric.ID)
		if err != nil {
			return nil, errors.Wrap(err, "GetGauge")
		}
	case string(models.Counter):
		metricResp, err = s.store.GetCounter(metric.ID)
		if err != nil {
			return nil, errors.Wrap(err, "GetCounter")
		}
	}

	return metricToResponse(metricResp), nil
}

func (s *Service) UpdateMetric(metric *domain.Metrics) error {
	metricUpdate := requestToMetric(metric)
	switch metric.MType {
	case string(models.Gauge):
		err := s.store.UpdateGauge(metricUpdate)
		if err != nil {
			return errors.Wrap(err, "UpdateGauge")
		}
	case string(models.Counter):
		err := s.store.UpdateCounter(metricUpdate)
		if err != nil {
			return errors.Wrap(err, "UpdateCounter")
		}
	}

	return nil
}

func (s *Service) GetMetrics() (*[]domain.Metrics, error) {
	counterMetrics, err := s.store.GetCounterMetrics()
	if err != nil {
		return nil, errors.Wrap(err, "GetCounterMetrics")
	}
	gaugeMetrics, err := s.store.GetGaugeMetrics()
	if err != nil {
		return nil, errors.Wrap(err, "GetGaugeMetrics")
	}
	resp := make([]domain.Metrics, 0, len(counterMetrics)+len(gaugeMetrics))
	for _, v := range counterMetrics {
		resp = append(resp, *metricToResponse(&v))
	}
	for _, v := range gaugeMetrics {
		resp = append(resp, *metricToResponse(&v))
	}

	return &resp, nil
}

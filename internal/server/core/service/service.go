package service

import (
	"github.com/VoevodinAnton/metrics/internal/models"
	"github.com/pkg/errors"
)

type Store interface {
	UpdateGauge(metric models.Metric) error
	GetGauge(name string) (models.Metric, error)
	UpdateCounter(metric models.Metric) error
	GetCounter(name string) (models.Metric, error)
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

func (s *Service) GetMetric(req models.Metric) (models.Metric, error) {
	var metricResp models.Metric
	switch req.Type {
	case models.Gauge:
		gauge, err := s.store.GetGauge(req.Name)
		if err != nil {
			return models.Metric{}, errors.Wrap(err, "GetGauge")
		}
		metricResp = gauge
	case models.Counter:
		counter, err := s.store.GetCounter(req.Name)
		if err != nil {
			return models.Metric{}, errors.Wrap(err, "GetCounter")
		}
		metricResp = counter
	}

	return metricResp, nil
}

func (s *Service) UpdateMetric(req models.Metric) error {
	switch req.Type {
	case models.Gauge:
		err := s.store.UpdateGauge(req)
		if err != nil {
			return errors.Wrap(err, "UpdateGauge")
		}
	case models.Counter:
		err := s.store.UpdateCounter(req)
		if err != nil {
			return errors.Wrap(err, "UpdateCounter")
		}
	}

	return nil
}

func (s *Service) GetMetrics() ([]models.Metric, error) {
	counterMetrics, err := s.store.GetCounterMetrics()
	if err != nil {
		return nil, errors.Wrap(err, "GetCounterMetrics")
	}
	gaugeMetrics, err := s.store.GetGaugeMetrics()
	if err != nil {
		return nil, errors.Wrap(err, "GetGaugeMetrics")
	}
	resp := make([]models.Metric, 0, len(counterMetrics)+len(gaugeMetrics))
	for _, v := range counterMetrics {
		resp = append(resp, v)
	}
	for _, v := range gaugeMetrics {
		resp = append(resp, v)
	}

	return resp, nil
}

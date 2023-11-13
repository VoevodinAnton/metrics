package service

import (
	"github.com/VoevodinAnton/metrics/internal/models"
)

type Store interface {
	UpdateGauge(metric *models.GaugeMetric) error
	GetGauge(name string) (*models.GaugeMetric, error)
	UpdateCounter(metric *models.CounterMetric) error
	GetCounter(name string) (*models.CounterMetric, error)
	GetCounterMetrics() (map[string]*models.CounterMetric, error)
	GetGaugeMetrics() (map[string]*models.GaugeMetric, error)
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetMetric(req *models.MetricReq) (*models.MetricResp, error) {
	var metricResp *models.MetricResp
	switch req.Type {
	case models.Gauge:
		gauge, err := s.store.GetGauge(req.Name)
		if err != nil {
			return nil, err // nolint: wrapheck
		}
		metricResp = gaugeModelToAPI(gauge)
	case models.Counter:
		counter, err := s.store.GetCounter(req.Name)
		if err != nil {
			return nil, err // nolint: wrapheck
		}
		metricResp = counterModelToAPI(counter)
	}

	return metricResp, nil
}

func (s *Service) UpdateMetric(req *models.MetricReq) error {
	switch req.Type {
	case models.Gauge:
		err := s.store.UpdateGauge(gaugeAPIToModel(req))
		if err != nil {
			return err // nolint: wrapheck
		}
	case models.Counter:
		err := s.store.UpdateCounter(counterAPIToModel(req))
		if err != nil {
			return err // nolint: wrapheck
		}
	}

	return nil
}

func (s *Service) GetMetrics() ([]*models.MetricResp, error) {
	var resp []*models.MetricResp
	counterMetrics, err := s.store.GetCounterMetrics()
	if err != nil {
		return nil, err // nolint: wrapheck
	}
	gaugeMetrics, err := s.store.GetGaugeMetrics()
	if err != nil {
		return nil, err // nolint: wrapheck
	}
	for _, v := range counterMetrics {
		resp = append(resp, counterModelToAPI(v))
	}
	for _, v := range gaugeMetrics {
		resp = append(resp, gaugeModelToAPI(v))
	}

	return resp, nil
}

func counterAPIToModel(c *models.MetricReq) *models.CounterMetric {
	return &models.CounterMetric{
		Name:  c.Name,
		Type:  c.Type,
		Value: int64(c.Value),
	}
}

func counterModelToAPI(m *models.CounterMetric) *models.MetricResp {
	return &models.MetricResp{
		Name:  m.Name,
		Type:  m.Type,
		Value: float64(m.Value),
	}
}

func gaugeAPIToModel(c *models.MetricReq) *models.GaugeMetric {
	return &models.GaugeMetric{
		Name:  c.Name,
		Type:  c.Type,
		Value: c.Value,
	}
}

func gaugeModelToAPI(m *models.GaugeMetric) *models.MetricResp {
	return &models.MetricResp{
		Name:  m.Name,
		Type:  m.Type,
		Value: float64(m.Value),
	}
}

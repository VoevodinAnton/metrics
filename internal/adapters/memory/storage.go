package memory

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/VoevodinAnton/metrics/internal/models"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

type Storage struct {
	gaugeMetrics   map[string]*models.GaugeMetric
	counterMetrics map[string]*models.CounterMetric

	sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		gaugeMetrics:   make(map[string]*models.GaugeMetric),
		counterMetrics: make(map[string]*models.CounterMetric),
	}
}

func (s *Storage) UpdateGauge(metric *models.GaugeMetric) error {
	s.Lock()
	defer s.Unlock()
	s.gaugeMetrics[metric.Name] = metric

	return nil
}

func (s *Storage) GetGauge(name string) (*models.GaugeMetric, error) {
	s.Lock()
	defer s.Unlock()
	value, ok := s.gaugeMetrics[name]
	if !ok {
		return nil, errors.Wrap(ErrMetricNotFound, name)
	}

	return value, nil
}

func (s *Storage) UpdateCounter(metric *models.CounterMetric) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.counterMetrics[metric.Name]
	if !ok {
		s.counterMetrics[metric.Name] = metric
	} else {
		s.counterMetrics[metric.Name].Value += metric.Value
	}

	return nil
}

func (s *Storage) GetCounter(name string) (*models.CounterMetric, error) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.counterMetrics[name]
	if !ok {
		return nil, errors.Wrap(ErrMetricNotFound, name)
	}

	return value, nil
}

func (s *Storage) GetCounterMetrics() (map[string]*models.CounterMetric, error) {
	s.Lock()
	defer s.Unlock()
	return s.counterMetrics, nil
}

func (s *Storage) GetGaugeMetrics() (map[string]*models.GaugeMetric, error) {
	s.Lock()
	defer s.Unlock()
	return s.gaugeMetrics, nil
}

func (s *Storage) ResetCounter(name string) error {
	s.Lock()
	defer s.Unlock()

	s.counterMetrics[name] = &models.CounterMetric{}

	return nil
}

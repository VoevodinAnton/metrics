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
	gaugeMetrics   map[string]*models.Metric
	counterMetrics map[string]*models.Metric

	sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		gaugeMetrics:   make(map[string]*models.Metric),
		counterMetrics: make(map[string]*models.Metric),
	}
}

func (s *Storage) UpdateGauge(metric *models.Metric) error {
	s.Lock()
	defer s.Unlock()
	s.gaugeMetrics[metric.Name] = metric

	return nil
}

func (s *Storage) GetGauge(name string) (*models.Metric, error) {
	s.Lock()
	defer s.Unlock()
	value, ok := s.gaugeMetrics[name]
	if !ok {
		return nil, errors.Wrap(ErrMetricNotFound, name)
	}

	return value, nil
}

func (s *Storage) UpdateCounter(update *models.Metric) error {
	s.Lock()
	defer s.Unlock()
	metric, ok := s.counterMetrics[update.Name]
	if !ok {
		metric = update
	} else {
		value, _ := metric.Value.(int64)
		if newValue, ok := update.Value.(int64); ok {
			metric.Value = value + newValue
		}
	}

	s.counterMetrics[update.Name] = metric

	return nil
}

func (s *Storage) GetCounter(name string) (*models.Metric, error) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.counterMetrics[name]
	if !ok {
		return nil, errors.Wrap(ErrMetricNotFound, name)
	}

	return value, nil
}

func (s *Storage) GetCounterMetrics() (map[string]*models.Metric, error) {
	s.Lock()
	defer s.Unlock()
	return s.counterMetrics, nil
}

func (s *Storage) GetGaugeMetrics() (map[string]*models.Metric, error) {
	s.Lock()
	defer s.Unlock()
	return s.gaugeMetrics, nil
}

func (s *Storage) ResetCounter(name string) error {
	s.Lock()
	defer s.Unlock()

	s.counterMetrics[name] = &models.Metric{}

	return nil
}

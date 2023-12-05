package memory

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/VoevodinAnton/metrics/internal/server/models"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

type Storage struct {
	gaugeMetrics   sync.Map
	counterMetrics sync.Map
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) UpdateGauge(metric *models.Metric) error {
	s.gaugeMetrics.Store(metric.Name, metric)
	return nil
}

func (s *Storage) GetGauge(name string) (*models.Metric, error) {
	value, ok := s.gaugeMetrics.Load(name)
	if !ok {
		return nil, errors.Wrap(ErrMetricNotFound, name)
	}

	return value.(*models.Metric), nil
}

func (s *Storage) UpdateCounter(update *models.Metric) error {
	m, ok := s.counterMetrics.Load(update.Name)
	if !ok {
		s.counterMetrics.Store(update.Name, update)
		return nil
	}
	metric, _ := m.(*models.Metric)
	value, _ := metric.Value.(*int64)
	if newValue, ok := update.Value.(*int64); ok {
		sum := *value + *newValue
		metric.Value = &sum
	} else {
		fmt.Println(reflect.TypeOf(update.Value))
	}
	s.counterMetrics.Store(update.Name, metric)

	return nil
}

func (s *Storage) GetCounter(name string) (*models.Metric, error) {
	value, ok := s.counterMetrics.Load(name)
	if !ok {
		return nil, errors.Wrap(ErrMetricNotFound, name)
	}

	return value.(*models.Metric), nil
}

func (s *Storage) GetCounterMetrics() (map[string]models.Metric, error) {
	data := make(map[string]models.Metric)
	s.counterMetrics.Range(func(key, value any) bool {
		keyStr, _ := key.(string)
		valueMetric, _ := value.(*models.Metric)
		data[keyStr] = *valueMetric
		return true
	})

	return data, nil
}

func (s *Storage) GetGaugeMetrics() (map[string]models.Metric, error) {
	data := make(map[string]models.Metric)
	s.gaugeMetrics.Range(func(key, value any) bool {
		keyStr, _ := key.(string)
		valueMetric, _ := value.(*models.Metric)

		data[keyStr] = *valueMetric
		return true
	})

	return data, nil
}

func (s *Storage) ResetCounter(name string) error {
	metric := &models.Metric{
		Name:  name,
		Value: int64(0),
		Type:  models.Counter,
	}

	s.counterMetrics.Store(name, metric)

	return nil
}

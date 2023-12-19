package memory

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/VoevodinAnton/metrics/internal/server/models"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

type Store struct {
	gaugeMetrics   sync.Map
	counterMetrics sync.Map
	sync.Mutex
}

func NewStorage() *Store {
	return &Store{}
}

func (s *Store) GetGauge(ctx context.Context, name string) (models.Metric, error) {
	value, ok := s.gaugeMetrics.Load(name)
	if !ok {
		return models.Metric{}, errors.Wrap(ErrMetricNotFound, name)
	}

	return value.(models.Metric), nil
}

func (s *Store) GetCounter(ctx context.Context, name string) (models.Metric, error) {
	value, ok := s.counterMetrics.Load(name)
	if !ok {
		return models.Metric{}, errors.Wrap(ErrMetricNotFound, name)
	}

	return value.(models.Metric), nil
}

func (s *Store) PutCounter(ctx context.Context, update models.Metric) error {
	zap.L().Debug("store.counter.putGauge", zap.Reflect("counterMetricPut", update))
	m, ok := s.counterMetrics.Load(update.Name)
	if !ok {
		s.counterMetrics.Store(update.Name, update)
		return nil
	}
	metric, _ := m.(models.Metric)
	value, _ := metric.Value.(int64)
	if newValue, ok := update.Value.(int64); ok {
		metric.Value = value + newValue
	} else {
		return errors.New("expected int64 type")
	}
	s.counterMetrics.Store(update.Name, metric)

	return nil
}

func (s *Store) PutGauge(ctx context.Context, update models.Metric) error {
	zap.L().Debug("store.memory.putGauge", zap.Reflect("gaugeMetricPut", update))
	s.gaugeMetrics.Store(update.Name, update)
	return nil
}

func (s *Store) GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error) {
	data := make(map[string]models.Metric)
	s.counterMetrics.Range(func(key, value any) bool {
		s.Lock()
		keyStr, _ := key.(string)
		valueMetric, _ := value.(models.Metric)
		data[keyStr] = valueMetric
		s.Unlock()
		return true
	})

	return data, nil
}

func (s *Store) GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error) {
	data := make(map[string]models.Metric)
	s.gaugeMetrics.Range(func(key, value any) bool {
		s.Lock()
		keyStr, _ := key.(string)
		valueMetric, _ := value.(models.Metric)
		data[keyStr] = valueMetric
		s.Unlock()
		return true
	})

	return data, nil
}

func (s *Store) Ping(ctx context.Context) error {
	return nil
}

func (s *Store) Close() {
}

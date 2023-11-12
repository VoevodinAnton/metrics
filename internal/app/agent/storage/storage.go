package storage

import (
	"fmt"
	"sync"

	"github.com/VoevodinAnton/metrics/internal/app/agent/repository"
)

type Storage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	sync.Mutex
}

func New() *Storage {
	return &Storage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (s *Storage) UpdateGaugeMetric(name string, metricValue float64) {
	s.Lock()
	defer s.Unlock()

	s.gaugeMetrics[name] = metricValue
}

func (s *Storage) UpdateCounterMetric(name string, metricValue int64) {
	s.Lock()
	defer s.Unlock()

	s.counterMetrics[name] += metricValue
}

func (s *Storage) GetCounterMetric(name string) (int64, error) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.counterMetrics[name]
	if !ok {
		fmt.Println(value)
		return 0, repository.ErrMetricNotFound
	}
	return value, nil
}

func (s *Storage) GetGaugeMetric(name string) (float64, error) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.gaugeMetrics[name]
	if !ok {
		fmt.Println(name)
		return 0, repository.ErrMetricNotFound
	}
	return value, nil
}

func (s *Storage) ResetCounterMetric(name string) {
	s.Lock()
	defer s.Unlock()

	s.counterMetrics[name] = 0
}

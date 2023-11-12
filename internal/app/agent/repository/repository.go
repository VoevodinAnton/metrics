package repository

import (
	"errors"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

type MetricRepository interface {
	UpdateGaugeMetric(name string, metricValue float64)
	UpdateCounterMetric(name string, metricValue int64)
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	ResetCounterMetric(name string)
}

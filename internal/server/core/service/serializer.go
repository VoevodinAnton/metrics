package service

import (
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/models"
)

func requestToMetric(m *domain.Metrics) *models.Metric {
	metric := &models.Metric{
		Name: m.ID,
		Type: models.MetricType(m.MType),
	}

	switch m.MType {
	case string(models.Gauge):
		metric.Value = m.Value
	case string(models.Counter):
		metric.Value = m.Delta
	}

	return metric
}

func metricToResponse(m *models.Metric) *domain.Metrics {
	metric := &domain.Metrics{
		ID:    m.Name,
		MType: string(m.Type),
	}

	switch m.Type {
	case models.MetricType(domain.Counter):
		v, _ := m.Value.(*int64)
		metric.Delta = v
	case models.MetricType(domain.Gauge):
		v, _ := m.Value.(*float64)
		metric.Value = v
	}

	return metric
}

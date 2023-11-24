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
		metric.Delta = m.Value.(*int64)
	case models.MetricType(domain.Gauge):
		metric.Value = m.Value.(*float64)
	}

	return metric
}

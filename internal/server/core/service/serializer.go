package service

import (
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/models"
)

func requestToMetric(m *domain.Metrics) models.Metric {
	metric := models.Metric{
		Name: m.ID,
		Type: m.MType,
	}

	switch m.MType {
	case models.Gauge:
		metric.Value = *m.Value
	case models.Counter:
		metric.Value = *m.Delta
	}

	return metric
}

func metricToResponse(m models.Metric) *domain.Metrics {
	metric := &domain.Metrics{
		ID:    m.Name,
		MType: m.Type,
	}

	switch m.Type {
	case domain.Counter:
		v, _ := m.Value.(int64)
		metric.Delta = &v
	case domain.Gauge:
		v, _ := m.Value.(float64)
		metric.Value = &v
	}

	return metric
}

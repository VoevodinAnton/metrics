package memory

import (
	"testing"

	"github.com/VoevodinAnton/metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage_UpdateGauge(t *testing.T) {
	type args struct {
		Metric models.Metric
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "gauge positive value",
			args: args{
				Metric: models.Metric{
					Name:  "SomeGaugeMetric",
					Type:  0,
					Value: 10.0,
				},
			},
			want: 10.0,
		},
		{
			name: "gauge negative value",
			args: args{
				Metric: models.Metric{
					Name:  "SomeGaugeMetric",
					Type:  0,
					Value: -10.0,
				},
			},
			want: -10.0,
		},
		{
			name: "gauge zero value",
			args: args{
				Metric: models.Metric{
					Name:  "SomeGaugeMetric",
					Type:  0,
					Value: 0.0,
				},
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			_ = s.UpdateGauge(tt.args.Metric)

			m, _ := s.gaugeMetrics.Load(tt.args.Metric.Name)
			metric, _ := m.(models.Metric)

			assert.Equal(t, tt.want, metric.Value)
		})
	}
}

func TestStorage_UpdateCounter(t *testing.T) {
	type args struct {
		Metric models.Metric
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "counter positive value",
			args: args{
				Metric: models.Metric{
					Name:  "SomeCounterMetric",
					Type:  1,
					Value: int64(10),
				},
			},
			want: 10,
		},
		{
			name: "counter zero value",
			args: args{
				Metric: models.Metric{
					Name:  "SomeCounterMetric",
					Type:  1,
					Value: int64(0),
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			_ = s.UpdateCounter(tt.args.Metric)

			m, _ := s.counterMetrics.Load(tt.args.Metric.Name)
			metric, _ := m.(models.Metric)

			assert.Equal(t, tt.want, metric.Value)
		})
	}
}

package memory

import (
	"testing"

	"github.com/VoevodinAnton/metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]*models.Metric
		CounterMetrics map[string]*models.Metric
	}
	type args struct {
		Metric models.Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "gauge positive value",
			fields: fields{
				GaugeMetrics:   make(map[string]*models.Metric),
				CounterMetrics: make(map[string]*models.Metric),
			},
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
			fields: fields{
				GaugeMetrics:   make(map[string]*models.Metric),
				CounterMetrics: make(map[string]*models.Metric),
			},
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
			fields: fields{
				GaugeMetrics:   make(map[string]*models.Metric),
				CounterMetrics: make(map[string]*models.Metric),
			},
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
			s := &Storage{
				gaugeMetrics:   tt.fields.GaugeMetrics,
				counterMetrics: tt.fields.CounterMetrics,
			}
			_ = s.UpdateGauge(tt.args.Metric)

			assert.Equal(t, tt.want, s.gaugeMetrics[tt.args.Metric.Name].Value)
		})
	}
}

func TestStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]*models.Metric
		CounterMetrics map[string]*models.Metric
	}
	type args struct {
		Metric models.Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "counter positive value",
			fields: fields{
				GaugeMetrics:   make(map[string]*models.Metric),
				CounterMetrics: make(map[string]*models.Metric),
			},
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
			fields: fields{
				GaugeMetrics:   make(map[string]*models.Metric),
				CounterMetrics: make(map[string]*models.Metric),
			},
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
			s := &Storage{
				gaugeMetrics:   tt.fields.GaugeMetrics,
				counterMetrics: tt.fields.CounterMetrics,
			}
			_ = s.UpdateCounter(tt.args.Metric)

			assert.Equal(t, tt.want, s.counterMetrics[tt.args.Metric.Name].Value)
		})
	}
}

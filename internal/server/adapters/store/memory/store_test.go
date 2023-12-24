package memory

import (
	"context"
	"testing"

	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage_PutGaugeMetric(t *testing.T) {
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
					Type:  models.Gauge,
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
					Type:  models.Gauge,
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
					Type:  models.Gauge,
					Value: 0.0,
				},
			},
			want: 0.0,
		},
	}
	for _, tt := range tests { //nolint: dupl // this is test
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{}
			err := s.PutGaugeMetric(context.Background(), tt.args.Metric)
			if err != nil {
				t.Errorf("Failed update counter: %v", err)
			}

			m, found := s.gaugeMetrics.Load(tt.args.Metric.Name)
			if !found {
				t.Errorf("Metric %s not found", tt.args.Metric.Name)
			}
			metric, ok := m.(models.Metric)
			if !ok {
				t.Error("Metric type expected")
			}

			assert.Equal(t, tt.want, metric.Value)
		})
	}
}

func TestStorage_PutCounterMetric(t *testing.T) {
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
					Type:  models.Counter,
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
					Type:  models.Counter,
					Value: int64(0),
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests { //nolint: dupl // this is test
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{}
			err := s.PutCounterMetric(context.Background(), tt.args.Metric)
			if err != nil {
				t.Errorf("Failed update counter: %v", err)
			}

			m, found := s.counterMetrics.Load(tt.args.Metric.Name)
			if !found {
				t.Errorf("Metric %s not found", tt.args.Metric.Name)
			}
			metric, ok := m.(models.Metric)
			if !ok {
				t.Error("Metric type expected")
			}

			assert.Equal(t, tt.want, metric.Value)
		})
	}
}

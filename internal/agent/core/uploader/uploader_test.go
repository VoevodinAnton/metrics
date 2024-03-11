package uploader

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

type TestCollector struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func (c *TestCollector) GetGaugeMetrics() map[string]float64 {
	return c.gaugeMetrics
}

func (c *TestCollector) GetCounterMetrics() map[string]int64 {
	return c.counterMetrics
}

func (c *TestCollector) ResetCounter() {
	for k := range c.counterMetrics {
		c.counterMetrics[k] = 0
	}
}

func NewServer(t *testing.T, expectedMetrics []domain.Metrics) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Errorf("Failed read body %v", err)
		}
		defer func() {
			_ = gz.Close()
		}()
		body, err := io.ReadAll(gz)
		if err != nil {
			t.Errorf("Failed read body %v", err)
		}

		t.Log(string(body))

		var metrics []domain.Metrics
		if err := json.Unmarshal(body, &metrics); err != nil {
			t.Fatalf("json.Unmarshal: %v", err)
		}

		require.Equal(t, expectedMetrics, metrics)
	}))
}

func TestUploader_sendCounterMetrics(t *testing.T) { //nolint: dupl // this is test
	tests := []struct {
		name            string
		expectedMetrics []domain.Metrics
		sendMetric      map[string]int64
	}{
		{
			name: "test counter",
			expectedMetrics: []domain.Metrics{
				{
					ID:    "TestCounter",
					MType: domain.Counter,
					Delta: toInt64Pointer(1),
				},
			},
			sendMetric: map[string]int64{
				"TestCounter": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := NewServer(t, tt.expectedMetrics)
			defer svr.Close()

			var cfg = &config.Config{
				ServerAddress: strings.TrimPrefix(svr.URL, "http://"),
				RateLimit:     10,
			}

			var collector = &TestCollector{
				counterMetrics: tt.sendMetric,
			}

			u := NewUploader(cfg, collector)

			u.sendCounterMetrics()
			select {
			case r := <-u.results:
				if r.Err != nil {
					t.Error(r.Err)
				}
			default:
			}
		})
	}
}

func TestUploader_sendGaugeMetrics(t *testing.T) { //nolint: dupl // this is test
	tests := []struct {
		name            string
		expectedMetrics []domain.Metrics
		sendMetric      map[string]float64
	}{
		{
			name: "test gauge",
			expectedMetrics: []domain.Metrics{
				{
					ID:    "TestGauge",
					MType: domain.Gauge,
					Value: toFloat64Pointer(2.5),
				},
			},
			sendMetric: map[string]float64{
				"TestGauge": 2.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := NewServer(t, tt.expectedMetrics)
			defer svr.Close()

			var cfg = &config.Config{
				ServerAddress: strings.TrimPrefix(svr.URL, "http://"),
				RateLimit:     10,
			}

			var collector = &TestCollector{
				gaugeMetrics: tt.sendMetric,
			}

			u := NewUploader(cfg, collector)

			u.sendGaugeMetrics()
			select {
			case r := <-u.results:
				if r.Err != nil {
					t.Error(r.Err)
				}
			default:
			}
		})
	}
}

func toInt64Pointer(i int64) *int64 {
	return &i
}

func toFloat64Pointer(f float64) *float64 {
	return &f
}

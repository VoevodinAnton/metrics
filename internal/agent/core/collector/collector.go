package collector

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
)

type Collector struct {
	metrics map[string]domain.Metrics
	cfg     *config.Config
}

func NewCollector(cfg *config.Config) *Collector {
	return &Collector{
		cfg:     cfg,
		metrics: make(map[string]domain.Metrics),
	}
}

func (c *Collector) Run() {
	ticker := time.NewTicker(c.cfg.PollInterval)
	for range ticker.C {
		c.updateMetrics()
	}
}

func (c *Collector) updateMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memStatsValue := reflect.ValueOf(memStats)

	for metricName := range c.cfg.RuntimeMetrics {
		v := memStatsValue.FieldByName(metricName)
		var gaugeMetric domain.Metrics
		if v.CanUint() {
			floatValue := float64(v.Uint())
			gaugeMetric = domain.Metrics{ID: metricName, MType: domain.Gauge, Value: toFloat64(floatValue)}
		}
		if v.CanFloat() {
			gaugeMetric = domain.Metrics{ID: metricName, MType: domain.Gauge, Value: toFloat64(v.Float())}
		}

		c.metrics[metricName] = gaugeMetric
	}

	c.metrics["RandomValue"] = domain.Metrics{ID: "RandomValue", MType: domain.Gauge, Value: getRandomValue()}
	c.metrics["PollCount"] = domain.Metrics{ID: "PollCount", MType: domain.Counter, Delta: toInt64(int64(1))}
}

func (c *Collector) GetMetrics() map[string]domain.Metrics {
	return c.metrics
}

func getRandomValue() *float64 {
	const value = 100
	randValue := float64(rand.Intn(value))
	return &randValue
}

func toInt64(i int64) *int64 {
	return &i
}

func toFloat64(f float64) *float64 {
	return &f
}

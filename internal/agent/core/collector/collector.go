package collector

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
)

type Collector struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	cfg            *config.Config
	sync.Mutex
}

func NewCollector(cfg *config.Config) *Collector {
	return &Collector{
		cfg:            cfg,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
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

	c.Lock()
	defer c.Unlock()
	for metricName := range c.cfg.RuntimeMetrics {
		v := memStatsValue.FieldByName(metricName)

		var floatValue float64
		if v.CanUint() {
			floatValue = float64(v.Uint())
		}
		if v.CanFloat() {
			floatValue = v.Float()
		}

		c.gaugeMetrics[metricName] = floatValue
	}

	c.gaugeMetrics["RandomValue"] = getRandomValue()
	c.counterMetrics["PollCount"]++
}

func (c *Collector) GetGaugeMetrics() map[string]float64 {
	return c.gaugeMetrics
}

func (c *Collector) GetCounterMetrics() map[string]int64 {
	return c.counterMetrics
}

func (c *Collector) ResetCounter() {
	c.Lock()
	defer c.Unlock()
	for k := range c.counterMetrics {
		c.counterMetrics[k] = 0
	}
}

func getRandomValue() float64 {
	const value = 100
	randValue := float64(rand.Intn(value))
	return randValue
}

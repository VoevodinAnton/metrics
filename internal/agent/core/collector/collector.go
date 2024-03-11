package collector

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	reflect_copy "golang.design/x/reflect"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/shirou/gopsutil/v3/mem"
)

type Collector struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	results        chan domain.UploadResult
	cfg            *config.Config
	sync.Mutex
}

func NewCollector(cfg *config.Config) *Collector {
	return &Collector{
		cfg:            cfg,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
		results:        make(chan domain.UploadResult),
	}
}

func (c *Collector) Run() {
	ticker := time.NewTicker(c.cfg.PollInterval)
	for range ticker.C {
		go c.updateMetrics()
		go c.updateGopsutilMetrics()

		select {
		case r := <-c.results:
			if r.Err != nil {
				zap.L().Error("Error sending metrics", zap.Error(r.Err))
			}
		default:
		}
	}
}

func (c *Collector) updateGopsutilMetrics() {
	c.Lock()
	defer c.Unlock()
	v, err := mem.VirtualMemory()
	if err != nil {
		c.results <- domain.UploadResult{Err: err}
	}

	totalMemoryValue := float64(v.Total)
	freeMemoryValue := float64(v.Free)
	usePersentValue := v.UsedPercent

	c.gaugeMetrics["TotalMemory"] = totalMemoryValue
	c.gaugeMetrics["FreeMemory"] = freeMemoryValue
	c.gaugeMetrics["CPUutilization1"] = usePersentValue
}

func (c *Collector) updateMetrics() {
	c.Lock()
	defer c.Unlock()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memStatsValue := reflect.ValueOf(memStats)
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
	c.Lock()
	defer c.Unlock()
	return reflect_copy.DeepCopy[map[string]float64](c.gaugeMetrics)
}

func (c *Collector) GetCounterMetrics() map[string]int64 {
	c.Lock()
	defer c.Unlock()
	return reflect_copy.DeepCopy[map[string]int64](c.counterMetrics)
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

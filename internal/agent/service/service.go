package service

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/models"
)

type Store interface {
	UpdateGauge(metric *models.GaugeMetric) error
	GetGauge(name string) (*models.GaugeMetric, error)
	UpdateCounter(metric *models.CounterMetric) error
	GetCounter(name string) (*models.CounterMetric, error)
	ResetCounter(name string) error
}

type service struct {
	store Store
	cfg   *config.Config
}

func New(cfg *config.Config, store Store) *service {
	return &service{
		store: store,
		cfg:   cfg,
	}
}

func (s *service) Start() {
	s.runAgent()
}

func (s *service) runAgent() {
	go func() {
		for {
			s.updateMetrics()
			time.Sleep(s.cfg.PollInterval)
		}
	}()
	for {
		time.Sleep(s.cfg.ReportInterval)
		err := s.sendMetrics()
		if err != nil {
			fmt.Println("Error sending metrics:", err)
		}
	}
}

func (s *service) updateMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memStatsValue := reflect.ValueOf(memStats)

	for metricName := range s.cfg.RuntimeMetrics {
		v := memStatsValue.FieldByName(metricName)
		if v.CanUint() {
			err := s.store.UpdateGauge(&models.GaugeMetric{Name: metricName, Type: models.Counter, Value: float64(v.Uint())})
			if err != nil {
				log.Println(err)
				continue
			}
		}
		if v.CanFloat() {
			err := s.store.UpdateGauge(&models.GaugeMetric{Name: metricName, Type: models.Gauge, Value: v.Float()})
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	_ = s.store.UpdateGauge(&models.GaugeMetric{Name: "RandomValue", Type: models.Gauge, Value: getRandomValue()})
	_ = s.store.UpdateCounter(&models.CounterMetric{Name: "PollCount", Type: models.Counter, Value: 1})
}

func (s *service) sendMetrics() error {
	for metricName, metricType := range s.cfg.RuntimeMetrics {
		metric, err := s.store.GetGauge(metricName)
		if err != nil {
			log.Println(err)
			continue
		}

		url := fmt.Sprintf("http://%s/update/%s/%s/%f", s.cfg.ServerAddress, metricType, metricName, metric.Value)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}
		_ = resp.Body.Close()
	}

	for metricName, metricType := range s.cfg.CustomMetrics {
		var url string
		switch metricType {
		case "gauge":
			metric, err := s.store.GetGauge(metricName)
			if err != nil {
				log.Println(err)
				continue
			}
			url = fmt.Sprintf("http://%s/update/%s/%s/%f", s.cfg.ServerAddress, metricType, metricName, metric.Value)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				continue
			}
			_ = resp.Body.Close()
		case "counter":
			metric, err := s.store.GetCounter(metricName)
			if err != nil {
				log.Println(err)
				continue
			}
			url = fmt.Sprintf("http://%s/update/%s/%s/%d", s.cfg.ServerAddress, metricType, metricName, metric.Value)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				continue
			}
			_ = resp.Body.Close()
			_ = s.store.ResetCounter(metricName)
		}
	}

	return nil
}

func getRandomValue() float64 {
	const value = 100
	return float64(rand.Intn(value))
}

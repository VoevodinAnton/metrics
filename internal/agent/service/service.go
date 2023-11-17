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
	"github.com/pkg/errors"
)

const (
	ContentTypeText = "text/plain"
)

type Store interface {
	UpdateGauge(metric models.Metric) error
	UpdateCounter(metric models.Metric) error
	GetCounterMetrics() (map[string]*models.Metric, error)
	GetGaugeMetrics() (map[string]*models.Metric, error)

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
	s.Run()
}

func (s *service) Run() {
	go func() {
		for {
			s.updateMetrics()
			time.Sleep(s.cfg.PollInterval)
		}
	}()
	ticker := time.NewTicker(s.cfg.ReportInterval)
	for range ticker.C {
		if err := s.SendCounterMetrics(); err != nil {
			log.Println("[ERROR] sendCounterMetrics error", err)
		}
		if err := s.SendGaugeMetrics(); err != nil {
			log.Println("[ERROR] sendGaugeMetrics error", err)
		}
	}
}

func (s *service) updateMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memStatsValue := reflect.ValueOf(memStats)

	for metricName := range s.cfg.RuntimeMetrics {
		v := memStatsValue.FieldByName(metricName)
		var gaugeMetric models.Metric
		if v.CanUint() {
			gaugeMetric = models.Metric{Name: metricName, Type: models.Counter, Value: float64(v.Uint())}
		}
		if v.CanFloat() {
			gaugeMetric = models.Metric{Name: metricName, Type: models.Gauge, Value: v.Float()}
		}

		err := s.store.UpdateGauge(gaugeMetric)
		if err != nil {
			log.Println("[ERROR] updateGauge error", err)
			continue
		}
	}

	_ = s.store.UpdateGauge(models.Metric{Name: "RandomValue", Type: models.Gauge, Value: getRandomValue()})
	_ = s.store.UpdateCounter(models.Metric{Name: "PollCount", Type: models.Counter, Value: int64(1)})
}

func (s *service) SendCounterMetrics() error {
	counterMetrics, err := s.store.GetCounterMetrics()
	if err != nil {
		return errors.Wrap(err, "GetCounterMetrics")
	}
	for metricName, metric := range counterMetrics {
		url := fmt.Sprintf("http://%s/update/counter/%s/%d", s.cfg.ServerAddress, metricName, metric.Value)
		if err := s.sendMetrics(url); err != nil {
			log.Println("[ERROR] send counter metrics error", err)
			continue
		}
		_ = s.store.ResetCounter(metricName)
	}

	return nil
}

func (s *service) SendGaugeMetrics() error {
	gaugeMetrics, err := s.store.GetGaugeMetrics()
	if err != nil {
		return errors.Wrap(err, "GetGaugeMetrics")
	}
	for metricName, metric := range gaugeMetrics {
		url := fmt.Sprintf("http://%s/update/gauge/%s/%f", s.cfg.ServerAddress, metricName, metric.Value)
		if err := s.sendMetrics(url); err != nil {
			log.Println("[ERROR] send gauge metrics error", err)
			continue
		}
	}

	return nil
}

func (s *service) sendMetrics(url string) error {
	resp, err := http.Post(url, ContentTypeText, nil)
	if err != nil {
		return errors.Wrap(err, "http.Get")
	}
	if resp.StatusCode > http.StatusBadRequest {
		return errors.Wrap(errors.New("status code > 400"), resp.Status)
	}
	if err = resp.Body.Close(); err != nil {
		return errors.Wrap(err, "body.Close")
	}

	return nil
}

func getRandomValue() float64 {
	const value = 100
	return float64(rand.Intn(value))
}

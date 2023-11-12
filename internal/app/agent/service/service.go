package service

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/VoevodinAnton/metrics/internal/app/agent/config"
	"github.com/VoevodinAnton/metrics/internal/app/agent/repository"
)

type Service interface {
	Start()
}

type service struct {
	metricRepo repository.MetricRepository
	cfg        *config.Config
}

func New(cfg *config.Config, metricRepo repository.MetricRepository) *service {
	return &service{
		metricRepo: metricRepo,
		cfg:        cfg,
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

	for metricName, _ := range s.cfg.RuntimeMetrics {
		v := memStatsValue.FieldByName(metricName)
		if v.CanFloat() {
			s.metricRepo.UpdateGaugeMetric(metricName, v.Float())
		}
		if v.CanUint() {
			s.metricRepo.UpdateGaugeMetric(metricName, float64(v.Uint()))
		}
	}

	s.metricRepo.UpdateGaugeMetric("RandomValue", getRandomValue())
	s.metricRepo.UpdateCounterMetric("PollCount", 1)
}

func (s *service) sendMetrics() error {
	for metricName, metricType := range s.cfg.RuntimeMetrics {
		metricValue, err := s.metricRepo.GetGaugeMetric(metricName)
		if err != nil {
			log.Println(err)
			continue
		}

		url := fmt.Sprintf("http://%s/update/%s/%s/%f", s.cfg.Server, metricType, metricName, metricValue)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}
		resp.Body.Close()
	}

	for metricName, metricType := range s.cfg.CustomMetrics {
		var url string
		switch metricType {
		case "gauge":
			metricValue, err := s.metricRepo.GetGaugeMetric(metricName)
			if err != nil {
				log.Println(err)
				continue
			}
			url = fmt.Sprintf("http://%s/update/%s/%s/%f", s.cfg.Server, metricType, metricName, metricValue)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				continue
			}
			resp.Body.Close()
		case "counter":
			metricValue, err := s.metricRepo.GetCounterMetric(metricName)
			if err != nil {
				log.Println(err)
				continue
			}
			url = fmt.Sprintf("http://%s/update/%s/%s/%d", s.cfg.Server, metricType, metricName, metricValue)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				continue
			}
			resp.Body.Close()
			s.metricRepo.ResetCounterMetric(metricName)
		}
	}

	return nil
}

func getRandomValue() float64 {
	return float64(rand.Intn(100))
}

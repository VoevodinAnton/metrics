package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

type Metric struct {
	Type  MetricType
	Value float64
}

type MemStorage struct {
	metrics map[string]Metric
	mu      sync.Mutex
}

type MetricStorage interface {
	GetMetric(name string) (Metric, bool)
	UpdateMetric(name string, metric Metric)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]Metric),
	}
}

func (m *MemStorage) GetMetric(name string) (Metric, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.metrics[name]
	return val, ok
}

func (m *MemStorage) UpdateMetric(name string, metric Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch metric.Type {
	case Gauge:
		m.metrics[name] = metric
	case Counter:
		prevMetric, ok := m.metrics[name]
		if ok {
			metric.Value += prevMetric.Value
		}
		m.metrics[name] = metric
	}
}

func handleUpdate(storage MetricStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Invalid URL", http.StatusNotFound)
			return
		}

		metricType := parts[2]
		metricName := parts[3]
		metricValue := parts[4]

		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		var metricTypeVal MetricType
		switch metricType {
		case "gauge":
			metricTypeVal = Gauge
		case "counter":
			metricTypeVal = Counter
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		storage.UpdateMetric(metricName, Metric{Type: metricTypeVal, Value: value})

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	storage := NewMemStorage()

	http.HandleFunc("/update/", handleUpdate(storage))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

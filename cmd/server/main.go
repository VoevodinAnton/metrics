package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi"
)

var pageTemplate = template.Must(template.New("metrics").Parse(`
<html>
<head>
  <title>Metric List</title>
</head>
<body>
  <h1>Metric List</h1>
  <ul>
    {{range $name, $metric := .}}
      <li><strong>{{$name}}:</strong> {{$metric.Value}}</li>
    {{end}}
  </ul>
</body>
</html>
`))

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
	GetMetrics() map[string]Metric
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

func (m *MemStorage) GetMetrics() map[string]Metric {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.metrics
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

func handleUpdateMetric(storage MetricStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := strings.ToLower(chi.URLParam(r, "metricType"))
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

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

func handleGetMetric(storage MetricStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := strings.ToLower(chi.URLParam(r, "metricType"))
		metricName := strings.ToLower(chi.URLParam(r, "metricName"))

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

		_ = metricTypeVal

		metric, ok := storage.GetMetric(metricName)
		if !ok {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("%f", metric.Value)))
		w.WriteHeader(http.StatusOK)
	}
}

func handleGetMetrics(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pageBuffer bytes.Buffer
		pageTemplate.Execute(&pageBuffer, storage.GetMetrics())

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(pageBuffer.Bytes())
	}
}

func main() {
	var serverAddress string
	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.Parse()

	envServerAddress := os.Getenv("ADDRESS")
	if envServerAddress != "" {
		serverAddress = envServerAddress
	}

	storage := NewMemStorage()

	r := chi.NewRouter()
	r.Get("/update/{metricType}/{metricName}/{metricValue}", handleUpdateMetric(storage))
	r.Get("/value/{metricType}/{metricName}", handleGetMetric(storage))
	r.Get("/", handleGetMetrics(storage))

	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

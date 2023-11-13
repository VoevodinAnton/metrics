package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/VoevodinAnton/metrics/internal/models"
	"github.com/go-chi/chi/v5"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type Handler struct {
	Service Service
}

func (h *Handler) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "Invalid metric value", http.StatusBadRequest)
		return
	}

	var metricTypeVal models.MetricType
	switch metricType {
	case gauge:
		metricTypeVal = models.Gauge
	case counter:
		metricTypeVal = models.Counter
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateMetric(&models.MetricReq{Name: metricName, Type: metricTypeVal, Value: value})
	if err != nil {
		http.Error(w, "Failed update metric", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	var metricTypeVal models.MetricType
	switch metricType {
	case gauge:
		metricTypeVal = models.Gauge
	case counter:
		metricTypeVal = models.Counter
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	metric, err := h.Service.GetMetric(&models.MetricReq{Name: metricName, Type: metricTypeVal})
	if err != nil {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	if metricTypeVal == models.Gauge {
		fmt.Fprint(w, metric.Value)
	} else if metricTypeVal == models.Counter {
		fmt.Fprint(w, int64(metric.Value))
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var pageTemplate = template.Must(template.New("metrics").Parse(`
	<html>
	<head>
	  <title>Metric List</title>
	</head>
	<body>
	  <h1>Metric List</h1>
	  <ul>
		{{range $metric := .}}
		  <li><strong>{{$metric.Name}}:</strong> {{$metric.Value}}</li>
		{{end}}
	  </ul>
	</body>
	</html>
	`))

	metrics, err := h.Service.GetMetrics()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var pageBuffer bytes.Buffer
	err = pageTemplate.Execute(&pageBuffer, metrics)
	if err != nil {
		http.Error(w, "Failed execute page template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pageBuffer.Bytes())
}

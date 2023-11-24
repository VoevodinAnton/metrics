package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/go-chi/chi/v5"
)

const (
	metricTypeURLParam  = "metricType"
	metricNameURLParam  = "metricName"
	metricValueURLParam = "metricValue"

	ContentTypeHeader = "Content-Type"
	ContentTypeText   = "text/plain; charset=utf-8"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeJSON   = "application/json"
)

type Handler struct {
	Service Service
}

func (h *Handler) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, metricTypeURLParam)
	metricName := chi.URLParam(r, metricNameURLParam)
	metricValue := chi.URLParam(r, metricValueURLParam)

	req := domain.Metrics{ID: metricName}

	var err error
	switch metricType {
	case domain.Gauge:
		var metricVal float64
		metricVal, err = strconv.ParseFloat(metricValue, 64)
		req.Value = &metricVal
	case domain.Counter:
		var metricVal int64
		metricVal, err = strconv.ParseInt(metricValue, 10, 64)
		req.Delta = &metricVal
	default:
		http.Error(w, ErrInvalidMetricType.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, ErrInvalidMetricValue.Error(), http.StatusBadRequest)
		return
	}
	req.MType = metricType

	err = h.Service.UpdateMetric(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(ContentTypeHeader, ContentTypeText)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, metricTypeURLParam)
	metricName := chi.URLParam(r, metricNameURLParam)

	metricReq := &domain.Metrics{ID: metricName, MType: metricType}
	metric, err := h.Service.GetMetric(metricReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set(ContentTypeHeader, ContentTypeText)

	switch metricType {
	case domain.Counter:
		fmt.Fprint(w, *metric.Delta)
	case domain.Gauge:
		fmt.Fprint(w, *metric.Value)
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
		  <li><strong>{{$metric.ID}}:</strong> {{if $metric.Value}} {{$metric.Value}} {{else}} {{$metric.Delta}} {{end}}</li>
		{{end}}
	  </ul>
	</body>
	</html>
	`))

	metrics, err := h.Service.GetMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var pageBuffer bytes.Buffer
	err = pageTemplate.Execute(&pageBuffer, *metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeHTML)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pageBuffer.Bytes())
}

func (h *Handler) UpdateJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metricUpdate domain.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err := h.Service.UpdateMetric(&metricUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metricReq domain.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	metric, err := h.Service.GetMetric(&metricReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metricResp, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(metricResp)
}

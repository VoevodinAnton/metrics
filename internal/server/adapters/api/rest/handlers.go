package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/VoevodinAnton/metrics/internal/pkg/constants"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/middlewares"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	metricTypeURLParam  = "metricType"
	metricNameURLParam  = "metricName"
	metricValueURLParam = "metricValue"
)

type Handler struct {
	service Service
	mw      middlewares.MiddlewareManager
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

	err = h.service.UpdateMetric(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentTypeHeader, constants.ContentTypeText)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, metricTypeURLParam)
	metricName := chi.URLParam(r, metricNameURLParam)

	metricReq := &domain.Metrics{ID: metricName, MType: metricType}
	metric, err := h.service.GetMetric(metricReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set(constants.ContentTypeHeader, constants.ContentTypeText)

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

	metrics, err := h.service.GetMetrics()
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

	w.Header().Set(constants.ContentTypeHeader, constants.ContentTypeHTML)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pageBuffer.Bytes())
}

func (h *Handler) UpdateJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metricUpdate domain.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricUpdate); err != nil {
		zap.L().Error("UpdateJSONMetricHandler json.NewDecoder", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	zap.L().Debug("UpdateJSONMetricHandler", zap.Reflect("json.NewDecoder1", metricUpdate))
	err := h.service.UpdateMetric(&metricUpdate)
	if err != nil {
		zap.L().Error("UpdateJSONMetricHandler service.UpdateMetric", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	var metricReq domain.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricReq); err != nil {
		zap.L().Error("GetJSONMetricHandler json.NewDecoder", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	zap.L().Debug("GetJSONMetricHandler", zap.Reflect("json.NewDecoder2", metricReq))
	metric, err := h.service.GetMetric(&metricReq)
	if err != nil {
		zap.L().Error("GetJSONMetricHandler service.GetMetric", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metricResp, err := json.Marshal(metric)
	if err != nil {
		zap.L().Error("GetJSONMetricHandler json.Marshal", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zap.L().Debug("GetJSONMetricHandler1", zap.Reflect("json.Marshal", metricResp))
	w.Header().Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(metricResp)
}

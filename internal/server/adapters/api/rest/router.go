package api

import (
	"net/http"

	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/middlewares"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

type Service interface {
	GetMetric(metric *domain.Metrics) (*domain.Metrics, error)
	UpdateMetric(metric *domain.Metrics) error
	GetMetrics() (*[]domain.Metrics, error)
}

type Router struct {
	cfg *config.Config
	r   *chi.Mux
}

func NewRouter(cfg *config.Config, service Service, mw middlewares.MiddlewareManager) *Router {
	h := Handler{
		service: service,
		mw:      mw,
	}
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		logging.WithLogging,
	)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetMetricHandler)

	gzipGroup := r.Group(nil)
	gzipGroup.Use(mw.GzipCompressHandle, mw.GzipDecompressHandle)
	gzipGroup.Post("/update", h.UpdateJSONMetricHandler)
	gzipGroup.Get("/", h.GetMetricsHandler)
	gzipGroup.Post("/value/", h.GetJSONMetricHandler)

	return &Router{
		r:   r,
		cfg: cfg,
	}
}

func (r *Router) ServeRouter() error {
	err := http.ListenAndServe(r.cfg.Server.Address, r.r)
	return errors.Wrap(err, "http.ListenAndServe")
}

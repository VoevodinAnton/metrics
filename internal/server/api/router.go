package api

import (
	"net/http"

	"github.com/VoevodinAnton/metrics/internal/models"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

type Service interface {
	GetMetric(req *models.Metric) (*models.Metric, error)
	UpdateMetric(req models.Metric) error
	GetMetrics() ([]*models.Metric, error)
}

type Router struct {
	cfg *config.Config
	r   *chi.Mux
}

func NewRouter(cfg *config.Config, service Service) *Router {
	h := Handler{
		Service: service,
	}
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
	)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetricHandler)
	r.Get("/value/{metricType}/{metricName}", h.GetMetricHandler)
	r.Get("/", h.GetMetricsHandler)

	return &Router{
		r:   r,
		cfg: cfg,
	}
}

func (r *Router) ServeRouter() error {
	err := http.ListenAndServe(r.cfg.Server.Address, r.r)
	return errors.Wrap(err, "http.ListenAndServe")
}

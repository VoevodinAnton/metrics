package api

import (
	"net/http"

	"github.com/VoevodinAnton/metrics/internal/models"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Service interface {
	GetMetric(req *models.MetricReq) (*models.MetricResp, error)
	UpdateMetric(req *models.MetricReq) error
	GetMetrics() ([]*models.MetricResp, error)
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
		render.SetContentType(render.ContentTypeJSON), // forces Content-type
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
	return http.ListenAndServe(r.cfg.Server.Address, r.r) // nolint: wrapcheck
}

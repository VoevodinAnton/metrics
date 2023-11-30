package uploader

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	ContentTypeText = "text/plain"
)

type Store interface {
	GetMetrics() map[string]domain.Metrics
}

type Uploader struct {
	cfg   *config.Config
	store Store
	sync.Mutex
}

func NewUploader(cfg *config.Config, store Store) *Uploader {
	return &Uploader{
		cfg:   cfg,
		store: store,
	}
}

func (u *Uploader) Run() {
	ticker := time.NewTicker(u.cfg.ReportInterval)
	for range ticker.C {
		u.sendMetrics()
	}
}

func (u *Uploader) sendMetrics() {
	metrics := u.store.GetMetrics()
	for _, m := range metrics {
		url := fmt.Sprintf("http://%s/update", u.cfg.ServerAddress)
		err := u.Upload(url, m)
		if err != nil {
			zap.L().Error("svc.sendMetrics gauge", zap.Error(err))
			continue
		}
	}
}

func (u *Uploader) Upload(url string, m domain.Metrics) error {
	metricReq, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err = w.Write(metricReq)
	if err != nil {
		return errors.Wrap(err, "writer.Write")
	}
	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "writer.Close")
	}

	resp, err := http.Post(url, ContentTypeText, &b)
	if err != nil {
		return errors.Wrap(err, "http.Get")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("status code != 200"), resp.Status)
	}
	if err = resp.Body.Close(); err != nil {
		return errors.Wrap(err, "body.Close")
	}

	return nil
}

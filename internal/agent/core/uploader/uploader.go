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
	"github.com/VoevodinAnton/metrics/internal/pkg/constants"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	updateURLTemplate = "http://%s/update"
	clientTimeout     = 60 * time.Second
)

type Store interface {
	GetGaugeMetrics() map[string]float64
	GetCounterMetrics() map[string]int64
	ResetCounter()
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
		u.sendGaugeMetrics()
		u.sendCounterMetrics()
	}
}

func (u *Uploader) sendGaugeMetrics() {
	metrics := u.store.GetGaugeMetrics()
	for name, value := range metrics {
		value := value
		url := fmt.Sprintf(updateURLTemplate, u.cfg.ServerAddress)

		m := domain.Metrics{
			ID:    name,
			MType: domain.Gauge,
			Value: &value,
		}

		if m.ID == "StackSys" {
			fmt.Println(*m.Value)
		}

		err := u.Upload(url, m)
		if err != nil {
			zap.L().Error("upload gauge", zap.Error(err))
			continue
		}
	}
}

func (u *Uploader) sendCounterMetrics() {
	metrics := u.store.GetCounterMetrics()
	for name, value := range metrics {
		value := value
		url := fmt.Sprintf(updateURLTemplate, u.cfg.ServerAddress)

		m := domain.Metrics{
			ID:    name,
			MType: domain.Counter,
			Delta: &value,
		}

		err := u.Upload(url, m)
		if err != nil {
			zap.L().Error("upload counter", zap.Error(err))
			continue
		}

		u.store.ResetCounter()
	}
}

func (u *Uploader) Upload(url string, m domain.Metrics) error {
	client := http.Client{
		Timeout: clientTimeout,
	}
	metricsJSON, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err = w.Write(metricsJSON)
	if err != nil {
		return errors.Wrap(err, "writer.Write")
	}
	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "writer.Close")
	}
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.ContentEncodingHeader, constants.GzipEncoding)
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("status code != 200"), resp.Status)
	}
	if err = resp.Body.Close(); err != nil {
		return errors.Wrap(err, "body.Close")
	}

	return nil
}

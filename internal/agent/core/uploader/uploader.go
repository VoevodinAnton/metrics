package uploader

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/pkg/constants"
	"github.com/VoevodinAnton/metrics/internal/pkg/domain"
	"github.com/VoevodinAnton/metrics/internal/pkg/semaphore"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

const (
	updateURLTemplate = "http://%s/updates"
	clientTimeout     = 10 * time.Second
	batchSize         = 5
	metricsToSendSize = 1000
)

type Store interface {
	GetGaugeMetrics() map[string]float64
	GetCounterMetrics() map[string]int64
	ResetCounter()
}

type Uploader struct {
	cfg           *config.Config
	cb            *gobreaker.CircuitBreaker
	semaphore     *semaphore.Semaphore
	wg            *sync.WaitGroup
	results       chan domain.UploadResult
	metricsToSend chan domain.Metrics
	store         Store
	sync.Mutex
}

func NewUploader(cfg *config.Config, store Store) *Uploader {
	var st gobreaker.Settings
	st.Name = "HTTP REQUEST"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests > 20 && failureRatio >= 0.7
	}
	return &Uploader{
		cfg:           cfg,
		store:         store,
		semaphore:     semaphore.NewSemaphore(cfg.RateLimit),
		wg:            &sync.WaitGroup{},
		results:       make(chan domain.UploadResult),
		metricsToSend: make(chan domain.Metrics, metricsToSendSize),
		cb:            gobreaker.NewCircuitBreaker(st),
	}
}

func (u *Uploader) Run() {
	go u.sendMetrics()

	ticker := time.NewTicker(u.cfg.ReportInterval)

	for {
		select {
		case <-ticker.C:
			go u.report()
		case r := <-u.results:
			if r.Err != nil {
				zap.L().Error("Error sending metrics", zap.Error(r.Err))
			}
		}
	}
}

func (u *Uploader) report() {
	u.Lock()
	defer u.Unlock()
	gaugeMetrics := u.store.GetGaugeMetrics()
	for name, value := range gaugeMetrics {
		value := value
		m := domain.Metrics{
			ID:    name,
			MType: domain.Gauge,
			Value: &value,
		}
		u.metricsToSend <- m
	}
	counterMetrics := u.store.GetCounterMetrics()

	for name, value := range counterMetrics {
		value := value
		m := domain.Metrics{
			ID:    name,
			MType: domain.Counter,
			Delta: &value,
		}
		u.metricsToSend <- m
	}
}

func (u *Uploader) sendMetrics() {
	metricsBatch := make([]domain.Metrics, 0)
	url := fmt.Sprintf(updateURLTemplate, u.cfg.ServerAddress)
	for idx := 0; idx < 2; idx++ {
		go func() {
			for metric := range u.metricsToSend {
				metricsBatch = append(metricsBatch, metric)
				if len(metricsBatch) == batchSize {
					u.semaphore.Acquire()
					err := u.Upload(url, metricsBatch)
					if err != nil {
						u.results <- domain.UploadResult{Err: err}
					}
					metricsBatch = make([]domain.Metrics, 0)
					u.semaphore.Release()
				}
			}
		}()
	}

	u.store.ResetCounter()
}

func (u *Uploader) Upload(url string, m []domain.Metrics) error {
	_, err := u.cb.Execute(func() (interface{}, error) {
		client := http.Client{
			Timeout: clientTimeout,
		}
		metricsJSON, err := json.Marshal(m)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal")
		}
		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		_, err = w.Write(metricsJSON)
		if err != nil {
			return nil, errors.Wrap(err, "writer.Write")
		}
		err = w.Close()
		if err != nil {
			return nil, errors.Wrap(err, "writer.Close")
		}
		req, err := http.NewRequest(http.MethodPost, url, &b)
		if err != nil {
			return nil, errors.Wrap(err, "http.NewRequest")
		}

		if u.cfg.Key != "" {
			h := hmac.New(sha256.New, []byte(u.cfg.Key))

			h.Write(metricsJSON)
			metricsHash := h.Sum(nil)

			hashString := hex.EncodeToString(metricsHash)

			req.Header.Add(constants.HashSHA256, hashString)
		}

		req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
		req.Header.Set(constants.ContentEncodingHeader, constants.GzipEncoding)
		resp, err := client.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "client.Do")
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.Wrap(errors.New("status code != 200"), resp.Status)
		}
		if err = resp.Body.Close(); err != nil {
			return nil, errors.Wrap(err, "body.Close")
		}

		return nil, nil //nolint: nilnil // currect return
	})
	if err != nil {
		return errors.Wrap(err, "cb.Execute")
	}

	return nil
}

package backup

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/internal/server/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	writeFilePerm = 0600
)

type Store interface {
	PutCounter(ctx context.Context, update models.Metric) error
	PutGauge(ctx context.Context, update models.Metric) error
	GetCounterMetrics(ctx context.Context) (map[string]models.Metric, error)
	GetGaugeMetrics(ctx context.Context) (map[string]models.Metric, error)
}

type Backuper struct {
	store Store
	cfg   *config.Config
}

func New(cfg *config.Config, store Store) *Backuper {
	return &Backuper{
		store: store,
		cfg:   cfg,
	}
}

func (b *Backuper) Run(ctx context.Context) {
	ticker := time.NewTicker(b.cfg.StoreInterval)
	for range ticker.C {
		err := b.SaveMetricsToFile(ctx)
		if err != nil {
			zap.L().Error("saveMetricsToFile", zap.Error(err))
			continue
		}
		zap.L().Sugar().Infof("metrics saved to file %s", b.cfg.FilePath)
	}
}

func (b *Backuper) SaveMetricsToFile(ctx context.Context) error {
	metrics := make(map[string]models.Metric)
	gaugeMetrics, err := b.store.GetGaugeMetrics(ctx)
	if err != nil {
		return errors.Wrap(err, "store.GetGaugeMetrics")
	}
	counterMetrics, err := b.store.GetCounterMetrics(ctx)
	if err != nil {
		return errors.Wrap(err, "store.GetCounterMetrics")
	}
	for k, v := range gaugeMetrics {
		metrics[k] = v
	}
	for k, v := range counterMetrics {
		metrics[k] = v
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	err = os.WriteFile(b.cfg.FilePath, data, writeFilePerm)
	if err != nil {
		return errors.Wrap(err, "os.WriteFile")
	}

	return nil
}

func (b *Backuper) RestoreMetricsFromFile(ctx context.Context) error {
	file, err := os.Open(b.cfg.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return errors.Wrap(err, "os.Open")
		}
	}
	defer func() {
		err := file.Close()
		if err != nil {
			zap.L().Error("file.Close", zap.Error(err))
		}
	}()

	metrics := make(map[string]models.Metric)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metrics); err != nil {
		return errors.Wrap(err, "decoder.Decode")
	}

	for _, metric := range metrics {
		if metric.Type == models.Counter {
			v, _ := metric.Value.(float64)
			metric.Value = int64(v)
			_ = b.store.PutCounter(ctx, metric)
		}
		if metric.Type == models.Gauge {
			_ = b.store.PutGauge(ctx, metric)
		}
	}

	return nil
}

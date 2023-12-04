package main

import (
	"os"
	"os/signal"
	"syscall"

	api "github.com/VoevodinAnton/metrics/internal/server/adapters/api/rest"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/backup"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/middlewares"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/store/memory"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/internal/server/core/service"
	logger "github.com/VoevodinAnton/metrics/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	cfg := config.InitConfig()
	logger.NewLogger(cfg.Logger)
	defer logger.Close()
	mw := middlewares.NewMiddlewareManager()
	storage := memory.NewStorage()
	backup := backup.New(cfg, storage)

	if cfg.Restore {
		err := backup.RestoreMetricsFromFile()
		if err != nil {
			zap.L().Error("empty start", zap.Error(err))
		}
	}
	if cfg.FilePath != "" {
		go func() {
			backup.Run()
		}()
		defer func() {
			err := backup.SaveMetricsToFile()
			if err != nil {
				zap.L().Error("backup.SaveMetricsToFile", zap.Error(err))
			}
		}()
	}

	service := service.New(storage)
	r := api.NewRouter(cfg, service, mw)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	listenErr := make(chan error, 1)
	listenSignals := make(chan os.Signal, 1)
	signal.Notify(listenSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		listenErr <- r.ServeRouter()
	}()

	select {
	case sig := <-listenSignals:
		zap.L().Warn("received signal", zap.String("signal", sig.String()))
	case err := <-listenErr:
		zap.L().Error("", zap.Error(err))
	}
}

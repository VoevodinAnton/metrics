package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	api "github.com/VoevodinAnton/metrics/internal/server/adapters/api/rest"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/backup"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/middlewares"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/store"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/internal/server/core/service"
	logger "github.com/VoevodinAnton/metrics/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	logger.NewLogger(cfg.Logger)
	defer logger.Close()
	mw := middlewares.NewMiddlewareManager()
	ctx := context.Background()
	storage, err := store.NewStore(cfg)
	if err != nil {
		zap.L().Fatal("store.NewStore", zap.Error(err))
	}
	defer storage.Close()

	backup := backup.New(cfg, storage)
	if cfg.Restore {
		err := backup.RestoreMetricsFromFile(ctx)
		if err != nil {
			zap.L().Error("empty start", zap.Error(err))
		}
	}
	if cfg.FilePath != "" {
		go func() {
			backup.Run(ctx)
		}()
		defer func() {
			err := backup.SaveMetricsToFile(ctx)
			if err != nil {
				zap.L().Error("backup.SaveMetricsToFile", zap.Error(err))
			}
		}()
	}

	service := service.New(storage)
	r := api.NewRouter(cfg, service, mw)

	listenErr := make(chan error, 1)
	listenSignals := make(chan os.Signal, 1)
	signal.Notify(listenSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		listenErr <- r.ServeRouter()
	}()

	zap.L().Sugar().Infof("The server is listening and serving the address %s", cfg.Server.Address)
	select {
	case sig := <-listenSignals:
		zap.L().Warn("received signal", zap.String("signal", sig.String()))
		storage.Close()
	case err := <-listenErr:
		zap.L().Error("", zap.Error(err))
		storage.Close()
	}
}

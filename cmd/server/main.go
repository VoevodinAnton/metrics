package main

import (
	"context"
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
	"github.com/VoevodinAnton/metrics/pkg/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
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
	storage := memory.NewStorage()
	backup := backup.New(cfg, storage)

	var db *pgxpool.Pool
	if cfg.Postgres.DatabaseDSN != "" {
		db, err = postgres.NewPgxConn(context.Background(), cfg.Postgres)
		if err != nil {
			zap.L().Fatal("postgres.NewPgxConn", zap.Error(err))
		}
	}

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
	r := api.NewRouter(cfg, service, mw, db)

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
		db.Close()
	case err := <-listenErr:
		zap.L().Error("", zap.Error(err))
		db.Close()
	}
}

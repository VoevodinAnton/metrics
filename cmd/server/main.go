package main

import (
	"github.com/VoevodinAnton/metrics/internal/server/adapters/api"
	"github.com/VoevodinAnton/metrics/internal/server/adapters/memory"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/internal/server/core/service"
	logger "github.com/VoevodinAnton/metrics/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	cfg := config.InitConfig()
	logger.NewLogger(cfg.Logger)
	defer logger.Close()
	storage := memory.NewStorage()
	service := service.New(storage)
	r := api.NewRouter(cfg, service)
	err := r.ServeRouter()
	if err != nil {
		zap.L().Fatal("Error starting server", zap.Error(err))
	}

	// zap.L().Info("Server started", zap.String("addr", cfg.Server.Address))
}

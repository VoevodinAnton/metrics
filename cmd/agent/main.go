package main

import (
	"github.com/VoevodinAnton/metrics/internal/adapters/memory"
	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/agent/service"
)

func main() {
	cfg := config.InitConfig()

	storage := memory.NewStorage()

	svc := service.New(cfg, storage)
	svc.Start()
}

package main

import (
	"github.com/VoevodinAnton/metrics/internal/app/agent/config"
	"github.com/VoevodinAnton/metrics/internal/app/agent/service"
	"github.com/VoevodinAnton/metrics/internal/app/agent/storage"
)

func main() {
	cfg := config.New()

	storage := storage.New()

	svc := service.New(cfg, storage)
	svc.Start()
}

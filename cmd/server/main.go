package main

import (
	"log"

	"github.com/VoevodinAnton/metrics/internal/adapters/memory"
	"github.com/VoevodinAnton/metrics/internal/server/api"
	"github.com/VoevodinAnton/metrics/internal/server/config"
	"github.com/VoevodinAnton/metrics/internal/server/core/service"
)

func main() {
	cfg := config.InitConfig()
	storage := memory.NewStorage()
	service := service.New(storage)
	r := api.NewRouter(cfg, service)
	err := r.ServeRouter()
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}

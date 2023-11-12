package main

import (
	"flag"
	"time"

	"github.com/VoevodinAnton/metrics/internal/app/agent/config"
	"github.com/VoevodinAnton/metrics/internal/app/agent/service"
	"github.com/VoevodinAnton/metrics/internal/app/agent/storage"
)

func main() {
	var serverAddress string
	var reportInterval, pollInterval int64

	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.Int64Var(&reportInterval, "r", 10, "Report interval in seconds")
	flag.Int64Var(&pollInterval, "p", 2, "Poll interval in seconds")
	flag.Parse()

	cfg := config.New()
	cfg.Server = serverAddress
	cfg.ReportInterval = time.Duration(reportInterval) * time.Second
	cfg.PollInterval = time.Duration(pollInterval) * time.Second

	storage := storage.New()

	svc := service.New(cfg, storage)
	svc.Start()
}

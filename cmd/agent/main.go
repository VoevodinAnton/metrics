package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/VoevodinAnton/metrics/internal/app/agent/config"
	"github.com/VoevodinAnton/metrics/internal/app/agent/service"
	"github.com/VoevodinAnton/metrics/internal/app/agent/storage"
)

func main() {
	var serverAddress string
	var reportInterval, pollInterval int

	envServerAddress := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")

	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&reportInterval, "r", 10, "Report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "Poll interval in seconds")
	flag.Parse()

	if envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envReportInterval != "" {
		reportInterval, _ = strconv.Atoi(envReportInterval)
	}
	if envPollInterval != "" {
		pollInterval, _ = strconv.Atoi(envPollInterval)
	}

	cfg := config.New()
	cfg.Server = serverAddress
	cfg.ReportInterval = time.Duration(reportInterval) * time.Second
	cfg.PollInterval = time.Duration(pollInterval) * time.Second

	storage := storage.New()

	svc := service.New(cfg, storage)
	svc.Start()
}

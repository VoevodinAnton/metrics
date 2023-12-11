package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/VoevodinAnton/metrics/pkg/config"
)

const (
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	Logger         *config.Logger
	CustomMetrics  map[string]string
	RuntimeMetrics map[string]string
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func InitConfig() *Config {
	var serverAddress string
	var reportInterval, pollInterval int

	envServerAddress := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")

	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&reportInterval, "r", defaultReportInterval, "Report interval in seconds")
	flag.IntVar(&pollInterval, "p", defaultPollInterval, "Poll interval in seconds")
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

	return &Config{
		ServerAddress:  serverAddress,
		PollInterval:   time.Duration(pollInterval) * time.Second,
		ReportInterval: time.Duration(reportInterval) * time.Second,
		RuntimeMetrics: map[string]string{
			"Alloc":         "gauge",
			"BuckHashSys":   "gauge",
			"Frees":         "gauge",
			"GCCPUFraction": "gauge",
			"GCSys":         "gauge",
			"HeapAlloc":     "gauge",
			"HeapIdle":      "gauge",
			"HeapInuse":     "gauge",
			"HeapObjects":   "gauge",
			"HeapReleased":  "gauge",
			"HeapSys":       "gauge",
			"LastGC":        "gauge",
			"Lookups":       "gauge",
			"NumGC":         "gauge",
			"MCacheInuse":   "gauge",
			"MCacheSys":     "gauge",
			"MSpanInuse":    "gauge",
			"MSpanSys":      "gauge",
			"Mallocs":       "gauge",
			"NextGC":        "gauge",
			"NumForcedGC":   "gauge",
			"OtherSys":      "gauge",
			"PauseTotalNs":  "gauge",
			"StackInuse":    "gauge",
			"StackSys":      "gauge",
			"Sys":           "gauge",
			"TotalAlloc":    "gauge",
		},
		CustomMetrics: map[string]string{
			"PollCount":   "counter",
			"RandomValue": "gauge",
		},
		Logger: &config.Logger{
			Development: true,
			Level:       "debug",
		},
	}
}

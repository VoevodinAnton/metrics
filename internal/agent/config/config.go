package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server         Server
	PollInterval   time.Duration
	ReportInterval time.Duration
	RuntimeMetrics map[string]string
	CustomMetrics  map[string]string
}

type Server struct {
	Address string
}

func InitConfig() *Config {
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

	return &Config{
		Server: Server{
			Address: serverAddress,
		},
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
	}
}

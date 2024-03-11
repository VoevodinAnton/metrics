package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/VoevodinAnton/metrics/pkg/config"
	"github.com/pkg/errors"
)

const (
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultRateLimit      = 10
)

type Config struct {
	Logger         *config.Logger
	CustomMetrics  map[string]string
	RuntimeMetrics map[string]string
	ServerAddress  string
	Key            string
	PollInterval   time.Duration
	ReportInterval time.Duration
	RateLimit      int
}

func InitConfig() (*Config, error) {
	var err error
	var serverAddress string
	var reportInterval, pollInterval, rateLimit int
	var key string

	envServerAddress := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")
	envKey := os.Getenv("KEY")
	envRateLimit := os.Getenv("RATE_LIMIT")

	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&reportInterval, "r", defaultReportInterval, "Report interval in seconds")
	flag.IntVar(&pollInterval, "p", defaultPollInterval, "Poll interval in seconds")
	flag.IntVar(&rateLimit, "l", defaultRateLimit, "Rate limit")
	flag.StringVar(&key, "k", "", "secret key for signing data")
	flag.Parse()

	if envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envReportInterval != "" {
		reportInterval, err = strconv.Atoi(envReportInterval)
		if err != nil {
			return nil, errors.Wrap(err, "error parse report interval")
		}
	}
	if envPollInterval != "" {
		pollInterval, err = strconv.Atoi(envPollInterval)
		if err != nil {
			return nil, errors.Wrap(err, "error parse poll interval")
		}
	}
	if envKey != "" {
		key = envKey
	}
	if envRateLimit != "" {
		rateLimit, err = strconv.Atoi(envRateLimit)
		if err != nil {
			return nil, errors.Wrap(err, "error parse rate limit")
		}
	}

	return &Config{
		ServerAddress:  serverAddress,
		PollInterval:   time.Duration(pollInterval) * time.Second,
		ReportInterval: time.Duration(reportInterval) * time.Second,
		Key:            key,
		RateLimit:      rateLimit,
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
	}, nil
}

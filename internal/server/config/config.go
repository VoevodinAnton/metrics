package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/VoevodinAnton/metrics/pkg/config"
)

const (
	defaultStoreInterval = 300
)

type Config struct {
	Logger        *config.Logger
	Server        Server
	FilePath      string
	StoreInterval time.Duration
	Restore       bool
}

type Server struct {
	Address string
}

func InitConfig() *Config {
	var serverAddress string
	var storeInterval int
	var restore bool
	var filePath string

	envServerAddress := os.Getenv("ADDRESS")
	envStoreInterval := os.Getenv("STORE_INTERVAL")
	envFilePath := os.Getenv("FILE_STORAGE_PATH")
	envRestore := os.Getenv("RESTORE")

	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&storeInterval, "i", defaultStoreInterval, "Interval in seconds to save metrics to disk")
	flag.StringVar(&filePath, "f", "/tmp/metrics-db.json", "Path to file where metrics are saved")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file on start")
	flag.Parse()

	if envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envStoreInterval != "" {
		storeInterval, _ = strconv.Atoi(envStoreInterval)
	}
	if envFilePath != "" {
		filePath = envFilePath
	}
	if envRestore != "" {
		restore, _ = strconv.ParseBool(envRestore)
	}

	cfg := Config{
		Server: Server{
			Address: serverAddress,
		},
		Logger: &config.Logger{
			Development: true,
			Level:       "debug",
		},
		StoreInterval: time.Duration(storeInterval) * time.Second,
		FilePath:      filePath,
		Restore:       restore,
	}

	return &cfg
}

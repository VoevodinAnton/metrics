package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/VoevodinAnton/metrics/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Reader microservice config path")
}

const (
	defaultStoreInterval = 300

	configPathEnv      = "CONFIG_PATH"
	serverAddressEnv   = "ADDRESS"
	storeIntervalEnv   = "STORE_INTERVAL"
	fileStoragePassEnv = "FILE_STORAGE_PATH"
	restoreEnv         = "RESTORE"
	databaseDSNEnv     = "DATABASE_DSN"

	yaml = "yaml"
)

type Config struct {
	Logger        *config.Logger `mapstructure:"logger"`
	Postgres      *config.Postgres
	Server        *config.Server
	FilePath      string
	StoreInterval time.Duration
	Restore       bool
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(configPathEnv)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/internal/server/config/config.yml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	var serverAddress string
	var storeInterval int
	var restore bool
	var filePath string
	var databaseDSN string

	envServerAddress := os.Getenv(serverAddress)
	envStoreInterval := os.Getenv(storeIntervalEnv)
	envFilePath := os.Getenv(fileStoragePassEnv)
	envRestore := os.Getenv(restoreEnv)
	envDatabaseDSN := os.Getenv(databaseDSNEnv)

	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&storeInterval, "i", defaultStoreInterval, "Interval in seconds to save metrics to disk")
	flag.StringVar(&filePath, "f", "/tmp/metrics-db.json", "Path to file where metrics are saved")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file on start")
	flag.StringVar(&databaseDSN, "d", "", "Connection string to postgres")
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
	if envDatabaseDSN != "" {
		databaseDSN = envDatabaseDSN
	}
	fmt.Println(databaseDSN)
	cfg.Server = &config.Server{
		Address: serverAddress,
	}
	cfg.StoreInterval = time.Duration(storeInterval) * time.Second
	cfg.FilePath = filePath
	cfg.Restore = restore
	cfg.Postgres = &config.Postgres{
		DatabaseDSN: databaseDSN,
	}

	return cfg, nil
}

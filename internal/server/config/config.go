package config

import (
	"flag"
	"os"

	"github.com/VoevodinAnton/metrics/pkg/config"
)

type Config struct {
	Logger *config.Logger
	Server Server
}

type Server struct {
	Address string
}

func InitConfig() *Config {
	var serverAddress string
	flag.StringVar(&serverAddress, "a", "localhost:8080", "HTTP server endpoint address")
	flag.Parse()

	envServerAddress := os.Getenv("ADDRESS")
	if envServerAddress != "" {
		serverAddress = envServerAddress
	}

	cfg := Config{
		Server: Server{
			Address: serverAddress,
		},
		Logger: &config.Logger{
			Development: true,
			Level:       "debug",
		},
	}

	return &cfg
}

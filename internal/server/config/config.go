package config

import (
	"flag"
	"os"
)

type Config struct {
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
	}

	return &cfg
}

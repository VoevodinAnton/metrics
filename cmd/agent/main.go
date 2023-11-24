package main

import (
	"github.com/VoevodinAnton/metrics/internal/agent/config"
	"github.com/VoevodinAnton/metrics/internal/agent/core/collector"
	"github.com/VoevodinAnton/metrics/internal/agent/core/uploader"
	logger "github.com/VoevodinAnton/metrics/pkg/logging"
)

func main() {
	cfg := config.InitConfig()

	logger.NewLogger(cfg.Logger)
	defer logger.Close()

	c := collector.NewCollector(cfg)
	u := uploader.NewUploader(cfg, c)

	go c.Run()
	u.Run()
}

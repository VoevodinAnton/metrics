package logging

import (
	"fmt"

	"github.com/VoevodinAnton/metrics/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOpts func(*zap.Logger)

func NewLogger(config *config.Logger, opts ...LoggerOpts) {
	logLevel, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		panic(fmt.Sprintf("Unknown log level %v", logLevel))
	}

	var cfg zap.Config
	if config.Development {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	for _, opt := range opts {
		opt(logger)
	}

	zap.ReplaceGlobals(logger)
}

func WithLoggerName(name string) LoggerOpts {
	return func(l *zap.Logger) {
		l.Named(name)
	}
}

func WithOptions(opts ...zap.Option) LoggerOpts {
	return func(l *zap.Logger) {
		l.WithOptions(opts...)
	}
}

func WithHooks(hooks ...func(e zapcore.Entry) error) LoggerOpts {
	return func(l *zap.Logger) {
		l.WithOptions(zap.AddCaller(), zap.Hooks(hooks...))
	}
}

func Close() {
	defer func() {
		_ = zap.L().Sync()
	}()
}

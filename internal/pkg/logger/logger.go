package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogMod string

const (
	DevelopmentMod LogMod = "development"
	ProductionMod  LogMod = "production"
)

type Config struct {
	LogMod   LogMod `mapstructure:"mod"` // LOGGER_MOD_PRODUCTION
	LogLevel string `mapstructure:"level"`
}

var globalLogger = newDefault() //nolint:gochecknoglobals // global logger needed for all app.

// newDefault creates new default logger.
func newDefault(opts ...zap.Option) *zap.Logger {
	cfg := newZapCfg(DevelopmentMod, zapcore.DebugLevel)
	logger, _ := cfg.Build(opts...)

	return logger
}

// newZapCfg creates new zap config.
func newZapCfg(mod LogMod, logLevel zapcore.Level) zap.Config {
	var cfg zap.Config

	switch mod {
	case ProductionMod:
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(logLevel)
	case DevelopmentMod:
		cfg = zap.NewDevelopmentConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	return cfg
}

// NewFromConfig creates new logger from config.
func NewFromConfig(cfg *Config, opts ...zap.Option) *zap.Logger {
	level, parseErr := zapcore.ParseLevel(cfg.LogLevel)
	if parseErr != nil {
		panic(parseErr)
	}

	zapCfg := newZapCfg(cfg.LogMod, level)
	logger, _ := zapCfg.Build(opts...)

	globalLogger = logger

	return globalLogger
}

// Global returns the global logger.
func Global() *zap.Logger {
	return globalLogger
}

// StructLogger is an alias for *zap.Logger included in project struct.
type StructLogger = *zap.Logger

// NewStructLogger returns a new StructLogger with the given name.
func NewStructLogger(name string) StructLogger {
	return globalLogger.Named(name)
}

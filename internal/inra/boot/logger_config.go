package boot

import (
	"fmt"
	"log/slog"
)

type LoggerConfig struct {
	Level          string `env:"LOG_LEVEL" envDefault:"info"`
	Pretty         bool   `env:"LOG_PRETTY"`
	IncludeSources bool   `env:"LOG_SOURCES"`
}

const (
	LoggerDebugLevel = "debug"
	LoggerInfoLevel  = "info"
	LoggerWarnLevel  = "warn"
	LoggerErrorLevel = "error"
)

func (r *LoggerConfig) Validate() error {
	switch r.Level {
	case LoggerDebugLevel, LoggerInfoLevel, LoggerWarnLevel, LoggerErrorLevel:
		return nil
	default:
		return fmt.Errorf("Failed to parse log level %s, unknown log level", r.Level)
	}
}

func (r *LoggerConfig) SlogLevel() slog.Level {
	switch r.Level {
	case LoggerDebugLevel:
		return slog.LevelDebug
	case LoggerInfoLevel:
		return slog.LevelInfo
	case LoggerWarnLevel:
		return slog.LevelWarn
	case LoggerErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

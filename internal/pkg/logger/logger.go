package logger

import (
	"doctormakarhina/lumos/internal/pkg/errs"
	"log/slog"
	"os"
)

const DebugLevel = "debug"
const InfoLevel = "info"
const WarnLevel = "warn"
const ErrorLevel = "error"

type Config struct {
	Level          string
	Pretty         bool
	IncludeSources bool
}

func (r *Config) validate() error {
	switch r.Level {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel:
		return nil
	default:
		return errs.NewErrorf(errs.ErrCodeInvalidArgument, "Failed to parse log level %s, unknown log level", r.Level)
	}
}

func (r *Config) slogLevel() slog.Level {
	switch r.Level {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func ConfigureDefault() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}

func NewRoot(config Config) (*slog.Logger, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	var handler slog.Handler
	handlerConfig := &slog.HandlerOptions{
		AddSource: config.IncludeSources,
		Level:     config.slogLevel(),
	}
	if config.Pretty {
		handler = slog.NewTextHandler(os.Stdout, handlerConfig)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, handlerConfig)
	}
	return slog.New(handler), nil
}

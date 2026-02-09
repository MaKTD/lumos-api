package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func ConfigureDefault() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func SetToDefault(logger *slog.Logger) {
	slog.SetDefault(logger)
}

func New(
	level slog.Level,
	pretty bool,
	includeSources bool,
	slogArgs ...any,
) (*slog.Logger, error) {
	var handler slog.Handler
	if pretty {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource: includeSources,
			Level:     level,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: includeSources,
			Level:     level,
		})
	}
	return slog.New(handler).With(slogArgs...), nil
}

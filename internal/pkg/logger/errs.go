package logger

import (
	"doctormakarhina/lumos/internal/pkg/errs"
	"errors"
	"log/slog"
)

func LogForHandler(
	message string,
	err error,
	logger *slog.Logger,
) {
	if err == nil {
		return
	}

	var codeErr *errs.CodeError
	if ok := errors.As(err, &codeErr); !ok {
		logger.Error(
			message,
			slog.String("err", err.Error()),
			slog.String("code", "unspecified"),
		)
	}

	switch codeErr.Code() {
	case errs.ErrCodeUnknown, errs.ErrCodeInternal:
		logger.Error(
			message,
			slog.String("err", codeErr.Error()),
			slog.String("code", errs.CodeToString(codeErr.Code())),
		)
	case errs.ErrCodeTimeout:
		logger.Warn(
			message,
			slog.String("err", codeErr.Error()),
			slog.String("code", errs.CodeToString(codeErr.Code())),
		)
	default:
		logger.Debug(
			message,
			slog.String("err", codeErr.Error()),
			slog.String("code", errs.CodeToString(codeErr.Code())),
		)
	}
}

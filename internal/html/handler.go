package html

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type GlobalHandler struct {
	logger *slog.Logger
}

func NewGlobalHandler(
	rootLogger *slog.Logger,
) *GlobalHandler {
	return &GlobalHandler{
		logger: rootLogger.With(slog.String("context", "htmlGlobalHandle")),
	}
}

func (h *GlobalHandler) RegIn(router chi.Router) {

}

package api

import (
	"doctormakarhina/lumos/internal/pkg/errs"
	"doctormakarhina/lumos/internal/pkg/modules/auth_logs"
	"doctormakarhina/lumos/internal/pkg/modules/search_logs"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type GlobalHandler struct {
	logger            *slog.Logger
	authLogsService   *auth_logs.Service
	searchLogsService *search_logs.Service
}

func NewGlobalHandler(
	rootLogger *slog.Logger,
	authLogsService *auth_logs.Service,
	queryLogsService *search_logs.Service,
) *GlobalHandler {
	return &GlobalHandler{
		authLogsService:   authLogsService,
		searchLogsService: queryLogsService,
		logger:            rootLogger.With(slog.String("context", "apiGlobalHandler")),
	}
}

func (h *GlobalHandler) RegIn(router chi.Router) {
	router.Post("/auth/telemetry", h.SaveAuthTelemetry)
	router.Post("/search/telemetry", h.SaveSearchTelemetry)
}

func (h *GlobalHandler) SaveAuthTelemetry(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	var req SaveAuthTelemetryRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errs.WrapErrorf(err, errs.ErrCodeParsingFailed, "failed to parse request body")
		h.logger.Error(
			"failed to parse request",
			slog.String("err", err.Error()),
		)

		return
	}

	err := h.authLogsService.SaveAuthTelemetry(
		r.Context(),
		req.Login,
		r.UserAgent(),
		r.RemoteAddr,
		req.Fingerprint,
		req.ConfidenceScore,
	)
	if err != nil {
		h.logger.Error(
			"failed to SaveAuthTelemetry",
			slog.String("err", err.Error()),
		)
	}
}

func (h *GlobalHandler) SaveSearchTelemetry(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	var req SaveSearchTelemetryRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errs.WrapErrorf(err, errs.ErrCodeParsingFailed, "failed to parse request body")
		h.logger.Error(
			"failed to parse request",
			slog.String("err", err.Error()),
		)

		return
	}

	err := h.searchLogsService.SaveLogsQuery(r.Context(), req.Query)
	if err != nil {
		h.logger.Error(
			"failed to SaveLogsQuery",
			slog.String("err", err.Error()),
		)
	}
}

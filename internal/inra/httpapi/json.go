package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeJSON(w http.ResponseWriter, logger *slog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("writeJSON error", slog.String("err", err.Error()))
	}
}

func writeJsonError(w http.ResponseWriter, logger *slog.Logger, status int, msg string) {
	writeJSON(w, logger, status, ErrorRes{Error: msg})
}

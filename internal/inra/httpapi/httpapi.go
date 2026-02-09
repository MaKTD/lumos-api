package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return r
}

func RegInPing(r chi.Router) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func RegInHealthz(r chi.Router, logger *slog.Logger) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, logger, http.StatusOK, HealthRes{Status: "ok"})
	})
}

func RegInAuthLogs(r chi.Router, db *sqlx.DB, rootLogger *slog.Logger) {
	srv := &authLogs{db: db, logger: rootLogger.With(slog.String("context", "AuthLogs"))}

	r.Post("/auth/telemetry", srv.Handle)
}

func RegInSearchLogs(r chi.Router, db *sqlx.DB, rootLogger *slog.Logger) {
	srv := &searchLogs{db: db, logger: rootLogger.With(slog.String("context", "SearchLogs"))}

	r.Post("/search/telemetry", srv.Handle)
}

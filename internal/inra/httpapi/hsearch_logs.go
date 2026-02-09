package httpapi

import (
	"context"
	"doctormakarhina/lumos/internal/pkg/errs"
	"encoding/json"
	"log/slog"
	"net/http"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
)

type searchLogs struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func (s *searchLogs) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	var req SaveSearchTelemetryRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errs.WrapErrorf(err, errs.ErrCodeParsingFailed, "failed to parse request body")
		s.logger.Error(
			"failed to parse request",
			slog.String("err", err.Error()),
		)

		return
	}

	err := s.saveLogsQuery(r.Context(), req.Query)
	if err != nil {
		s.logger.Error(
			"failed to SaveLogsQuery",
			slog.String("err", err.Error()),
		)
	}
}

func (s *searchLogs) saveLogsQuery(
	ctx context.Context,
	searchQuery string,
) error {
	if searchQuery == "" {
		return nil
	}
	if utf8.RuneCountInString(searchQuery) <= 2 {
		return nil
	}

	query := `INSERT INTO lumos.search_queries (query) VALUES ($1)`

	_, err := s.db.ExecContext(ctx, query, searchQuery)

	return err
}

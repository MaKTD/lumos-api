package search_logs

import (
	"context"
	"database/sql"
	"unicode/utf8"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (r *Service) SaveLogsQuery(
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

	_, err := r.db.ExecContext(ctx, query, searchQuery)

	return err
}

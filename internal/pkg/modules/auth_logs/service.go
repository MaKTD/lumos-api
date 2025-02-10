package auth_logs

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (r *Service) SaveAuthTelemetry(
	ctx context.Context,
	login string,
	userAgent string,
	ipAddrRaw string,
	fingerprintId string,
	fingerprintScore string,
) error {
	if login == "" {
		return nil
	}

	ipAddr := "unknown"
	if ipAddrRaw != "" {
		parts := strings.Split(ipAddrRaw, ":")
		if len(parts) > 0 {
			ipAddr = parts[0]
		} else {
			ipAddr = ipAddrRaw
		}
	}

	fmt.Println(login, ipAddr, fingerprintId, fingerprintScore, userAgent)

	query := `
		INSERT INTO lumos.auth_logs
		(login, ip, fingerprint, confidencescore, useragent)
		VALUES 
		($1, $2, $3, $4, $5);
    `

	_, err := r.db.ExecContext(ctx, query, login, ipAddr, fingerprintId, fingerprintScore, userAgent)

	return err
}

package httpapi

import (
	"context"
	"doctormakarhina/lumos/internal/pkg/errs"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type authLogs struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func (s *authLogs) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	var req SaveAuthTelemetryRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errs.WrapErrorf(err, errs.ErrCodeParsingFailed, "failed to parse request body")
		s.logger.Error(
			"failed to parse request",
			slog.String("err", err.Error()),
		)

		return
	}

	err := s.saveAuthTelemetry(
		r.Context(),
		req.Login,
		r.UserAgent(),
		r.RemoteAddr,
		req.Fingerprint,
		req.ConfidenceScore,
	)

	if err != nil {
		s.logger.Error(
			"failed to SaveAuthTelemetry",
			slog.String("err", err.Error()),
		)
	}
}

func (s *authLogs) saveAuthTelemetry(
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

	query := `
		INSERT INTO lumos.auth_logs
		(login, ip, fingerprint, confidencescore, useragent)
		VALUES
		($1, $2, $3, $4, $5);
    `

	_, err := s.db.ExecContext(ctx, query, login, ipAddr, fingerprintId, fingerprintScore, userAgent)

	return err
}

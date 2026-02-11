package httpapi

import (
	"doctormakarhina/lumos/internal/core/payments"
	"log/slog"
	"net/http"
)

type userAccess struct {
	srv    payments.Service
	logger *slog.Logger
}

func (s *userAccess) Handle(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	projectID := r.URL.Query().Get("project_id")

	if email == "" || projectID == "" {
		writeJSON(w, s.logger, 400, ErrMsgRes{Message: "missing email or project_id in search params"})
		return
	}

	allowed, err := s.srv.IsAccessAlowed(r.Context(), email, projectID)
	if err != nil {
		writeJSON(w, s.logger, 500, ErrMsgRes{Message: "internal error"})
		return
	}

	if !allowed {
		writeJSON(w, s.logger, 403, ErrMsgRes{Message: "forbidden"})
		return
	}

	writeJSON(w, s.logger, 200, UserAccessRes{Success: true})
}

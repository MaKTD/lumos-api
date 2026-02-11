package httpapi

import (
	"doctormakarhina/lumos/internal/core/payments"
	"errors"
	"log/slog"
	"net/http"
)

type userInfo struct {
	srv    payments.Service
	logger *slog.Logger
}

func (s *userInfo) Handle(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	projectID := r.URL.Query().Get("project_id")

	if email == "" || projectID == "" {
		writeJSON(w, s.logger, 400, ErrMsgRes{Message: "missing email or project_id in search params"})
		return
	}

	user, err := s.srv.User(r.Context(), email, projectID)
	if err != nil {
		if errors.Is(err, payments.ErrInvalidProjectId) {
			writeJSON(w, s.logger, 403, ErrMsgRes{Message: "forbidden"})
			return
		}

		writeJSON(w, s.logger, 500, ErrMsgRes{Message: "internal server error"})
		return
	}

	if user == nil {
		writeJSON(w, s.logger, 404, ErrMsgRes{Message: "not found"})
		return
	}

	writeJSON(w, s.logger, 200, UserInfoRes{
		Email:           user.Email,
		Tariff:          user.Tariff,
		ExpiresAt:       user.ExpiresAt,
		SubStatus:       user.SubscriptionStatus,
		SubscriptionsId: user.SubscriptionID,
		NextPaymentSum:  user.LastSubPrice,
	})
}

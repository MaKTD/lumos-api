package httpapi

import (
	"doctormakarhina/lumos/internal/core/notify"
	"doctormakarhina/lumos/internal/core/payments"
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

func RegInTrialPayments(
	r chi.Router,
	routeHash string,
	srv payments.Service,
	notifer notify.Service,

) {
	trialSrv := paymentRegTrial{srv: srv, notifer: notifer}

	r.Post("/payments/trial/"+routeHash, trialSrv.Handle)
}

func RegInProdamusPayWebHook(
	r chi.Router,
	routeHash string,
	srv payments.Service,
	notifer notify.Service,
	logger *slog.Logger,
) {
	prodamusSrv := prodamusPayNotification{
		srv:     srv,
		notifer: notifer,
		logger:  logger.With(slog.String("context", "ProdamusPayNotificationHandler")),
	}

	r.Post("/payments/prodamus/webhook/pay/"+routeHash, prodamusSrv.Handle)
}

func RegInCloudPaymentsPayHook(
	r chi.Router,
	routeHash string,
	srv payments.Service,
	notifer notify.Service,
	logger *slog.Logger,
) {
	cpSrv := cloudPaymentsPayNotification{
		srv:     srv,
		notifer: notifer,
		logger:  logger.With(slog.String("context", "CloudPaymentsPayNotificationHandler")),
	}

	r.Post("/payments/cloudpayment/webhook/pay/"+routeHash, cpSrv.Handle)
}

func RegInCloudPaymentsReccurentNotif(
	r chi.Router,
	routeHash string,
	srv payments.Service,
	notifer notify.Service,
	logger *slog.Logger,
) {
	cpSrv := cloudPaymentReccurentNotif{
		srv:     srv,
		notifer: notifer,
		logger:  logger.With(slog.String("context", "CloudPaymentsReccurentNotificationHandler")),
	}

	r.Post("/payments/cloudpayment/webhook/reccurent/"+routeHash, cpSrv.Handle)
}

func RegInUserAccessRoute(
	r chi.Router,
	srv payments.Service,
	logger *slog.Logger,
) {
	userAccessSrv := userAccess{srv: srv, logger: logger.With(slog.String("context", "UserAccessHandler"))}

	r.Get("/payments/user/access", userAccessSrv.Handle)
}

func RegInUserInfoRoute(
	r chi.Router,
	srv payments.Service,
	logger *slog.Logger,
) {
	userInfoSrv := userInfo{srv: srv, logger: logger.With(slog.String("context", "UserInfoHandler"))}

	r.Get("/payments/user/info", userInfoSrv.Handle)
}

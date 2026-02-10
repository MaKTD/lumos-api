package httpapi

import (
	"doctormakarhina/lumos/internal/core/notify"
	"doctormakarhina/lumos/internal/core/payments"
	"fmt"
	"log/slog"
	"net/http"
)

type cloudPaymentReccurentNotif struct {
	logger  *slog.Logger
	srv     payments.Service
	notifer notify.Service
}

func (s *cloudPaymentReccurentNotif) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.notifer.ForAdmin("[CloudPaymentsPayHandler] recieve invalid form")
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	subscriptionID := r.FormValue("Id")
	email := r.FormValue("Email")
	status := r.FormValue("Status")

	if subscriptionID == "" || email == "" {
		s.notifer.ForAdmin(fmt.Sprintf("[cloudPaymentReccurentNotif] recieve form without required fields: subscriptionID=%s, email=%s", subscriptionID, email))
		writeJSON(w, s.logger, 200, CloudPaymentsNotificationRes{Code: 0})
		return
	}

	err := s.srv.RegisterCloudPaymentReccurent(r.Context(), subscriptionID, email, status)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, s.logger, 200, CloudPaymentsNotificationRes{Code: 0})

	//    Id caflksdjfklasjfaskl0
	//    AccountId tima98tima@gmail.com
	//    Description Доступ на месяц
	//    Email tima98tima@gmail.com
	//    Amount 20.00
	//    Currency RUB
	//    RequireConfirmation 0
	//    StartDate 2026-04-05 09:59:34
	//    Interval Month
	//    Period 1
	//    Status Cancelled
	//    SuccessfulTransactionsNumber 0
	//    FailedTransactionsNumber 0
	//    NextTransactionDate 2026-04-05 09:59:34
}

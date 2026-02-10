package httpapi

import (
	"doctormakarhina/lumos/internal/core/notify"
	"doctormakarhina/lumos/internal/core/payments"
	"fmt"
	"log/slog"
	"net/http"
)

type cloudPaymentsPayNotification struct {
	logger  *slog.Logger
	srv     payments.Service
	notifer notify.Service
}

func (s *cloudPaymentsPayNotification) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.notifer.ForAdmin("[CloudPaymentsPayHandler] recieve invalid form")
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	tariff := r.FormValue("Description")
	email := r.FormValue("Email")
	name := r.FormValue("Name")
	priceStr := r.FormValue("Amount")
	subscriptionID := r.FormValue("SubscriptionId")
	opType := r.FormValue("OperationType")

	if opType != "Payment" {
		s.notifer.ForAdmin(fmt.Sprintf("[CloudPaymentsPayHandler] recieve payment notification with invalid operation type = %s. email = %s, tariff = %s", opType, email, tariff))
		writeJSON(w, s.logger, 200, CloudPaymentsNotificationRes{Code: 0})
		return
	}

	if email == "" || tariff == "" {
		s.notifer.ForAdmin(fmt.Sprintf("[CloudPaymentsPayHandler] recieve payment notification with empty required fields email = %s, tariff = %s", email, tariff))
		writeJSON(w, s.logger, 200, CloudPaymentsNotificationRes{Code: 0})
		return
	}

	priceParsed, err := parsePrice(priceStr)
	if err != nil {
		s.notifer.ForAdmin(fmt.Sprintf("[CloudPaymentsPayHandler] recieve payment notification with invalid price = %s. email = %s, tariff = %s", priceStr, email, tariff))
		writeJSON(w, s.logger, 200, CloudPaymentsNotificationRes{Code: 0})
		return
	}

	err = s.srv.RegisterFromCloudPayments(
		r.Context(),
		tariff,
		email,
		name,
		float32(priceParsed),
		subscriptionID,
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, s.logger, 200, CloudPaymentsNotificationRes{Code: 0})

	//  TransactionId 2147332467
	// Amount 20.00
	// Currency RUB
	// PaymentAmount 20.00
	// PaymentCurrency RUB
	// OperationType Payment
	// InvoiceId 1707856125
	// AccountId tima98tima@gmail.com
	// SubscriptionId sc_25239b18a21e921468455cba5aa4f
	// Name
	// Email tima98tima@gmail.com
	// DateTime 2026-02-10 13:46:43
	// IpAddress 34.76.167.208
	// IpCountry US
	// IpCity Маунтин-Вью
	// IpRegion Калифорния
	// IpDistrict Маунтин-Вью
	// IpLatitude 37.38605
	// IpLongitude -122.08385
	// CardId
	// CardFirstSix 220070
	// CardLastFour 4560
	// CardType MIR
	// CardExpDate 09/33
	// Issuer T-Bank (Tinkoff)
	// IssuerBankCountry RU
	// Description Доступ на месяц
	// AuthCode 194364
	// Token tk_ec56e535818b5cd81b7ccfc114476
	// TestMode 0
	// Status Completed
	// GatewayName Tbank
	// DataLong String
	// TotalFee 3.90
	// CardProduct TKN
	// PaymentMethod TinkoffPay
	// Rrn 129057103428
	//InstallmentTerm InstallmentMonthlyPayment
}

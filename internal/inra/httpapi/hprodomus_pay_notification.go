package httpapi

import (
	"doctormakarhina/lumos/internal/core/notify"
	"doctormakarhina/lumos/internal/core/payments"
	"fmt"
	"log/slog"
	"net/http"
)

type prodamusPayNotification struct {
	logger  *slog.Logger
	srv     payments.Service
	notifer notify.Service
}

func (s *prodamusPayNotification) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.notifer.ForAdmin("[ProdamusPayHandler] recieve invalid form")
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	tariff := r.FormValue("customer_extra")
	email := r.FormValue("customer_email")
	// phone := r.FormValue("customer_phone")
	price := r.FormValue("sum")
	name := r.FormValue("_param_name")
	paymentStatus := r.FormValue("payment_status")

	if paymentStatus != "success" {
		s.notifer.ForAdmin(fmt.Sprintf("[ProdamusPayHandler] recieve payment notification with not success status = %s, email = %s, tariff = %s, price = %s, name = %s", paymentStatus, email, tariff, price, name))
		writeJSON(w, s.logger, 200, ProdamusPayNotificationRes{Success: true})
		return
	}

	if email == "" || tariff == "" {
		s.notifer.ForAdmin(fmt.Sprintf("[ProdamusPayHandler] recieve payment notification with empty required fields email = %s, tariff = %s", email, tariff))
		writeJSON(w, s.logger, 200, ProdamusPayNotificationRes{Success: true})
		return
	}

	priceParsed, err := parsePrice(price)
	if err != nil {
		s.notifer.ForAdmin(fmt.Sprintf("[ProdamusPayHandler] recieve payment notification with invalid price = %s. email = %s, tariff = %s", price, email, tariff))
		writeJSON(w, s.logger, 200, ProdamusPayNotificationRes{Success: true})
		return
	}

	err = s.srv.RegisterFromProdamus(
		r.Context(),
		tariff,
		email,
		name,
		float32(priceParsed),
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, s.logger, 200, ProdamusPayNotificationRes{Success: true})

	// date 2026-02-07T20:27:09+03:00
	//    order_id 41227995
	//    order_num 9313321:1224829704
	//    domain makarshina.payform.ru
	//    sum 2290.00
	//    currency eur
	//    currency_sum 31.97
	//    currency_commission_sum 3.20
	//    customer_phone +375298712954
	//    customer_email lubov.5@mail.ru
	//    customer_extra Продление 3 месяца
	//    payment_type Карты банков мира кроме России
	//    commission 10
	//    commission_sum 229.00
	//    attempt 1
	//    sys tilda
	//    _param_name Любовь
	//    productsArray 1Collection
	//            name Продление 3 месяца
	//            price 2290.00
	//            quantity 1
	//            sum 2290.00
	//    payment_status success
	//    payment_status_description Успешная оплата
	//    payment_init manual
}

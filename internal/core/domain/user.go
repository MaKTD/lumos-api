package domain

import (
	"fmt"
	"time"
)

const (
	UserTariffUnlimited string = "Бессрочный"
	UserTariffTrial     string = "Пробный период"
	UserTariff1Month    string = "1 месяц"
	UserTariff3Months   string = "3 месяца"
	UserTariff6Months   string = "6 месяцев"
)

const (
	statusOk      string = "🆗"
	statusExpired string = "⚠️"
)

const (
	UserSubStatusActive   string = "Active"
	UserSubStatusCanceled string = "Cancelled"
	UserSubStatusRejected string = "Rejected"
	UserSubStatusPastDue  string = "PastDue"
)

type User struct {
	ID                 string    `db:"id" json:"id"`
	Email              string    `db:"email" json:"email"`
	Name               string    `db:"name" json:"name"`
	Tariff             string    `db:"tariff" json:"tariff"`
	ExpiresAt          time.Time `db:"expires_at" json:"expires_at"`
	SubscriptionID     string    `db:"subscription_id" json:"subscription_id"`
	SubscriptionStatus string    `db:"subscription_status" json:"subscription_status"`
	LastSubPrice       float32   `db:"last_sub_price" json:"last_sub_price"`
}

func (u *User) SubExpired(now time.Time) bool {
	if u.Tariff == UserTariffUnlimited {
		return false
	}

	if u.Tariff == UserTariffTrial {
		return now.After(u.ExpiresAt)
	}

	return now.After(u.ExpiresAt)
}

func (u *User) StatusInfo(now time.Time) string {
	if u.Tariff == UserTariffUnlimited {
		return "Тариф: " + u.Tariff
	}

	status := statusOk
	if u.SubExpired(now) {
		status = statusExpired
	}

	subStatus := "Подписка отменена"
	if u.SubscriptionStatus == UserSubStatusActive {
		subStatus = "Подписка активна"
	}

	return "Тариф: " + u.Tariff +
		"\n" + "Оплачено до: " + u.ExpiresAt.Format("02.01.2006 ") + " " + status +
		"\n" + subStatus +
		"\n" + "Сумма следующего платежа: " + fmt.Sprintf("%.2f", u.LastSubPrice)
}

// redirect url
// IF(
//   IsExpired = "false",
//   "<script></script>",
//   "<script>window.location.replace('/subscription/payment')</script>"
// )
//
//

// cancel subscription url
//"https://my.cloudpayments.ru/unsubscribe/from/" & SubscriptionId

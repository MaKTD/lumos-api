package httpapi

import "time"

type ErrorRes struct {
	Error string `json:"error"`
}

type HealthRes struct {
	Status string `json:"status"`
}

type SaveAuthTelemetryRequestBody struct {
	Login           string `json:"login"`
	Fingerprint     string `json:"fingerprint"`
	ConfidenceScore string `json:"confidenceScore"`
}

type SaveSearchTelemetryRequestBody struct {
	Query string `json:"query"`
}

type ProdamusPayNotificationRes struct {
	Success bool `json:"success"`
}

type CloudPaymentsNotificationRes struct {
	Code int `json:"code"`
}

type UserAccessRes struct {
	Success bool `json:"success"`
}

type UserInfoRes struct {
	Email           string    `json:"email"`
	Tariff          string    `json:"tariff"`
	ExpiresAt       time.Time `json:"expires_at"`
	SubStatus       string    `json:"sub_status"`
	NextPaymentSum  float32   `json:"next_payment_sum"`
	SubscriptionsId string    `json:"sub_id"`
}

type ErrMsgRes struct {
	Message string `json:"message"`
}

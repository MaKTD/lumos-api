package httpapi

// import "lumos/search/internal/core/domain"

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

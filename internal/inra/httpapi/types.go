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

// type ProductSearchRes struct {
// 	Products []domain.Product `json:"products"`
// 	Total    int              `json:"total"`
// }

package api

type SaveAuthTelemetryRequestBody struct {
	Login           string `json:"login"`
	Fingerprint     string `json:"fingerprint"`
	ConfidenceScore string `json:"confidenceScore"`
}

type SaveSearchTelemetryRequestBody struct {
	Query string `json:"query"`
}

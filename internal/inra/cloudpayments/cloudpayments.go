package cloudpayments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	corepayments "doctormakarhina/lumos/internal/core/payments"
)

var _ corepayments.CloudPayments = (*CloudPaymentsClient)(nil)

const (
	defaultBaseURL = "https://api.cloudpayments.ru"

	endpointSubscriptionsUpdate = "/subscriptions/update"
	endpointSubscriptionsCancel = "/subscriptions/cancel"
)

type Config struct {
	// PublicID is CloudPayments "Public ID" (used as Basic Auth username).
	PublicID string
	// APISecret is CloudPayments "API Secret" (used as Basic Auth password).
	APISecret string

	// BaseURL is CloudPayments API base URL. If empty, defaults to https://api.cloudpayments.ru
	BaseURL string

	// HTTPClient is used to make requests. If nil, http.DefaultClient is used.
	HTTPClient *http.Client
}

type CloudPaymentsClient struct {
	publicID  string
	apiSecret string
	baseURL   string
	http      *http.Client
}

func New(cfg Config) (*CloudPaymentsClient, error) {
	if strings.TrimSpace(cfg.PublicID) == "" {
		return nil, fmt.Errorf("cloudpayments: PublicID is required")
	}
	if strings.TrimSpace(cfg.APISecret) == "" {
		return nil, fmt.Errorf("cloudpayments: APISecret is required")
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &CloudPaymentsClient{
		publicID:  cfg.PublicID,
		apiSecret: cfg.APISecret,
		baseURL:   baseURL,
		http:      httpClient,
	}, nil
}

// UpdateSubscription updates an existing subscription schedule in CloudPayments and returns the subscription Status
// from the response Model (e.g. "Active").
func (c *CloudPaymentsClient) UpdateSubscription(ctx context.Context, ID string, startDate time.Time, interval string, period int) (string, error) {
	id := strings.TrimSpace(ID)
	if id == "" {
		return "", fmt.Errorf("cloudpayments: UpdateSubscription: ID is empty")
	}
	interval = strings.TrimSpace(interval)
	if interval == "" {
		return "", fmt.Errorf("cloudpayments: UpdateSubscription: interval is empty")
	}
	if period <= 0 {
		return "", fmt.Errorf("cloudpayments: UpdateSubscription: period must be > 0")
	}

	return c.updateBySubscriptionID(ctx, id, startDate, interval, period)
}

func (c *CloudPaymentsClient) CancelSubscription(ctx context.Context, ID string) error {
	id := strings.TrimSpace(ID)
	if id == "" {
		return fmt.Errorf("cloudpayments: CancelSubscription: ID is empty")
	}

	req := struct {
		Id string `json:"Id"`
	}{
		Id: id,
	}

	var resp apiResponse[json.RawMessage]
	if err := c.doJSON(ctx, http.MethodPost, endpointSubscriptionsCancel, req, &resp); err != nil {
		return err
	}
	if !resp.Success {
		return newAPIError(http.StatusOK, resp.Message, resp.rawBody)
	}

	return nil
}

func (c *CloudPaymentsClient) updateBySubscriptionID(ctx context.Context, subscriptionID string, startDate time.Time, interval string, period int) (string, error) {
	req := struct {
		Id        string `json:"Id"`
		StartDate string `json:"StartDate"`
		Interval  string `json:"Interval"`
		Period    int    `json:"Period"`
	}{
		Id:        subscriptionID,
		StartDate: startDate.UTC().Format(time.RFC3339),
		Interval:  interval,
		Period:    period,
	}

	var resp apiResponse[subscriptionModel]
	if err := c.doJSON(ctx, http.MethodPost, endpointSubscriptionsUpdate, req, &resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", newAPIError(http.StatusOK, resp.Message, resp.rawBody)
	}

	status := strings.TrimSpace(resp.Model.Status)
	if status == "" {
		// Be defensive: if Model is empty but success=true, return a sensible default.
		status = "Active"
	}
	return status, nil
}

type subscriptionModel struct {
	ID                     string `json:"Id"`
	AccountID              string `json:"AccountId"`
	Status                 string `json:"Status"`
	NextTransactionDate    string `json:"NextTransactionDate"`
	NextTransactionDateIso string `json:"NextTransactionDateIso"`
}

type apiResponse[T any] struct {
	Model   T       `json:"Model"`
	Success bool    `json:"Success"`
	Message *string `json:"Message"`

	// rawBody is captured by doJSON for richer errors.
	rawBody string `json:"-"`
}

type APIError struct {
	HTTPStatus int
	Message    string
	Raw        string
}

func (e *APIError) Error() string {
	msg := strings.TrimSpace(e.Message)
	if msg == "" {
		msg = "request failed"
	}
	if e.HTTPStatus > 0 {
		msg = fmt.Sprintf("cloudpayments: %s (http=%d)", msg, e.HTTPStatus)
	}
	if raw := strings.TrimSpace(e.Raw); raw != "" {
		// Keep this short; raw is mainly for diagnostics.
		const max = 512
		if len(raw) > max {
			raw = raw[:max] + "…"
		}
		msg = msg + ": " + raw
	}
	return msg
}

func newAPIError(httpStatus int, message *string, raw string) error {
	msg := ""
	if message != nil {
		msg = *message
	}
	return &APIError{
		HTTPStatus: httpStatus,
		Message:    msg,
		Raw:        raw,
	}
}

func (c *CloudPaymentsClient) doJSON(ctx context.Context, method, path string, req any, out any) error {
	if ctx == nil {
		return errors.New("cloudpayments: context is nil")
	}
	if strings.TrimSpace(method) == "" {
		return errors.New("cloudpayments: http method is empty")
	}

	url := c.baseURL + path

	var body io.Reader
	if req != nil {
		b, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("cloudpayments: failed to marshal request: %w", err)
		}
		body = bytes.NewReader(b)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("cloudpayments: failed to create request: %w", err)
	}

	httpReq.SetBasicAuth(c.publicID, c.apiSecret)
	httpReq.Header.Set("Accept", "application/json")
	if req != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	res, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("cloudpayments: request failed: %w", err)
	}
	defer res.Body.Close()

	rawBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("cloudpayments: failed to read response body: %w", readErr)
	}
	raw := strings.TrimSpace(string(rawBytes))

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// CloudPayments typically returns JSON even on errors, but we still surface raw.
		return &APIError{
			HTTPStatus: res.StatusCode,
			Message:    res.Status,
			Raw:        raw,
		}
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(rawBytes, out); err != nil {
		return fmt.Errorf("cloudpayments: failed to decode response: %w: %s", err, raw)
	}

	// If it's our apiResponse[*], attach raw body for better errors.
	switch v := out.(type) {
	case *apiResponse[json.RawMessage]:
		v.rawBody = raw
	case *apiResponse[subscriptionModel]:
		v.rawBody = raw
	case *apiResponse[[]subscriptionModel]:
		v.rawBody = raw
	default:
		// ignore
	}

	return nil
}

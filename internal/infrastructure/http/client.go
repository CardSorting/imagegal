package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"image/internal/domain/models"
	"image/internal/domain/ports"
	apperrors "image/pkg/errors"
)

// Client implements the HTTPClient interface
type Client struct {
	client     *http.Client
	baseURL    string
	apiKey     string
	maxRetries int
	logger     ports.Logger
}

// ClientOption defines a function type for client configuration
type ClientOption func(*Client)

// NewClient creates a new HTTP client with the provided options
func NewClient(baseURL, apiKey string, logger ports.Logger, opts ...ClientOption) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    baseURL,
		apiKey:     apiKey,
		maxRetries: 3,
		logger:     logger,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// Do executes an HTTP request with retries and error handling
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var attempt int

	for attempt = 1; attempt <= c.maxRetries; attempt++ {
		c.logger.Debug("Attempting request",
			"attempt", attempt,
			"url", req.URL.String(),
			"method", req.Method,
		)

		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt < c.maxRetries {
			delay := time.Duration(attempt) * time.Second
			time.Sleep(delay)
		}
	}

	if err != nil {
		return nil, apperrors.NewExternalAPIError(
			fmt.Sprintf("HTTP request failed after %d attempts", attempt-1),
			err,
		)
	}

	return resp, nil
}

// Get sends a GET request
func (c *Client) Get(ctx context.Context, path string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to create request", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Log request details
	c.logger.Debug("Preparing request",
		"url", req.URL.String(),
		"method", "GET",
	)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return apperrors.NewInternalServerError("Failed to decode response", err)
	}

	return nil
}

// Post sends a POST request with JSON body
func (c *Client) Post(ctx context.Context, path string, body interface{}, response interface{}) error {
	// Set API key in request body for ModelsLab API
	if bodyStruct, ok := body.(*models.ModelsLabAPIRequest); ok {
		bodyStruct.Key = c.apiKey
		c.logger.Debug("Setting API key for request",
			"key_length", len(c.apiKey),
			"path", path,
		)
	}

	// Add additional headers for Cloudflare
	headers := map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Origin":          "https://modelslab.com",
		"Referer":         "https://modelslab.com/",
		"Accept-Language": "en-US,en;q=0.9",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
	}

	// Log request details
	c.logger.Debug("Preparing request",
		"url", c.baseURL+path,
		"method", "POST",
	)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return apperrors.NewInvalidRequestError("Failed to marshal request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return apperrors.NewInternalServerError("Failed to create request", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Log request body for debugging
	c.logger.Debug("Request payload",
		"body", string(jsonBody),
		"url", req.URL.String(),
	)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return apperrors.NewInternalServerError("Failed to decode response", err)
	}

	return nil
}

// handleErrorResponse processes error responses from the API
func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	// Log raw error response for debugging
	c.logger.Debug("Received error response",
		"status_code", statusCode,
		"body", string(body),
	)

	// Try to parse as ModelsLab API error
	var errorResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Code    string `json:"code,omitempty"`
	}

	if err := json.Unmarshal(body, &errorResp); err != nil {
		// Check if it's a Cloudflare error page
		if bytes.Contains(body, []byte("cloudflare")) {
			return apperrors.NewExternalAPIError(
				"Request blocked by Cloudflare",
				fmt.Errorf("status code: %d, possible reasons: rate limiting, invalid headers, or bot detection", statusCode),
			)
		}

		return apperrors.NewExternalAPIError(
			fmt.Sprintf("API error with status code %d", statusCode),
			fmt.Errorf("raw response: %s", string(body)),
		)
	}

	// Log parsed error details
	c.logger.Debug("Parsed error response",
		"status", errorResp.Status,
		"message", errorResp.Message,
		"code", errorResp.Code,
	)

	switch statusCode {
	case http.StatusUnauthorized:
		return apperrors.NewUnauthorizedError("Invalid or missing API key", nil)
	case http.StatusBadRequest:
		if errorResp.Message != "" {
			return apperrors.NewInvalidRequestError(errorResp.Message, nil)
		}
		return apperrors.NewInvalidRequestError("Invalid request parameters", nil)
	case http.StatusTooManyRequests:
		return apperrors.NewExternalAPIError("Rate limit exceeded", nil)
	default:
		if errorResp.Message != "" {
			return apperrors.NewExternalAPIError(errorResp.Message, nil)
		}
		return apperrors.NewExternalAPIError("Unexpected API error", nil)
	}
}

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"image/internal/domain/ports"
	apperrors "image/pkg/errors"
)

// Client implements the HTTPClient interface
type Client struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	maxRetries int
	logger    ports.Logger
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

// Post sends a POST request with JSON body
func (c *Client) Post(ctx context.Context, path string, body interface{}, response interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return apperrors.NewInvalidRequestError("Failed to marshal request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return apperrors.NewInternalServerError("Failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
	var errorResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Code    string `json:"code,omitempty"`
	}

	if err := json.Unmarshal(body, &errorResp); err != nil {
		return apperrors.NewExternalAPIError(
			fmt.Sprintf("API error with status code %d", statusCode),
			fmt.Errorf("raw response: %s", string(body)),
		)
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return apperrors.NewUnauthorizedError(errorResp.Message, nil)
	case http.StatusBadRequest:
		return apperrors.NewInvalidRequestError(errorResp.Message, nil)
	default:
		return apperrors.NewExternalAPIError(errorResp.Message, nil)
	}
}

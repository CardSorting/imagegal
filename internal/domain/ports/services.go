package ports

import (
	"context"
	"net/http"

	"image/internal/domain/models"
)

// ModelsLabService defines the interface for interacting with the ModelsLab API
type ModelsLabService interface {
	// GenerateImage generates an image from text using the ModelsLab API
	GenerateImage(ctx context.Context, req *models.Text2ImgRequest) (*models.Text2ImgResponse, error)
}

// ImageGenerator defines the interface for image generation
type ImageGenerator interface {
	// Generate creates an image based on the provided parameters
	Generate(ctx context.Context, params *models.Text2ImgRequest) ([]string, error)
}

// Validator defines the interface for request validation
type Validator interface {
	// Validate validates the request parameters
	Validate(params interface{}) error
}

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	// Do executes an HTTP request and returns an HTTP response
	Do(req *http.Request) (*http.Response, error)
	// Post sends a POST request with JSON body
	Post(ctx context.Context, path string, body interface{}, response interface{}) error
	// Get sends a GET request
	Get(ctx context.Context, path string, response interface{}) error
}

// Logger defines the interface for logging operations
type Logger interface {
	// Info logs information messages
	Info(msg string, fields ...interface{})
	// Error logs error messages
	Error(msg string, err error, fields ...interface{})
	// Debug logs debug messages
	Debug(msg string, fields ...interface{})
}

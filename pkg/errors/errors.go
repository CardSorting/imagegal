package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a specific error type
type ErrorCode string

const (
	// ErrInvalidRequest represents validation or bad request errors
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
	// ErrUnauthorized represents authentication errors
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrInternalServer represents internal server errors
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	// ErrExternalAPI represents errors from external API calls
	ErrExternalAPI ErrorCode = "EXTERNAL_API_ERROR"
	// ErrTimeout represents timeout errors
	ErrTimeout ErrorCode = "TIMEOUT"
)

// AppError represents an application-specific error
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
	Status  int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewInvalidRequestError creates a new invalid request error
func NewInvalidRequestError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrInvalidRequest,
		Message: message,
		Err:     err,
		Status:  http.StatusBadRequest,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrUnauthorized,
		Message: message,
		Err:     err,
		Status:  http.StatusUnauthorized,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrInternalServer,
		Message: message,
		Err:     err,
		Status:  http.StatusInternalServerError,
	}
}

// NewExternalAPIError creates a new external API error
func NewExternalAPIError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrExternalAPI,
		Message: message,
		Err:     err,
		Status:  http.StatusBadGateway,
	}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrTimeout,
		Message: message,
		Err:     err,
		Status:  http.StatusGatewayTimeout,
	}
}

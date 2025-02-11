package health

import (
	"encoding/json"
	"net/http"
	"time"

	"image/internal/domain/ports"
)

// Handler implements the HealthHandler interface
type Handler struct {
	logger ports.Logger
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// NewHandler creates a new health check handler
func NewHandler(logger ports.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Handle processes health check requests
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode health check response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handle(w, r)
}

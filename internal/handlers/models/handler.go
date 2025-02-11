package modelshandler

import (
	"encoding/json"
	"net/http"

	"image/internal/domain/models"
	"image/internal/domain/ports"
)

// Handler handles model-related requests
type Handler struct {
	registry ports.ModelRegistry
	logger   ports.Logger
}

// NewHandler creates a new models handler
func NewHandler(registry ports.ModelRegistry, logger ports.Logger) *Handler {
	return &Handler{
		registry: registry,
		logger:   logger,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	modelsList := h.registry.List()

	// Convert models to response format
	response := models.ModelsResponse{
		Models: make([]models.ModelResponse, len(modelsList)),
	}

	for i, model := range modelsList {
		response.Models[i] = models.ToResponse(model)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode models response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

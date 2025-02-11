package text2img

import (
	"encoding/json"
	"net/http"

	"image/internal/domain/models"
	"image/internal/domain/ports"
	apperrors "image/pkg/errors"
)

// Handler implements the Text2ImgHandler interface
type Handler struct {
	service ports.ModelsLabService
	logger  ports.Logger
}

// NewHandler creates a new text-to-image handler instance
func NewHandler(service ports.ModelsLabService, logger ports.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Handle processes text-to-image generation requests
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.writeError(w, apperrors.NewInvalidRequestError(
			"Method not allowed",
			nil,
		))
		return
	}

	// Parse request body
	var req models.Text2ImgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, apperrors.NewInvalidRequestError(
			"Invalid request body",
			err,
		))
		return
	}

	// Generate image
	resp, err := h.service.GenerateImage(r.Context(), &req)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Write response
	h.writeJSON(w, resp)
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handle(w, r)
}

// writeJSON writes a JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		h.writeError(w, apperrors.NewInternalServerError(
			"Failed to encode response",
			err,
		))
	}
}

// writeError writes an error response
func (h *Handler) writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var response models.ErrorResponse
	var status int

	if appErr, ok := err.(*apperrors.AppError); ok {
		response = models.ErrorResponse{
			Status:  "error",
			Message: appErr.Message,
			Code:    string(appErr.Code),
		}
		status = appErr.Status
	} else {
		response = models.ErrorResponse{
			Status:  "error",
			Message: "Internal server error",
			Code:    string(apperrors.ErrInternalServer),
		}
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode error response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

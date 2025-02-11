package modelslab

import (
	"context"
	"fmt"

	"image/internal/domain/models"
	"image/internal/domain/ports"
	"image/internal/infrastructure/validation"
	apperrors "image/pkg/errors"
)

const (
	text2ImgEndpoint = "/images/text2img"
)

// Service implements the ModelsLabService interface
type Service struct {
	client    ports.HTTPClient
	validator *validation.Validator
	logger    ports.Logger
}

// NewService creates a new ModelsLab service instance
func NewService(client ports.HTTPClient, validator *validation.Validator, logger ports.Logger) *Service {
	return &Service{
		client:    client,
		validator: validator,
		logger:    logger,
	}
}

// GenerateImage generates an image from text using the ModelsLab API
func (s *Service) GenerateImage(ctx context.Context, req *models.Text2ImgRequest) (*models.Text2ImgResponse, error) {
	// Log the incoming request
	s.logger.Info("Processing text-to-image request",
		"model_id", req.ModelID,
		"width", req.Width,
		"height", req.Height,
		"samples", req.Samples,
	)

	// Validate the request
	if err := s.validateRequest(req); err != nil {
		s.logger.Error("Request validation failed", err)
		return nil, err
	}

	// Additional validations for specific fields
	if err := s.validator.ValidateScheduler(req.Scheduler); err != nil {
		return nil, err
	}

	if err := s.validator.ValidateEnhancePrompt(req.EnhancePrompt); err != nil {
		return nil, err
	}

	// Initialize response
	var response models.Text2ImgResponse

	// Call the ModelsLab API
	if err := s.client.Post(ctx, text2ImgEndpoint, req, &response); err != nil {
		s.logger.Error("Failed to generate image", err,
			"model_id", req.ModelID,
		)
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// Validate the response
	if err := s.validateResponse(&response); err != nil {
		s.logger.Error("Invalid response from API", err)
		return nil, err
	}

	s.logger.Info("Successfully generated image",
		"generation_time", response.GenerationTime,
		"image_count", len(response.Output),
	)

	return &response, nil
}

// validateRequest performs validation on the request
func (s *Service) validateRequest(req *models.Text2ImgRequest) error {
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	// Additional business logic validations
	if req.Width*req.Height > 1024*1024 {
		return apperrors.NewInvalidRequestError(
			"Image dimensions exceed maximum allowed size",
			nil,
		)
	}

	if req.Samples > 4 {
		return apperrors.NewInvalidRequestError(
			"Maximum number of samples exceeded (max: 4)",
			nil,
		)
	}

	return nil
}

// validateResponse validates the API response
func (s *Service) validateResponse(resp *models.Text2ImgResponse) error {
	if resp == nil {
		return apperrors.NewExternalAPIError("Empty response from API", nil)
	}

	if resp.Status != "success" {
		return apperrors.NewExternalAPIError(
			fmt.Sprintf("API returned non-success status: %s", resp.Status),
			nil,
		)
	}

	if len(resp.Output) == 0 {
		return apperrors.NewExternalAPIError("No images generated", nil)
	}

	return nil
}

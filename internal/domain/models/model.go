package models

import apperrors "image/pkg/errors"

// AIModel represents an AI model that can generate images
type AIModel interface {
	// ID returns the unique identifier of the model
	ID() string
	// Capabilities returns the model's capabilities
	Capabilities() ModelCapabilities
	// ValidateRequest validates a request against the model's capabilities
	ValidateRequest(req *Text2ImgRequest) error
}

// ModelCapabilities defines what a model can do
type ModelCapabilities struct {
	MaxWidth            int
	MaxHeight           int
	MaxSamples          int
	MinInferenceSteps   int
	MaxInferenceSteps   int
	SupportedSchedulers []string
	MinGuidanceScale    float64
	MaxGuidanceScale    float64
	SupportsUpscale     bool
	SupportsTomeSD      bool
	SupportsKarras      bool
}

// BaseModel provides common functionality for AI models
type BaseModel struct {
	id           string
	capabilities ModelCapabilities
}

// NewBaseModel creates a new BaseModel instance
func NewBaseModel(id string, capabilities ModelCapabilities) BaseModel {
	return BaseModel{
		id:           id,
		capabilities: capabilities,
	}
}

// ID returns the model's identifier
func (m BaseModel) ID() string {
	return m.id
}

// Capabilities returns the model's capabilities
func (m BaseModel) Capabilities() ModelCapabilities {
	return m.capabilities
}

// ValidateRequest performs basic validation against model capabilities
func (m BaseModel) ValidateRequest(req *Text2ImgRequest) error {
	caps := m.Capabilities()

	if req.Width > caps.MaxWidth {
		return apperrors.NewInvalidRequestError(
			"Width exceeds model maximum",
			nil,
		)
	}

	if req.Height > caps.MaxHeight {
		return apperrors.NewInvalidRequestError(
			"Height exceeds model maximum",
			nil,
		)
	}

	if req.Samples > caps.MaxSamples {
		return apperrors.NewInvalidRequestError(
			"Samples exceeds model maximum",
			nil,
		)
	}

	if req.NumInferenceSteps < caps.MinInferenceSteps || req.NumInferenceSteps > caps.MaxInferenceSteps {
		return apperrors.NewInvalidRequestError(
			"Number of inference steps outside model bounds",
			nil,
		)
	}

	if req.GuidanceScale < caps.MinGuidanceScale || req.GuidanceScale > caps.MaxGuidanceScale {
		return apperrors.NewInvalidRequestError(
			"Guidance scale outside model bounds",
			nil,
		)
	}

	if req.Scheduler != "" {
		validScheduler := false
		for _, s := range caps.SupportedSchedulers {
			if req.Scheduler == s {
				validScheduler = true
				break
			}
		}
		if !validScheduler {
			return apperrors.NewInvalidRequestError(
				"Unsupported scheduler for this model",
				nil,
			)
		}
	}

	if req.Upscale != "" && !caps.SupportsUpscale {
		return apperrors.NewInvalidRequestError(
			"Model does not support upscaling",
			nil,
		)
	}

	if req.Tomesd == "yes" && !caps.SupportsTomeSD {
		return apperrors.NewInvalidRequestError(
			"Model does not support TomeSD",
			nil,
		)
	}

	if req.UseKarrasSigmas == "yes" && !caps.SupportsKarras {
		return apperrors.NewInvalidRequestError(
			"Model does not support Karras sigmas",
			nil,
		)
	}

	return nil
}

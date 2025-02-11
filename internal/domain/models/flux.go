package models

// FluxModel represents the Flux AI model
type FluxModel struct {
	BaseModel
}

// NewFluxModel creates a new Flux model instance
func NewFluxModel() *FluxModel {
	return &FluxModel{
		BaseModel: NewBaseModel("flux", ModelCapabilities{
			MaxWidth:          768,
			MaxHeight:         768,
			MaxSamples:        4,
			MinInferenceSteps: 1,
			MaxInferenceSteps: 20,
			SupportedSchedulers: []string{
				"UniPCMultistepScheduler",
				"EulerAncestralDiscreteScheduler",
			},
			MinGuidanceScale: 1.0,
			MaxGuidanceScale: 20.0,
			SupportsUpscale:  false,
			SupportsTomeSD:   true,
			SupportsKarras:   true,
		}),
	}
}

// ValidateRequest extends base validation with Flux-specific rules
func (m *FluxModel) ValidateRequest(req *Text2ImgRequest) error {
	if err := m.BaseModel.ValidateRequest(req); err != nil {
		return err
	}

	// Add any Flux-specific validation rules here
	// For example, checking for specific prompt requirements or restrictions

	return nil
}

package models

// MidjourneyModel represents the Midjourney AI model
type MidjourneyModel struct {
	BaseModel
}

// NewMidjourneyModel creates a new Midjourney model instance
func NewMidjourneyModel() *MidjourneyModel {
	return &MidjourneyModel{
		BaseModel: NewBaseModel("midjourney", ModelCapabilities{
			MaxWidth:          1024,
			MaxHeight:         1024,
			MaxSamples:        4,
			MinInferenceSteps: 1,
			MaxInferenceSteps: 20,
			SupportedSchedulers: []string{
				"UniPCMultistepScheduler",
				"DDIMScheduler",
				"DPMSolverMultistepScheduler",
				"EulerAncestralDiscreteScheduler",
			},
			MinGuidanceScale: 1.0,
			MaxGuidanceScale: 20.0,
			SupportsUpscale:  true,
			SupportsTomeSD:   true,
			SupportsKarras:   true,
		}),
	}
}

// ValidateRequest extends base validation with Midjourney-specific rules
func (m *MidjourneyModel) ValidateRequest(req *Text2ImgRequest) error {
	if err := m.BaseModel.ValidateRequest(req); err != nil {
		return err
	}

	// Add any Midjourney-specific validation rules here
	// For example, checking for specific prompt requirements or restrictions

	return nil
}

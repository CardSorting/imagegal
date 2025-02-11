package models

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Text2ImgResponse represents the response from the text-to-image endpoint
type Text2ImgResponse struct {
	Status         string   `json:"status"`
	GenerationTime float64  `json:"generation_time,omitempty"`
	Output         []string `json:"output,omitempty"`
	Images         []string `json:"images,omitempty"`
	TaskID         string   `json:"task_id,omitempty"`
	Progress       float64  `json:"progress,omitempty"`
}

// IsProcessing returns true if the response indicates the request is still processing
func (r *Text2ImgResponse) IsProcessing() bool {
	return r.Status == "processing"
}

// IsSuccess returns true if the response indicates success
func (r *Text2ImgResponse) IsSuccess() bool {
	return r.Status == "success"
}

// ModelResponse represents a single model in the API response
type ModelResponse struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Capabilities CapabilitiesResponse `json:"capabilities"`
}

// CapabilitiesResponse represents model capabilities in the API response
type CapabilitiesResponse struct {
	MaxWidth            int      `json:"maxWidth"`
	MaxHeight           int      `json:"maxHeight"`
	MaxSamples          int      `json:"maxSamples"`
	MinInferenceSteps   int      `json:"minInferenceSteps"`
	MaxInferenceSteps   int      `json:"maxInferenceSteps"`
	SupportedSchedulers []string `json:"supportedSchedulers"`
	MinGuidanceScale    float64  `json:"minGuidanceScale"`
	MaxGuidanceScale    float64  `json:"maxGuidanceScale"`
	SupportsUpscale     bool     `json:"supportsUpscale"`
	SupportsTomeSD      bool     `json:"supportsTomeSD"`
	SupportsKarras      bool     `json:"supportsKarras"`
}

// ModelsResponse represents the response for the models endpoint
type ModelsResponse struct {
	Models []ModelResponse `json:"models"`
}

// ToResponse converts an AIModel to a ModelResponse
func ToResponse(model AIModel) ModelResponse {
	caps := model.Capabilities()
	return ModelResponse{
		ID:   model.ID(),
		Name: model.ID(), // Using ID as name since we don't have a separate name field
		Capabilities: CapabilitiesResponse{
			MaxWidth:            caps.MaxWidth,
			MaxHeight:           caps.MaxHeight,
			MaxSamples:          caps.MaxSamples,
			MinInferenceSteps:   caps.MinInferenceSteps,
			MaxInferenceSteps:   caps.MaxInferenceSteps,
			SupportedSchedulers: caps.SupportedSchedulers,
			MinGuidanceScale:    caps.MinGuidanceScale,
			MaxGuidanceScale:    caps.MaxGuidanceScale,
			SupportsUpscale:     caps.SupportsUpscale,
			SupportsTomeSD:      caps.SupportsTomeSD,
			SupportsKarras:      caps.SupportsKarras,
		},
	}
}

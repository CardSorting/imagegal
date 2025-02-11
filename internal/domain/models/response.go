package models

// Text2ImgResponse represents the response structure from the text-to-image generation API
type Text2ImgResponse struct {
	Status         string      `json:"status"`
	GenerationTime float64     `json:"generationTime"`
	ID            int         `json:"id"`
	Output        []string    `json:"output"`
	Meta          *MetaData   `json:"meta,omitempty"`
}

// MetaData represents the metadata included in the API response
type MetaData struct {
	Prompt          string      `json:"prompt"`
	ModelID         string      `json:"model_id"`
	NegativePrompt  string      `json:"negative_prompt,omitempty"`
	Scheduler       string      `json:"scheduler"`
	SafetyChecker   string      `json:"safetychecker"`
	Width           int         `json:"W"`
	Height          int         `json:"H"`
	GuidanceScale   float64     `json:"guidance_scale"`
	Seed            int64       `json:"seed"`
	Steps           int         `json:"steps"`
	Samples         int         `json:"n_samples"`
	FullURL         string      `json:"full_url"`
	Upscale         string      `json:"upscale"`
	Panorama        string      `json:"panorama"`
	SelfAttention   string      `json:"self_attention"`
	Embeddings      interface{} `json:"embeddings"`
	LoraModel       interface{} `json:"lora_model"`
	LoraStrength    interface{} `json:"lora_strength"`
	OutputDir       string      `json:"outdir"`
	FilePrefix      string      `json:"file_prefix"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

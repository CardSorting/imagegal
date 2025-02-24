package models

// Text2ImgRequest represents the request structure for text-to-image generation
type Text2ImgRequest struct {
	Key               string  `json:"key"`
	ModelID           string  `json:"model_id" validate:"required"`
	Prompt            string  `json:"prompt" validate:"required"`
	NegativePrompt    string  `json:"negative_prompt,omitempty"`
	Width             int     `json:"width" validate:"required,min=64,max=1024"`
	Height            int     `json:"height" validate:"required,min=64,max=1024"`
	Samples           int     `json:"samples" validate:"required,min=1,max=4"`
	NumInferenceSteps int     `json:"num_inference_steps" validate:"required,min=1,max=20"`
	SafetyChecker     string  `json:"safety_checker" validate:"omitempty,oneof=yes no"`
	EnhancePrompt     string  `json:"enhance_prompt" validate:"omitempty,oneof=yes no"`
	Seed              *int64  `json:"seed,omitempty"`
	GuidanceScale     float64 `json:"guidance_scale" validate:"omitempty,min=1,max=20"`
	Panorama          string  `json:"panorama" validate:"omitempty,oneof=yes no"`
	SelfAttention     string  `json:"self_attention" validate:"omitempty,oneof=yes no"`
	Upscale           string  `json:"upscale" validate:"omitempty,oneof=no 1 2 3"`
	EmbeddingsModel   string  `json:"embeddings_model,omitempty"`
	LoraModel         string  `json:"lora_model,omitempty"`
	Tomesd            string  `json:"tomesd" validate:"omitempty,oneof=yes no"`
	ClipSkip          string  `json:"clip_skip,omitempty"`
	UseKarrasSigmas   string  `json:"use_karras_sigmas" validate:"omitempty,oneof=yes no"`
	Vae               string  `json:"vae,omitempty"`
	LoraStrength      string  `json:"lora_strength,omitempty"`
	Scheduler         string  `json:"scheduler,omitempty"`
	Webhook           string  `json:"webhook,omitempty"`
	TrackID           string  `json:"track_id,omitempty"`
}

// ModelsLabAPIRequest represents the request structure expected by the ModelsLab API
type ModelsLabAPIRequest struct {
	Key               string  `json:"key"`
	ModelID           string  `json:"model_id"`
	Prompt            string  `json:"prompt"`
	NegativePrompt    string  `json:"negative_prompt,omitempty"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	Samples           int     `json:"samples"`
	NumInferenceSteps int     `json:"num_inference_steps"`
	SafetyChecker     bool    `json:"safety_checker,omitempty"`
	EnhancePrompt     bool    `json:"enhance_prompt,omitempty"`
	Seed              *int64  `json:"seed,omitempty"`
	GuidanceScale     float64 `json:"guidance_scale,omitempty"`
	Panorama          bool    `json:"panorama,omitempty"`
	SelfAttention     bool    `json:"self_attention,omitempty"`
	Upscale           string  `json:"upscale,omitempty"`
	EmbeddingsModel   string  `json:"embeddings_model,omitempty"`
	LoraModel         string  `json:"lora_model,omitempty"`
	Tomesd            bool    `json:"tomesd,omitempty"`
	ClipSkip          string  `json:"clip_skip,omitempty"`
	UseKarrasSigmas   bool    `json:"use_karras_sigmas,omitempty"`
	Vae               string  `json:"vae,omitempty"`
	LoraStrength      string  `json:"lora_strength,omitempty"`
	Scheduler         string  `json:"scheduler,omitempty"`
	Webhook           string  `json:"webhook,omitempty"`
	TrackID           string  `json:"track_id,omitempty"`
}

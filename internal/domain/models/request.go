package models

// Text2ImgRequest represents the request structure for text-to-image generation
type Text2ImgRequest struct {
	Key               string  `json:"key"`
	ModelID           string  `json:"model_id" validate:"required"`
	Prompt            string  `json:"prompt" validate:"required"`
	NegativePrompt    string  `json:"negative_prompt,omitempty"`
	Width             string  `json:"width" validate:"required"`
	Height            string  `json:"height" validate:"required"`
	Samples           string  `json:"samples" validate:"required"`
	NumInferenceSteps string  `json:"num_inference_steps" validate:"required"`
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

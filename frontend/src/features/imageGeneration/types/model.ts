export interface ModelCapabilities {
  maxWidth: number;
  maxHeight: number;
  maxSamples: number;
  minInferenceSteps: number;
  maxInferenceSteps: number;
  supportedSchedulers: string[];
  minGuidanceScale: number;
  maxGuidanceScale: number;
  supportsUpscale: boolean;
  supportsTomeSD: boolean;
  supportsKarras: boolean;
}

export interface AIModel {
  id: string;
  name: string;
  capabilities: ModelCapabilities;
}

export interface ModelsResponse {
  models: AIModel[];
}

// Default capabilities matching the most restrictive model (Flux)
export const defaultCapabilities: ModelCapabilities = {
  maxWidth: 768,
  maxHeight: 768,
  maxSamples: 4,
  minInferenceSteps: 1,
  maxInferenceSteps: 20,
  supportedSchedulers: [
    'UniPCMultistepScheduler',
    'EulerAncestralDiscreteScheduler',
  ],
  minGuidanceScale: 1,
  maxGuidanceScale: 20,
  supportsUpscale: false,
  supportsTomeSD: true,
  supportsKarras: true,
};

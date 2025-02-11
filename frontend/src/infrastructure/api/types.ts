export interface ApiError {
  status: string;
  message: string;
  code?: string;
}

export interface ApiResponse<T> {
  data: T;
  error?: ApiError;
}

export interface ImageGenerationRequest {
  model_id: string;
  prompt: string;
  negative_prompt?: string;
  width: number;
  height: number;
  samples: number;
  num_inference_steps: number;
  safety_checker?: 'yes' | 'no';
  enhance_prompt?: 'yes' | 'no';
  seed?: number | null;
  guidance_scale?: number;
  scheduler?: string;
  panorama?: 'yes' | 'no';
  self_attention?: 'yes' | 'no';
  upscale?: '1' | '2' | '3';
  tomesd?: 'yes' | 'no';
  use_karras_sigmas?: 'yes' | 'no';
}

export interface ImageGenerationResponse {
  status: string;
  generationTime: number;
  id: number;
  output: string[];
  meta?: {
    prompt: string;
    model_id: string;
    negative_prompt?: string;
    scheduler: string;
    safetychecker: string;
    W: number;
    H: number;
    guidance_scale: number;
    seed: number;
    steps: number;
    n_samples: number;
  };
}

export class ApiException extends Error {
  constructor(
    public status: number,
    public error: ApiError,
    public url: string
  ) {
    super(error.message);
    this.name = 'ApiException';
  }
}

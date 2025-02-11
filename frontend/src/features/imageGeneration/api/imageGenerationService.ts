import { apiClient } from '@/infrastructure/api/client';
import { ImageGenerationRequest, ImageGenerationResponse } from '@/infrastructure/api/types';
import { ModelsResponse, ModelCapabilities } from '../types/model';
import { env } from '@/infrastructure/config/env';

export interface ImageGenerationOptions extends Omit<ImageGenerationRequest, 'model_id'> {
  model_id?: string;
}

export class ImageGenerationService {
  private static instance: ImageGenerationService;
  private readonly baseUrl = '/api/v6';

  private constructor() {}

  public static getInstance(): ImageGenerationService {
    if (!ImageGenerationService.instance) {
      ImageGenerationService.instance = new ImageGenerationService();
    }
    return ImageGenerationService.instance;
  }

  public async getAvailableModels(): Promise<ModelsResponse> {
    const response = await apiClient.get<ModelsResponse>(
      `${this.baseUrl}/models`
    );
    return response.data;
  }

  public async generateImage(options: ImageGenerationOptions): Promise<ImageGenerationResponse> {
    const request: ImageGenerationRequest = {
      model_id: options.model_id || env.VITE_DEFAULT_MODEL,
      prompt: options.prompt,
      width: options.width,
      height: options.height,
      samples: options.samples,
      num_inference_steps: options.num_inference_steps,
      ...this.getOptionalParams(options),
    };

    const response = await apiClient.post<ImageGenerationResponse>(
      `${this.baseUrl}/images/text2img`,
      request
    );

    // If the response is in processing state, poll for completion
    if (response.data.status === 'processing' && response.data.task_id) {
      return this.pollForCompletion(response.data.task_id);
    }

    return response.data;
  }

  private async pollForCompletion(taskId: string): Promise<ImageGenerationResponse> {
    const maxAttempts = 30; // 30 attempts * 2 second delay = 60 seconds max
    let attempt = 0;

    while (attempt < maxAttempts) {
      await new Promise(resolve => setTimeout(resolve, 2000)); // Wait 2 seconds between attempts
      attempt++;

      const response = await apiClient.get<ImageGenerationResponse>(
        `${this.baseUrl}/images/status/${taskId}`
      );

      if (response.data.status === 'success') {
        return response.data;
      }

      if (response.data.status !== 'processing') {
        throw new Error(`Unexpected status: ${response.data.status}`);
      }
    }

    throw new Error('Image generation timed out');
  }

  private getOptionalParams(options: ImageGenerationOptions): Partial<ImageGenerationRequest> {
    const optionalParams: Partial<ImageGenerationRequest> = {};

    // Only include defined optional parameters
    if (options.negative_prompt) optionalParams.negative_prompt = options.negative_prompt;
    if (options.safety_checker) optionalParams.safety_checker = options.safety_checker;
    if (options.enhance_prompt) optionalParams.enhance_prompt = options.enhance_prompt;
    if (options.seed !== undefined) optionalParams.seed = options.seed;
    if (options.guidance_scale) optionalParams.guidance_scale = options.guidance_scale;
    if (options.scheduler) optionalParams.scheduler = options.scheduler;
    if (options.panorama) optionalParams.panorama = options.panorama;
    if (options.self_attention) optionalParams.self_attention = options.self_attention;
    if (options.upscale) optionalParams.upscale = options.upscale;
    if (options.tomesd) optionalParams.tomesd = options.tomesd;
    if (options.use_karras_sigmas) optionalParams.use_karras_sigmas = options.use_karras_sigmas;

    return optionalParams;
  }

  public async getDefaultSettings(capabilities: ModelCapabilities): Promise<Partial<ImageGenerationOptions>> {
    return {
      width: Math.min(512, capabilities.maxWidth),
      height: Math.min(512, capabilities.maxHeight),
      samples: 1,
      num_inference_steps: Math.min(30, capabilities.maxInferenceSteps),
      safety_checker: 'no',
      enhance_prompt: 'yes',
      guidance_scale: Math.min(7.5, capabilities.maxGuidanceScale),
      scheduler: capabilities.supportedSchedulers[0],
      tomesd: capabilities.supportsTomeSD ? 'yes' : undefined,
      use_karras_sigmas: capabilities.supportsKarras ? 'yes' : undefined,
    };
  }
}

// Export singleton instance
export const imageGenerationService = ImageGenerationService.getInstance();

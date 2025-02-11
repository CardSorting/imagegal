import { apiClient } from '@/infrastructure/api/client';
import { ImageGenerationRequest, ImageGenerationResponse } from '@/infrastructure/api/types';
import { env } from '@/infrastructure/config/env';

export interface ImageGenerationOptions extends Omit<ImageGenerationRequest, 'model_id'> {
  model_id?: string;
}

export class ImageGenerationService {
  private static instance: ImageGenerationService;
  private readonly baseUrl = '/api/v6/images';

  private constructor() {}

  public static getInstance(): ImageGenerationService {
    if (!ImageGenerationService.instance) {
      ImageGenerationService.instance = new ImageGenerationService();
    }
    return ImageGenerationService.instance;
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
      `${this.baseUrl}/text2img`,
      request
    );

    return response.data;
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

  public async getDefaultSettings(): Promise<Partial<ImageGenerationOptions>> {
    return {
      width: 512,
      height: 512,
      samples: 1,
      num_inference_steps: 30,
      safety_checker: 'no',
      enhance_prompt: 'yes',
      guidance_scale: 7.5,
      scheduler: 'UniPCMultistepScheduler',
      tomesd: 'yes',
      use_karras_sigmas: 'yes',
    };
  }
}

// Export singleton instance
export const imageGenerationService = ImageGenerationService.getInstance();

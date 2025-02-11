import { z } from 'zod';
import { ModelCapabilities } from '../types/model';

export const createImageGenerationSchema = (capabilities: ModelCapabilities) => z.object({
  prompt: z.string()
    .min(1, 'Prompt is required')
    .max(1000, 'Prompt must be less than 1000 characters'),
  
  negative_prompt: z.string()
    .max(1000, 'Negative prompt must be less than 1000 characters')
    .optional(),
  
  width: z.number()
    .int()
    .min(64, 'Width must be at least 64 pixels')
    .max(capabilities.maxWidth, `Width must be at most ${capabilities.maxWidth} pixels`),
  
  height: z.number()
    .int()
    .min(64, 'Height must be at least 64 pixels')
    .max(capabilities.maxHeight, `Height must be at most ${capabilities.maxHeight} pixels`),
  
  samples: z.number()
    .int()
    .min(1, 'Must generate at least 1 image')
    .max(capabilities.maxSamples, `Can generate at most ${capabilities.maxSamples} images`),
  
  num_inference_steps: z.number()
    .int()
    .min(capabilities.minInferenceSteps, `Steps must be at least ${capabilities.minInferenceSteps}`)
    .max(capabilities.maxInferenceSteps, `Steps must be at most ${capabilities.maxInferenceSteps}`),
  
  safety_checker: z.enum(['yes', 'no']).optional(),
  enhance_prompt: z.enum(['yes', 'no']).optional(),
  
  seed: z.number().int().nullable().optional(),
  
  guidance_scale: z.number()
    .min(capabilities.minGuidanceScale, `Guidance scale must be at least ${capabilities.minGuidanceScale}`)
    .max(capabilities.maxGuidanceScale, `Guidance scale must be at most ${capabilities.maxGuidanceScale}`)
    .optional(),
  
  scheduler: z.enum(capabilities.supportedSchedulers as [string, ...string[]]).optional(),
  
  panorama: z.enum(['yes', 'no']).optional(),
  self_attention: z.enum(['yes', 'no']).optional(),
  upscale: capabilities.supportsUpscale 
    ? z.enum(['1', '2', '3']).optional()
    : z.literal(undefined),
  tomesd: capabilities.supportsTomeSD 
    ? z.enum(['yes', 'no']).optional()
    : z.literal(undefined),
  
  use_karras_sigmas: capabilities.supportsKarras 
    ? z.enum(['yes', 'no']).optional()
    : z.literal(undefined),
});

export type ImageGenerationFormData = z.infer<ReturnType<typeof createImageGenerationSchema>>;

export const getDefaultFormValues = (capabilities: ModelCapabilities): Partial<ImageGenerationFormData> => ({
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
});

export const schedulerOptions = [
  'DDPMScheduler',
  'DDIMScheduler',
  'PNDMScheduler',
  'LMSDiscreteScheduler',
  'EulerDiscreteScheduler',
  'EulerAncestralDiscreteScheduler',
  'DPMSolverMultistepScheduler',
  'HeunDiscreteScheduler',
  'KDPM2DiscreteScheduler',
  'DPMSolverSinglestepScheduler',
  'KDPM2AncestralDiscreteScheduler',
  'UniPCMultistepScheduler',
  'DDIMInverseScheduler',
  'DEISMultistepScheduler',
  'IPNDMScheduler',
  'KarrasVeScheduler',
  'ScoreSdeVeScheduler',
  'LCMScheduler',
] as const;

export type Scheduler = typeof schedulerOptions[number];

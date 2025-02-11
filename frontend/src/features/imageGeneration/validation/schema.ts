import { z } from 'zod';

export const imageGenerationSchema = z.object({
  prompt: z.string()
    .min(1, 'Prompt is required')
    .max(1000, 'Prompt must be less than 1000 characters'),
  
  negative_prompt: z.string()
    .max(1000, 'Negative prompt must be less than 1000 characters')
    .optional(),
  
  width: z.number()
    .int()
    .min(64, 'Width must be at least 64 pixels')
    .max(1024, 'Width must be at most 1024 pixels'),
  
  height: z.number()
    .int()
    .min(64, 'Height must be at least 64 pixels')
    .max(1024, 'Height must be at most 1024 pixels'),
  
  samples: z.number()
    .int()
    .min(1, 'Must generate at least 1 image')
    .max(4, 'Can generate at most 4 images'),
  
  num_inference_steps: z.number()
    .int()
    .min(1, 'Steps must be at least 1')
    .max(20, 'Steps must be at most 20'),
  
  safety_checker: z.enum(['yes', 'no']).optional(),
  enhance_prompt: z.enum(['yes', 'no']).optional(),
  
  seed: z.number().int().nullable().optional(),
  
  guidance_scale: z.number()
    .min(1, 'Guidance scale must be at least 1')
    .max(20, 'Guidance scale must be at most 20')
    .optional(),
  
  scheduler: z.string().optional(),
  
  panorama: z.enum(['yes', 'no']).optional(),
  self_attention: z.enum(['yes', 'no']).optional(),
  upscale: z.enum(['1', '2', '3']).optional(),
  tomesd: z.enum(['yes', 'no']).optional(),
  use_karras_sigmas: z.enum(['yes', 'no']).optional(),
});

export type ImageGenerationFormData = z.infer<typeof imageGenerationSchema>;

export const defaultFormValues: Partial<ImageGenerationFormData> = {
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

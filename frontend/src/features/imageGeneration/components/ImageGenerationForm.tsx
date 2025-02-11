import React from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'react-hot-toast';
import { ImageGenerationResponse } from '@/infrastructure/api/types';

import { Button } from '@/shared/components/Button';
import { Input } from '@/shared/components/Input';
import { Slider } from '@/shared/components/Slider';
import { Textarea } from '@/shared/components/Textarea';
import { Select } from '@/shared/components/Select';
import { Switch } from '@/shared/components/Switch';
import { ProgressSpinner } from '@/shared/components/Spinner';

import { useImageGeneration } from '../hooks/useImageGeneration';
import { 
  imageGenerationSchema, 
  type ImageGenerationFormData,
  defaultFormValues,
  schedulerOptions 
} from '../validation/schema';

interface ImageGenerationFormProps {
  onSuccess?: (result: ImageGenerationResponse) => void;
}

export const ImageGenerationForm: React.FC<ImageGenerationFormProps> = ({ onSuccess }) => {
  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<ImageGenerationFormData>({
    resolver: zodResolver(imageGenerationSchema),
    defaultValues: defaultFormValues,
  });

  const { generate, isGenerating, progress } = useImageGeneration();

  const onSubmit = async (data: ImageGenerationFormData) => {
    try {
      const result = await generate(data);

      if (onSuccess) {
        onSuccess(result);
      }
      toast.success('Image generated successfully!');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to generate image');
    }
  };

  const schedulerSelectOptions = schedulerOptions.map(scheduler => ({
    value: scheduler,
    label: scheduler,
  }));

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 max-w-2xl mx-auto">
      <Controller
        name="prompt"
        control={control}
        render={({ field }) => (
          <Textarea
            label="Prompt"
            error={errors.prompt?.message}
            placeholder="Enter your image description..."
            rows={4}
            {...field}
          />
        )}
      />

      <Controller
        name="negative_prompt"
        control={control}
        render={({ field }) => (
          <Textarea
            label="Negative Prompt"
            error={errors.negative_prompt?.message}
            placeholder="Enter things to exclude from the image..."
            rows={2}
            {...field}
          />
        )}
      />

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Controller
          name="width"
          control={control}
          render={({ field: { onChange, ...field } }) => (
            <Slider
              label="Width"
              error={errors.width?.message}
              min={64}
              max={1024}
              step={64}
              valueSuffix="px"
              helperText="Image width in pixels (must be divisible by 64)"
              onChange={(e) => onChange(parseInt(e.target.value, 10))}
              {...field}
            />
          )}
        />

        <Controller
          name="height"
          control={control}
          render={({ field: { onChange, ...field } }) => (
            <Slider
              label="Height"
              error={errors.height?.message}
              min={64}
              max={1024}
              step={64}
              valueSuffix="px"
              helperText="Image height in pixels (must be divisible by 64)"
              onChange={(e) => onChange(parseInt(e.target.value, 10))}
              {...field}
            />
          )}
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Controller
          name="samples"
          control={control}
          render={({ field: { onChange, ...field } }) => (
            <Slider
              label="Number of Images"
              error={errors.samples?.message}
              min={1}
              max={4}
              step={1}
              helperText="Number of images to generate in one request"
              onChange={(e) => onChange(parseInt(e.target.value, 10))}
              {...field}
            />
          )}
        />

        <Controller
          name="num_inference_steps"
          control={control}
          render={({ field: { onChange, ...field } }) => (
            <Slider
              label="Inference Steps"
              error={errors.num_inference_steps?.message}
              min={1}
              max={20}
              step={1}
              helperText="More steps generally means higher quality but slower generation"
              onChange={(e) => onChange(parseInt(e.target.value, 10))}
              {...field}
            />
          )}
        />
      </div>

      <Controller
        name="scheduler"
        control={control}
        render={({ field }) => (
          <Select
            label="Scheduler"
            options={schedulerSelectOptions}
            error={errors.scheduler?.message}
            {...field}
          />
        )}
      />

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Controller
          name="guidance_scale"
          control={control}
          render={({ field: { onChange, ...field } }) => (
            <Slider
              label="Guidance Scale"
              error={errors.guidance_scale?.message}
              min={1}
              max={20}
              step={0.1}
              helperText="How closely to follow the prompt (higher = more strict)"
              onChange={(e) => onChange(parseFloat(e.target.value))}
              {...field}
            />
          )}
        />

        <Controller
          name="seed"
          control={control}
          render={({ field: { value, onChange, ...field } }) => (
            <Input
              type="number"
              label="Seed"
              placeholder="Random"
              error={errors.seed?.message}
              value={value === null ? '' : value}
              onChange={(e) => {
                const val = e.target.value;
                onChange(val === '' ? null : parseInt(val, 10));
              }}
              {...field}
            />
          )}
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Controller
          name="safety_checker"
          control={control}
          render={({ field: { value, onChange, ...field } }) => (
            <Switch
              label="Safety Checker"
              description="Check for NSFW content"
              checked={value === 'yes'}
              onChange={(e) => onChange(e.target.checked ? 'yes' : 'no')}
              {...field}
            />
          )}
        />

        <Controller
          name="enhance_prompt"
          control={control}
          render={({ field: { value, onChange, ...field } }) => (
            <Switch
              label="Enhance Prompt"
              description="Improve prompt automatically"
              checked={value === 'yes'}
              onChange={(e) => onChange(e.target.checked ? 'yes' : 'no')}
              {...field}
            />
          )}
        />
      </div>

      <div className="flex justify-between items-center pt-4">
        <Button
          type="button"
          variant="secondary"
          onClick={() => reset(defaultFormValues)}
          disabled={isGenerating}
        >
          Reset
        </Button>

        <div className="flex items-center gap-4">
          {isGenerating && <ProgressSpinner progress={progress} />}
          <Button
            type="submit"
            disabled={isGenerating}
            isLoading={isGenerating}
          >
            Generate Image
          </Button>
        </div>
      </div>
    </form>
  );
};

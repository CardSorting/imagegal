import { useMutation, useQueryClient } from '@tanstack/react-query';
import { imageGenerationService, ImageGenerationOptions } from '../api/imageGenerationService';
import { ImageGenerationResponse } from '@/infrastructure/api/types';
import { ApiException } from '@/infrastructure/api/types';
import { useState } from 'react';

export interface UseImageGenerationResult {
  generate: (options: ImageGenerationOptions) => Promise<ImageGenerationResponse>;
  isGenerating: boolean;
  error: ApiException | null;
  lastGeneratedImage: ImageGenerationResponse | null;
  progress: number;
  reset: () => void;
}

export function useImageGeneration(): UseImageGenerationResult {
  const queryClient = useQueryClient();
  const [lastGeneratedImage, setLastGeneratedImage] = useState<ImageGenerationResponse | null>(null);
  const [progress, setProgress] = useState<number>(0);

  const mutation = useMutation<ImageGenerationResponse, ApiException, ImageGenerationOptions>({
    mutationFn: async (options) => {
      try {
        // Reset state
        setProgress(0);
        
        // Start progress simulation
        const progressInterval = setInterval(() => {
          setProgress((prev) => {
            const next = prev + (100 - prev) * 0.1;
            return next > 95 ? 95 : next;
          });
        }, 500);

        // Generate image
        const response = await imageGenerationService.generateImage(options);
        
        // Cleanup and finalize progress
        clearInterval(progressInterval);
        setProgress(100);
        setLastGeneratedImage(response);

        // Invalidate relevant queries if needed
        await queryClient.invalidateQueries({ queryKey: ['images'] });

        return response;
      } catch (error) {
        setProgress(0);
        throw error;
      }
    },
  });

  const reset = () => {
    setLastGeneratedImage(null);
    setProgress(0);
  };

  return {
    generate: mutation.mutateAsync,
    isGenerating: mutation.isPending,
    error: mutation.error || null,
    lastGeneratedImage,
    progress,
    reset,
  };
}

// Custom hook for default settings
export function useDefaultSettings() {
  return {
    getDefaults: imageGenerationService.getDefaultSettings,
  };
}

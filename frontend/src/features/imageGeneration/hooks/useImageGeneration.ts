import { useMutation, useQueryClient } from '@tanstack/react-query';
import { imageGenerationService, ImageGenerationOptions } from '../api/imageGenerationService';
import { ImageGenerationResponse } from '@/infrastructure/api/types';
import { ApiException } from '@/infrastructure/api/types';
import { useState } from 'react';
import { useModel } from '../context/ModelContext';

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
  const { selectedModel } = useModel();
  const [lastGeneratedImage, setLastGeneratedImage] = useState<ImageGenerationResponse | null>(null);
  const [progress, setProgress] = useState<number>(0);

  const mutation = useMutation<ImageGenerationResponse, ApiException, ImageGenerationOptions>({
    mutationFn: async (options) => {
      if (!selectedModel) {
        throw new Error('No model selected');
      }

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

        // Generate image with selected model
        const response = await imageGenerationService.generateImage({
          ...options,
          model_id: selectedModel.id,
        });
        
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
  const { capabilities } = useModel();
  
  return {
    getDefaults: () => imageGenerationService.getDefaultSettings(capabilities),
  };
}

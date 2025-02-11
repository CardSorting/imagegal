import { useMutation, useQueryClient } from '@tanstack/react-query';
import { imageGenerationService, ImageGenerationOptions } from '../api/imageGenerationService';
import { ImageGenerationResponse } from '@/infrastructure/api/types';
import { ApiException } from '@/infrastructure/api/types';
import { useState } from 'react';
import { useModel } from '../context/ModelContext';
import { toast } from 'react-hot-toast';

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
        setLastGeneratedImage(null);

        // Generate image with selected model
        const response = await imageGenerationService.generateImage({
          ...options,
          model_id: selectedModel.id,
        });

        // Update progress based on API response
        if (response.status === 'processing') {
          setProgress(response.progress || 0);
        } else if (response.status === 'success') {
          setProgress(100);
          setLastGeneratedImage(response);
        }

        // Invalidate relevant queries if needed
        await queryClient.invalidateQueries({ queryKey: ['images'] });

        return response;
      } catch (error) {
        setProgress(0);
        throw error;
      }
    },
    onError: (error) => {
      setProgress(0);
      toast.error(error.message);
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

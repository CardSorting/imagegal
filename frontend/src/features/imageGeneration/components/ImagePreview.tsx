import React from 'react';
import { ImageGenerationResponse } from '@/infrastructure/api/types';
import { Button } from '@/shared/components/Button';

interface ImagePreviewProps {
  result: ImageGenerationResponse | null;
  onDownload?: (url: string) => void;
}

export const ImagePreview: React.FC<ImagePreviewProps> = ({ result, onDownload }) => {
  if (!result || !result.output.length) {
    return null;
  }

  const handleDownload = async (url: string) => {
    if (onDownload) {
      onDownload(url);
    } else {
      try {
        const response = await fetch(url);
        const blob = await response.blob();
        const downloadUrl = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = `generated-image-${result.id}-${Date.now()}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(downloadUrl);
      } catch (error) {
        console.error('Failed to download image:', error);
      }
    }
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {result.output.map((imageUrl, index) => (
          <div
            key={`${result.id}-${index}`}
            className="relative group rounded-lg overflow-hidden bg-muted"
          >
            <img
              src={imageUrl}
              alt={`Generated image ${index + 1}`}
              className="w-full h-auto object-cover"
              loading="lazy"
            />
            <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
              <Button
                variant="secondary"
                size="sm"
                onClick={() => handleDownload(imageUrl)}
              >
                Download
              </Button>
            </div>
          </div>
        ))}
      </div>

      {result.meta && (
        <div className="bg-muted rounded-lg p-4 text-sm space-y-2">
          <h3 className="font-medium">Generation Details</h3>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
            <div>
              <span className="text-muted-foreground">Model:</span>{' '}
              {result.meta.model_id}
            </div>
            <div>
              <span className="text-muted-foreground">Seed:</span>{' '}
              {result.meta.seed}
            </div>
            <div>
              <span className="text-muted-foreground">Steps:</span>{' '}
              {result.meta.steps}
            </div>
            <div>
              <span className="text-muted-foreground">CFG Scale:</span>{' '}
              {result.meta.guidance_scale}
            </div>
            <div>
              <span className="text-muted-foreground">Size:</span>{' '}
              {result.meta.W}x{result.meta.H}
            </div>
            <div>
              <span className="text-muted-foreground">Time:</span>{' '}
              {result.generationTime.toFixed(2)}s
            </div>
          </div>
          <div className="pt-2">
            <span className="text-muted-foreground">Prompt:</span>{' '}
            {result.meta.prompt}
          </div>
          {result.meta.negative_prompt && (
            <div>
              <span className="text-muted-foreground">Negative Prompt:</span>{' '}
              {result.meta.negative_prompt}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

import React from 'react';
import { Select } from '@/shared/components/Select';
import { useModel } from '../context/ModelContext';

export const ModelSelector: React.FC = () => {
  const { models, selectedModel, selectModel, isLoading } = useModel();

  const modelOptions = models.map(model => ({
    value: model.id,
    label: model.name,
  }));

  return (
    <div className="mb-6">
      <Select
        label="AI Model"
        value={selectedModel?.id || ''}
        options={modelOptions}
        onChange={(e) => selectModel(e.target.value)}
        disabled={isLoading}
        placeholder={isLoading ? 'Loading models...' : 'Select a model'}
      />
      {selectedModel && (
        <div className="mt-2 text-sm text-muted-foreground">
          <p>Max Resolution: {selectedModel.capabilities.maxWidth}x{selectedModel.capabilities.maxHeight}</p>
          <p>Max Samples: {selectedModel.capabilities.maxSamples}</p>
          {selectedModel.capabilities.supportsUpscale && <p>Supports Upscaling</p>}
          {selectedModel.capabilities.supportsTomeSD && <p>Supports TomeSD</p>}
          {selectedModel.capabilities.supportsKarras && <p>Supports Karras Scheduling</p>}
        </div>
      )}
    </div>
  );
};

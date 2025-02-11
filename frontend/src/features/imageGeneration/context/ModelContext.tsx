import React, { createContext, useContext, useState, useEffect } from 'react';
import { AIModel, ModelCapabilities, defaultCapabilities } from '../types/model';
import { imageGenerationService } from '../api/imageGenerationService';
import { toast } from 'react-hot-toast';

interface ModelContextType {
  models: AIModel[];
  selectedModel: AIModel | null;
  capabilities: ModelCapabilities;
  isLoading: boolean;
  selectModel: (modelId: string) => void;
}

const ModelContext = createContext<ModelContextType | undefined>(undefined);

export const ModelProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [models, setModels] = useState<AIModel[]>([]);
  const [selectedModel, setSelectedModel] = useState<AIModel | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchModels = async () => {
      try {
        const response = await imageGenerationService.getAvailableModels();
        setModels(response.models);
        
        // Select the first model by default
        if (response.models.length > 0) {
          setSelectedModel(response.models[0]);
        }
      } catch (error) {
        toast.error('Failed to fetch available models');
        console.error('Failed to fetch models:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchModels();
  }, []);

  const selectModel = (modelId: string) => {
    const model = models.find(m => m.id === modelId);
    if (model) {
      setSelectedModel(model);
    }
  };

  const value = {
    models,
    selectedModel,
    capabilities: selectedModel?.capabilities || defaultCapabilities,
    isLoading,
    selectModel,
  };

  return (
    <ModelContext.Provider value={value}>
      {children}
    </ModelContext.Provider>
  );
};

export const useModel = () => {
  const context = useContext(ModelContext);
  if (context === undefined) {
    throw new Error('useModel must be used within a ModelProvider');
  }
  return context;
};

import React, { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import { ImageGenerationForm } from '@/features/imageGeneration/components/ImageGenerationForm';
import { ImagePreview } from '@/features/imageGeneration/components/ImagePreview';
import { ImageGenerationResponse } from '@/infrastructure/api/types';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: false,
    },
  },
});

function App() {
  const [generatedImage, setGeneratedImage] = useState<ImageGenerationResponse | null>(null);

  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-background text-foreground">
        <div className="container mx-auto px-4 py-8">
          <header className="text-center mb-8">
            <h1 className="text-4xl font-bold mb-2">AI Image Generator</h1>
            <p className="text-muted-foreground">
              Create stunning images from text descriptions using advanced AI
            </p>
          </header>

          <main className="space-y-8">
            <ImageGenerationForm onSuccess={setGeneratedImage} />
            
            {generatedImage && (
              <section className="mt-8">
                <h2 className="text-2xl font-semibold mb-4">Generated Images</h2>
                <ImagePreview result={generatedImage} />
              </section>
            )}
          </main>

          <footer className="mt-16 text-center text-sm text-muted-foreground">
            <p>Powered by ModelsLab API</p>
          </footer>
        </div>
      </div>
      <Toaster position="top-right" />
    </QueryClientProvider>
  );
}

export default App;

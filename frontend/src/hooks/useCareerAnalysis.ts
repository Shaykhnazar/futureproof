import { useState } from 'react';
import { careersAPI } from '../api';
import { useCareerStore } from '../store';
import type { AnalysisRequest, AnalysisResult } from '../types';

export const useCareerAnalysis = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { setAnalysis, setAnalyzing } = useCareerStore();

  const analyzeCareer = async (request: AnalysisRequest): Promise<AnalysisResult | null> => {
    setIsLoading(true);
    setAnalyzing(true);
    setError(null);

    try {
      const result = await careersAPI.analyzeCareer(request);
      setAnalysis(result);
      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to analyze career';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
      setAnalyzing(false);
    }
  };

  return {
    analyzeCareer,
    isLoading,
    error,
  };
};

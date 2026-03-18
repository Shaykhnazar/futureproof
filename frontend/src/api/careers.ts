import { apiClient } from './client';
import type { Profession, CareerTransition, AnalysisRequest, AnalysisResult } from '../types';

export const careersAPI = {
  getAllProfessions: async (): Promise<Profession[]> => {
    const response = await apiClient.get<{ professions: Profession[] }>('/professions');
    return response.data.professions;
  },

  getProfessionBySlug: async (slug: string): Promise<Profession> => {
    const response = await apiClient.get<Profession>(`/professions/${slug}`);
    return response.data;
  },

  getFutureProfessions: async (): Promise<Profession[]> => {
    const response = await apiClient.get<{ professions: Profession[] }>('/professions/future');
    return response.data.professions;
  },

  getCareerTransitions: async (slug: string): Promise<CareerTransition[]> => {
    const response = await apiClient.get<{ transitions: CareerTransition[] }>(`/professions/${slug}/pivots`);
    return response.data.transitions;
  },

  analyzeCareer: async (request: AnalysisRequest): Promise<AnalysisResult> => {
    const response = await apiClient.post<AnalysisResult>('/analyze', request);
    return response.data;
  },

  saveCareer: async (professionId: string, notes?: string): Promise<void> => {
    await apiClient.post('/careers/save', {
      profession_id: professionId,
      notes: notes || '',
    });
  },

  getSavedCareers: async (): Promise<Profession[]> => {
    const response = await apiClient.get<{ professions: Profession[] }>('/careers/saved');
    return response.data.professions;
  },
};

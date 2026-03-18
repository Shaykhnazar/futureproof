import { create } from 'zustand';
import type { Profession, AnalysisResult } from '../types';

interface CareerStore {
  professions: Profession[];
  futureProfessions: Profession[];
  selectedProfession: Profession | null;
  analysis: AnalysisResult | null;
  savedCareers: Profession[];
  isAnalyzing: boolean;
  error: string | null;

  setProfessions: (professions: Profession[]) => void;
  setFutureProfessions: (professions: Profession[]) => void;
  setSelectedProfession: (profession: Profession | null) => void;
  setAnalysis: (analysis: AnalysisResult | null) => void;
  setSavedCareers: (careers: Profession[]) => void;
  setAnalyzing: (analyzing: boolean) => void;
  setError: (error: string | null) => void;
  addSavedCareer: (career: Profession) => void;
}

export const useCareerStore = create<CareerStore>((set) => ({
  professions: [],
  futureProfessions: [],
  selectedProfession: null,
  analysis: null,
  savedCareers: [],
  isAnalyzing: false,
  error: null,

  setProfessions: (professions) => set({ professions }),

  setFutureProfessions: (professions) => set({ futureProfessions: professions }),

  setSelectedProfession: (profession) => set({ selectedProfession: profession }),

  setAnalysis: (analysis) => set({ analysis }),

  setSavedCareers: (careers) => set({ savedCareers: careers }),

  setAnalyzing: (analyzing) => set({ isAnalyzing: analyzing }),

  setError: (error) => set({ error }),

  addSavedCareer: (career) =>
    set((state) => ({
      savedCareers: [...state.savedCareers, career],
    })),
}));

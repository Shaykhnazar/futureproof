import { create } from 'zustand';
import type { CityWithScore } from '../types';

interface GlobeStore {
  cities: CityWithScore[];
  selectedCity: CityWithScore | null;
  isLoading: boolean;
  error: string | null;
  setCities: (cities: CityWithScore[]) => void;
  setSelectedCity: (city: CityWithScore | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  updateCityScore: (cityId: string, newScore: number) => void;
}

export const useGlobeStore = create<GlobeStore>((set) => ({
  cities: [],
  selectedCity: null,
  isLoading: false,
  error: null,

  setCities: (cities) => set({ cities }),

  setSelectedCity: (city) => set({ selectedCity: city }),

  setLoading: (loading) => set({ isLoading: loading }),

  setError: (error) => set({ error }),

  updateCityScore: (cityId, newScore) =>
    set((state) => ({
      cities: state.cities.map((city) =>
        city.id === cityId ? { ...city, score: newScore } : city
      ),
      selectedCity:
        state.selectedCity?.id === cityId
          ? { ...state.selectedCity, score: newScore }
          : state.selectedCity,
    })),
}));

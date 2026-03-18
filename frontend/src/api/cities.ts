import { apiClient } from './client';
import type { CityWithScore } from '../types';

export const citiesAPI = {
  getAllCities: async (): Promise<CityWithScore[]> => {
    const response = await apiClient.get<{ cities: CityWithScore[] }>('/cities');
    return response.data.cities;
  },

  getCityById: async (id: string): Promise<CityWithScore> => {
    const response = await apiClient.get<CityWithScore>(`/cities/${id}`);
    return response.data;
  },

  getCitiesByRegion: async (region: string): Promise<CityWithScore[]> => {
    const response = await apiClient.get<{ cities: CityWithScore[] }>(`/cities/region/${region}`);
    return response.data.cities;
  },
};

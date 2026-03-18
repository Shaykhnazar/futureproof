import { apiClient } from './client';
import type { LoginRequest, RegisterRequest, LoginResponse, UserWithProfile, UserProfile } from '../types';

export const authAPI = {
  register: async (request: RegisterRequest): Promise<void> => {
    await apiClient.post('/auth/register', request);
  },

  login: async (request: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/login', request);
    const { access_token, refresh_token } = response.data;

    // Store tokens
    localStorage.setItem('access_token', access_token);
    localStorage.setItem('refresh_token', refresh_token);

    return response.data;
  },

  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  },

  getCurrentUser: async (): Promise<UserWithProfile> => {
    const response = await apiClient.get<UserWithProfile>('/users/me');
    return response.data;
  },

  updateProfile: async (profile: Partial<UserProfile>): Promise<UserProfile> => {
    const response = await apiClient.put<{ profile: UserProfile }>('/users/profile', profile);
    return response.data.profile;
  },

  isAuthenticated: (): boolean => {
    return !!localStorage.getItem('access_token');
  },
};

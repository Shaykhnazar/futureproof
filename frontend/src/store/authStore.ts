import { create } from 'zustand';
import type { UserWithProfile } from '../types';
import { authAPI } from '../api';

interface AuthStore {
  user: UserWithProfile | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  setUser: (user: UserWithProfile | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, name: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  isAuthenticated: authAPI.isAuthenticated(),
  isLoading: false,
  error: null,

  setUser: (user) => set({ user, isAuthenticated: !!user }),

  setLoading: (loading) => set({ isLoading: loading }),

  setError: (error) => set({ error }),

  login: async (email, password) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authAPI.login({ email, password });
      set({ user: response.user, isAuthenticated: true, isLoading: false });
    } catch (error) {
      set({
        error: 'Invalid email or password',
        isLoading: false,
        isAuthenticated: false,
      });
      throw error;
    }
  },

  register: async (email, name, password) => {
    set({ isLoading: true, error: null });
    try {
      await authAPI.register({ email, name, password });
      // After registration, log them in
      await authAPI.login({ email, password });
      const user = await authAPI.getCurrentUser();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch (error) {
      set({ error: 'Registration failed', isLoading: false });
      throw error;
    }
  },

  logout: () => {
    authAPI.logout();
    set({ user: null, isAuthenticated: false });
  },

  checkAuth: async () => {
    if (!authAPI.isAuthenticated()) {
      set({ isAuthenticated: false, user: null });
      return;
    }

    set({ isLoading: true });
    try {
      const user = await authAPI.getCurrentUser();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch (error) {
      authAPI.logout();
      set({ user: null, isAuthenticated: false, isLoading: false });
    }
  },
}));

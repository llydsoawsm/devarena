import { create } from 'zustand';
import { User } from '../lib/types';
import { authApi } from '../lib/api';

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (email: string, pass: string) => Promise<void>;
  register: (name: string, email: string, pass: string) => Promise<void>;
  logout: () => void;
  loadUser: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: typeof window !== 'undefined' ? localStorage.getItem('devcode_token') : null,
  isLoading: true,

  login: async (email, password) => {
    const res = await authApi.login({ email, password });
    localStorage.setItem('devcode_token', res.data.access_token);
    set({ user: res.data.user, token: res.data.access_token });
  },

  register: async (name, email, password) => {
    const res = await authApi.register({ name, email, password });
    localStorage.setItem('devcode_token', res.data.access_token);
    set({ user: res.data.user, token: res.data.access_token });
  },

  logout: () => {
    localStorage.removeItem('devcode_token');
    set({ user: null, token: null });
  },

  loadUser: async () => {
    const token = localStorage.getItem('devcode_token');
    if (!token) {
      set({ isLoading: false });
      return;
    }
    try {
      const res = await authApi.getMe();
      set({ user: res.data.user, token, isLoading: false });
    } catch {
      localStorage.removeItem('devcode_token');
      set({ user: null, token: null, isLoading: false });
    }
  },
}));

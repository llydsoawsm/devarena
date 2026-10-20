import axios from 'axios';
import { User, Problem, Match, Submission, ExecutionResult } from './types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('devcode_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Auth API
export const authApi = {
  register: (data: { name: string; email: string; password: string }) =>
    apiClient.post<{ user: User; access_token: string; refresh_token: string }>('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    apiClient.post<{ user: User; access_token: string; refresh_token: string }>('/auth/login', data),
  getMe: () => apiClient.get<{ user: User }>('/auth/me'),
};

// Problems API
export const problemsApi = {
  list: () => apiClient.get<{ problems: Problem[] }>('/problems'),
  getBySlug: (slug: string) => apiClient.get<{ problem: Problem }>(`/problems/${slug}`),
  getById: (id: number) => apiClient.get<{ problem: Problem }>(`/problems/id/${id}`),
};

// Matchmaking API
export const matchmakingApi = {
  enqueue: (mode: 'ranked' | 'casual' | 'tournament') =>
    apiClient.post<{ message: string; mode: string; rating: number }>('/matchmaking/enqueue', { mode }),
  dequeue: () => apiClient.post<{ message: string }>('/matchmaking/dequeue'),
  status: () => apiClient.get<{ is_queued: boolean }>('/matchmaking/status'),
};

// Matches API
export const matchesApi = {
  getByPublicId: (publicId: string) => apiClient.get<{ match: Match }>(`/matches/${publicId}`),
  forfeit: (publicId: string) => apiClient.post<{ message: string }>(`/matches/${publicId}/forfeit`),
  history: () => apiClient.get<{ matches: Match[] }>('/matches/user/history'),
  leaderboard: () => apiClient.get<{ leaderboard: User[] }>('/leaderboard'),
};

// Submissions API
export const submissionsApi = {
  runSample: (data: { problem_id: number; language: string; code: string }) =>
    apiClient.post<{ result: ExecutionResult }>('/submissions/run', data),
  submitMatch: (data: { match_public_id: string; language: string; code: string }) =>
    apiClient.post<{ submission: Submission; result: ExecutionResult }>('/submissions/submit', data),
};

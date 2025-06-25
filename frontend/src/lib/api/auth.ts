import api from './client';
import type { LoginResponse } from '@/types/auth';

export function login(data: { username: string; password: string }) {
  return api.post<LoginResponse>('/auth/login', data);
}

export function logout() {
  return api.post('/auth/logout');
}

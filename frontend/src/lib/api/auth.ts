import api from './client';
import axios from 'axios';
import type { LoginResponse } from '@/types/auth';

export function login(data: { username: string; password: string }) {
  // Call our Next.js proxy route which will forward the request
  // to the backend API defined via NEXT_PUBLIC_API_BASE_URL.
  return axios.post<LoginResponse>('/api/auth/login', data);
}

export function logout() {
  return api.post('/auth/logout');
}

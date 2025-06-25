'use client';
import { createContext, useContext } from 'react';
import type { User } from '@/types/auth';

export const AuthContext = createContext<User | null>(null);
export function useAuthContext() {
  return useContext(AuthContext);
}

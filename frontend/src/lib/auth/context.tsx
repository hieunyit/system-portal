'use client';
import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import * as authAPI from '@/lib/api/auth';
import tokenStorage from './storage';
import { getUser } from '@/lib/auth';
import type { User } from '@/types/auth';

export interface AuthContextValue {
  user: User | null;
  loading: boolean;
  login: (cred: { username: string; password: string }) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return ctx;
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const stored = getUser();
    if (stored) {
      setUser(stored as unknown as User);
    }
  }, []);

  const login = async (cred: { username: string; password: string }) => {
    setLoading(true);
    try {
      const res = await authAPI.login(cred);
      const body: any = res.data;
      if (body?.success?.data) {
        const { accessToken, refreshToken, user } = body.success.data;
        tokenStorage.setTokens(accessToken, refreshToken);
        setUser(user);
      } else if (body.access_token && body.refresh_token) {
        tokenStorage.setTokens(body.access_token, body.refresh_token);
        setUser(body.user);
      }
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    setLoading(true);
    try {
      await authAPI.logout();
    } catch {
      // ignore
    }
    tokenStorage.clearTokens();
    setUser(null);
    setLoading(false);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export default AuthContext;

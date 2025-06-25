'use client';
import { ReactNode, useEffect } from 'react';
import useAuth from '@/lib/hooks/useAuth';
import tokenStorage from '@/lib/auth/storage';

export default function AuthGuard({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();
  useEffect(() => {
    const token = tokenStorage.getAccessToken();
    if (!loading && !user && !token) {
      window.location.href = '/login';
    }
  }, [loading, user]);
  const token = tokenStorage.getAccessToken();
  if (!user && !token) return null;
  return <>{children}</>;
}

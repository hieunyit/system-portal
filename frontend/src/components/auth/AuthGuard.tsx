'use client';
import { ReactNode, useEffect } from 'react';
import useAuth from '@/lib/hooks/useAuth';

export default function AuthGuard({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();
  useEffect(() => {
    if (!loading && !user) {
      window.location.href = '/login';
    }
  }, [loading, user]);
  if (!user) return null;
  return <>{children}</>;
}

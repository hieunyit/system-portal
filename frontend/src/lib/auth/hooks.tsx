'use client';
import { useEffect, useState } from 'react';
import { login as apiLogin } from '@/lib/api/auth';
import tokenStorage from './storage';
import { AuthContext } from './context';
import type { User } from '@/types/auth';

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = tokenStorage.getAccessToken();
    if (!token) {
      setLoading(false);
      return;
    }
    // Placeholder: fetch profile
    setUser({ id: '1', username: 'demo', email: '', full_name: 'Demo' });
    setLoading(false);
  }, []);

  return { user, loading };
}

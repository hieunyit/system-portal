'use client';
import { useAuth } from './hooks';

export function usePermissions() {
  const { user } = useAuth();
  const hasPermission = (perm: string) => user?.permissions?.includes(perm) ?? false;
  return { hasPermission };
}

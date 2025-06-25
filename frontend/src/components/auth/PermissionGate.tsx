'use client';
import { ReactNode } from 'react';
import { usePermissions } from '@/lib/hooks/usePermissions';

export default function PermissionGate({ permission, children }: { permission: string; children: ReactNode }) {
  const { hasPermission } = usePermissions();
  if (!hasPermission(permission)) {
    return null;
  }
  return <>{children}</>;
}

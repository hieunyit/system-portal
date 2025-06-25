import api from '../client';
import type { Permission } from '@/types/portal';

export function listPermissions() {
  return api.get<Permission[]>('/api/portal/permissions');
}

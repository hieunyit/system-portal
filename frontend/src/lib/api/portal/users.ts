import api from '../client';
import type { PortalUser } from '@/types/portal';
import type { Pagination } from '@/types/global';

export function listUsers() {
  return api.get<Pagination<PortalUser>>('/api/portal/users');
}

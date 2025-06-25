import api from '../client';
import type { PortalGroup } from '@/types/portal';
import type { Pagination } from '@/types/global';

export function listGroups() {
  return api.get<Pagination<PortalGroup>>('/api/portal/groups');
}

import api from '../client';
import type { Pagination } from '@/types/global';

export function listLogs() {
  return api.get<Pagination<any>>('/api/portal/audit/logs');
}

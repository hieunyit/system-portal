import api from '../client';
import type { VpnUser } from '@/types/openvpn';
import type { Pagination } from '@/types/global';

export function listVpnUsers() {
  return api.get<Pagination<VpnUser>>('/api/openvpn/users');
}

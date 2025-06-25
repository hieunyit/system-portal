import api from '../client';
import type { VpnGroup } from '@/types/openvpn';
import type { Pagination } from '@/types/global';

export function listVpnGroups() {
  return api.get<Pagination<VpnGroup>>('/api/openvpn/groups');
}

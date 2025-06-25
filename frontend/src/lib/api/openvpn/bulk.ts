import api from '../client';

export function importUsers(data: FormData) {
  return api.post('/api/openvpn/bulk/import', data);
}

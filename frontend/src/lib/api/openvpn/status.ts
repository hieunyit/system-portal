import api from '../client';

export function getStatus() {
  return api.get('/api/openvpn/status');
}

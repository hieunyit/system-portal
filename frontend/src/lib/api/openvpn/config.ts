import api from '../client';

export function getVpnConfig() {
  return api.get('/api/openvpn/config');
}

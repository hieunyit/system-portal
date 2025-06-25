import api from '../client';

export function disconnectUser(username: string) {
  return api.post(`/api/openvpn/disconnect/${username}`);
}

import api from '../client';

export function getConfig() {
  return api.get('/api/portal/config');
}

'use client';
import api from '@/lib/api/client';
import tokenStorage from '@/lib/auth/storage';

export default function LogoutButton() {
  const logout = async () => {
    await api.post('/auth/logout');
    tokenStorage.clearTokens();
    window.location.href = '/login';
  };

  return (
    <button onClick={logout} className="px-3 py-2 bg-gray-200">Logout</button>
  );
}

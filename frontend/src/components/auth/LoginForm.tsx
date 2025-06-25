'use client';
import { useState } from 'react';
import api from '@/lib/api/client';
import tokenStorage from '@/lib/auth/storage';

export default function LoginForm() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await api.post('/auth/login', { username, password });
      tokenStorage.setTokens(res.data.access_token, res.data.refresh_token);
      window.location.href = '/dashboard';
    } catch (err) {
      setError('Login failed');
    }
  };

  return (
    <form onSubmit={submit} className="space-y-2">
      {error && <p className="text-red-500 text-sm">{error}</p>}
      <input value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" className="border p-2 w-full" />
      <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" className="border p-2 w-full" />
      <button className="bg-blue-600 text-white px-3 py-2 w-full">Login</button>
    </form>
  );
}

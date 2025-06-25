'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import useAuth from '@/lib/hooks/useAuth';

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await login({ username, password });
      router.push('/dashboard');
    } catch (err) {
      setError('Login failed');
    }
  };

  return (
    <div className="w-full max-w-md">
      <form
        onSubmit={handleSubmit}
        className="space-y-6 rounded-lg border bg-card p-8 shadow-lg"
      >
        <h1 className="text-2xl font-bold text-center text-primary">
          System Portal Login
        </h1>
        {error && <p className="text-destructive text-sm">{error}</p>}
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="w-full rounded border border-input bg-background p-2 focus:ring focus:ring-primary"
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full rounded border border-input bg-background p-2 focus:ring focus:ring-primary"
        />
        <button
          type="submit"
          className="w-full rounded bg-primary p-2 font-medium text-primary-foreground hover:bg-primary/90"
        >
          Login
        </button>
      </form>
    </div>
  );
}

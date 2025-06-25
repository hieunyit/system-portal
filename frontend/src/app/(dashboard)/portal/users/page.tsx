'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';

interface PortalUser {
  id: string;
  username: string;
  email: string;
  full_name: string;
}

export default function PortalUsersPage() {
  const [users, setUsers] = useState<PortalUser[]>([]);

  useEffect(() => {
    async function fetchUsers() {
      const res = await api.get('/api/portal/users');
      setUsers(res.data.data || []);
    }
    fetchUsers();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">Portal Users</h1>
      <table className="min-w-full border text-sm">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-3 py-2 text-left">Username</th>
            <th className="px-3 py-2 text-left">Email</th>
            <th className="px-3 py-2 text-left">Full Name</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id} className="border-t">
              <td className="px-3 py-2">{u.username}</td>
              <td className="px-3 py-2">{u.email}</td>
              <td className="px-3 py-2">{u.full_name}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

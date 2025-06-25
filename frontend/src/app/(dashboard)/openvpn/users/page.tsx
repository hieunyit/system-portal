'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';

interface VpnUser {
  username: string;
  group: string;
  status: string;
}

export default function VpnUsersPage() {
  const [users, setUsers] = useState<VpnUser[]>([]);

  useEffect(() => {
    async function fetchUsers() {
      const res = await api.get('/api/openvpn/users');
      setUsers(res.data.data || []);
    }
    fetchUsers();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">VPN Users</h1>
      <table className="min-w-full border">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-3 py-2 text-left">Username</th>
            <th className="px-3 py-2 text-left">Group</th>
            <th className="px-3 py-2 text-left">Status</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.username} className="border-t">
              <td className="px-3 py-2">{u.username}</td>
              <td className="px-3 py-2">{u.group}</td>
              <td className="px-3 py-2">{u.status}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

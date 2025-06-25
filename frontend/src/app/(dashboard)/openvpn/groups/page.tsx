'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';

interface VpnGroup {
  group_name: string;
  auth_method: string;
  mfa: boolean;
}

export default function VpnGroupsPage() {
  const [groups, setGroups] = useState<VpnGroup[]>([]);

  useEffect(() => {
    async function fetchGroups() {
      const res = await api.get('/api/openvpn/groups');
      setGroups(res.data.data || []);
    }
    fetchGroups();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">VPN Groups</h1>
      <table className="min-w-full border text-sm">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-3 py-2 text-left">Group Name</th>
            <th className="px-3 py-2 text-left">Auth Method</th>
            <th className="px-3 py-2 text-left">MFA</th>
          </tr>
        </thead>
        <tbody>
          {groups.map((g) => (
            <tr key={g.group_name} className="border-t">
              <td className="px-3 py-2">{g.group_name}</td>
              <td className="px-3 py-2">{g.auth_method}</td>
              <td className="px-3 py-2">{g.mfa ? 'Yes' : 'No'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';

interface PortalGroup {
  id: string;
  name: string;
  display_name: string;
}

export default function PortalGroupsPage() {
  const [groups, setGroups] = useState<PortalGroup[]>([]);

  useEffect(() => {
    async function fetchGroups() {
      const res = await api.get('/api/portal/groups');
      setGroups(res.data.data || []);
    }
    fetchGroups();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">Portal Groups</h1>
      <table className="min-w-full border text-sm">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-3 py-2 text-left">Name</th>
            <th className="px-3 py-2 text-left">Display Name</th>
          </tr>
        </thead>
        <tbody>
          {groups.map((g) => (
            <tr key={g.id} className="border-t">
              <td className="px-3 py-2">{g.name}</td>
              <td className="px-3 py-2">{g.display_name}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

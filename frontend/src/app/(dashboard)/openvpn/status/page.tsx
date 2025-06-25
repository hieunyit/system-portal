'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';

interface Session {
  username: string;
  ip: string;
  connected_since: string;
}

export default function VpnStatusPage() {
  const [sessions, setSessions] = useState<Session[]>([]);

  useEffect(() => {
    async function fetchStatus() {
      const res = await api.get('/api/openvpn/status');
      setSessions(res.data.active_sessions || []);
    }
    fetchStatus();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">Active VPN Sessions</h1>
      <table className="min-w-full border text-sm">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-3 py-2 text-left">Username</th>
            <th className="px-3 py-2 text-left">IP Address</th>
            <th className="px-3 py-2 text-left">Connected Since</th>
          </tr>
        </thead>
        <tbody>
          {sessions.map((s) => (
            <tr key={`${s.username}-${s.ip}`} className="border-t">
              <td className="px-3 py-2">{s.username}</td>
              <td className="px-3 py-2">{s.ip}</td>
              <td className="px-3 py-2">{s.connected_since}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

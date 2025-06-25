'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';
import PermissionTable from './components/PermissionTable';

interface Permission { id: string; resource: string; action: string; }

export default function PermissionsPage() {
  const [perms, setPerms] = useState<Permission[]>([]);

  useEffect(() => {
    api.get('/api/portal/permissions').then((res) => setPerms(res.data || []));
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">Permissions</h1>
      <PermissionTable permissions={perms} />
    </div>
  );
}

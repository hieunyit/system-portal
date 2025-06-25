'use client';
import { useEffect, useState } from 'react';
import api from '@/lib/api/client';
import AuditLogTable from '../components/AuditLogTable';
import AuditFilters from '../components/AuditFilters';
import AuditExport from '../components/AuditExport';

interface Log { id: string; username: string; action: string; ip: string; }

export default function AuditLogsPage() {
  const [logs, setLogs] = useState<Log[]>([]);

  useEffect(() => {
    api.get('/api/portal/audit/logs').then((res) => setLogs(res.data.data || []));
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">Audit Logs</h1>
      <AuditFilters />
      <AuditLogTable logs={logs} />
      <div className="mt-2">
        <AuditExport />
      </div>
    </div>
  );
}

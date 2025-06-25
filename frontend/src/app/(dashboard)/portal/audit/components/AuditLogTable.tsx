interface Log {
  id: string;
  username: string;
  action: string;
  ip: string;
}

export default function AuditLogTable({ logs }: { logs: Log[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">User</th>
          <th className="px-2 py-1 text-left">Action</th>
          <th className="px-2 py-1 text-left">IP</th>
        </tr>
      </thead>
      <tbody>
        {logs.map((l) => (
          <tr key={l.id} className="border-t">
            <td className="px-2 py-1">{l.username}</td>
            <td className="px-2 py-1">{l.action}</td>
            <td className="px-2 py-1">{l.ip}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

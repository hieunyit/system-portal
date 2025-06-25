interface Session {
  username: string;
  ip: string;
}

export default function ConnectionTable({ sessions }: { sessions: Session[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">User</th>
          <th className="px-2 py-1 text-left">IP</th>
        </tr>
      </thead>
      <tbody>
        {sessions.map((s, i) => (
          <tr key={i} className="border-t">
            <td className="px-2 py-1">{s.username}</td>
            <td className="px-2 py-1">{s.ip}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

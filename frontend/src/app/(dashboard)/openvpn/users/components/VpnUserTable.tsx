interface VpnUser {
  username: string;
  group: string;
}

export default function VpnUserTable({ users }: { users: VpnUser[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">Username</th>
          <th className="px-2 py-1 text-left">Group</th>
        </tr>
      </thead>
      <tbody>
        {users.map((u) => (
          <tr key={u.username} className="border-t">
            <td className="px-2 py-1">{u.username}</td>
            <td className="px-2 py-1">{u.group}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

interface PortalUser {
  id: string;
  username: string;
  email: string;
}

export default function UserTable({ users }: { users: PortalUser[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">Username</th>
          <th className="px-2 py-1 text-left">Email</th>
        </tr>
      </thead>
      <tbody>
        {users.map((u) => (
          <tr key={u.id} className="border-t">
            <td className="px-2 py-1">{u.username}</td>
            <td className="px-2 py-1">{u.email}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

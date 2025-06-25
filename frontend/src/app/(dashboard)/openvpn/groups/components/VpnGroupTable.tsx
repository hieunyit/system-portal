interface VpnGroup {
  group_name: string;
  auth_method: string;
}

export default function VpnGroupTable({ groups }: { groups: VpnGroup[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">Group Name</th>
          <th className="px-2 py-1 text-left">Auth Method</th>
        </tr>
      </thead>
      <tbody>
        {groups.map((g) => (
          <tr key={g.group_name} className="border-t">
            <td className="px-2 py-1">{g.group_name}</td>
            <td className="px-2 py-1">{g.auth_method}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

interface PortalGroup {
  id: string;
  name: string;
  display_name: string;
}

export default function GroupTable({ groups }: { groups: PortalGroup[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">Name</th>
          <th className="px-2 py-1 text-left">Display Name</th>
        </tr>
      </thead>
      <tbody>
        {groups.map((g) => (
          <tr key={g.id} className="border-t">
            <td className="px-2 py-1">{g.name}</td>
            <td className="px-2 py-1">{g.display_name}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

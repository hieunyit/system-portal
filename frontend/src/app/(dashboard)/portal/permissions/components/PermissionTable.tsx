interface Permission {
  id: string;
  resource: string;
  action: string;
}

export default function PermissionTable({ permissions }: { permissions: Permission[] }) {
  return (
    <table className="min-w-full border text-sm">
      <thead className="bg-gray-100">
        <tr>
          <th className="px-2 py-1 text-left">Resource</th>
          <th className="px-2 py-1 text-left">Action</th>
        </tr>
      </thead>
      <tbody>
        {permissions.map((p) => (
          <tr key={p.id} className="border-t">
            <td className="px-2 py-1">{p.resource}</td>
            <td className="px-2 py-1">{p.action}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

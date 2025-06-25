interface PermissionDetailProps {
  params: { id: string };
}

export default async function PermissionDetail({ params }: PermissionDetailProps) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/portal/permissions/${params.id}`);
  const permission = await res.json();

  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold">Permission Detail</h1>
      <p className="mt-2">Resource: {permission.resource}</p>
      <p>Action: {permission.action}</p>
    </div>
  );
}

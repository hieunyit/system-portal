import PermissionMatrix from '../../components/PermissionMatrix';

export default function GroupPermissionsPage() {
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">Group Permissions</h1>
      <PermissionMatrix />
    </div>
  );
}

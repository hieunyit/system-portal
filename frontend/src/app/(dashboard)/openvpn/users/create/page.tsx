import VpnUserForm from '../components/VpnUserForm';

export default function CreateVpnUserPage() {
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">Create VPN User</h1>
      <VpnUserForm />
    </div>
  );
}

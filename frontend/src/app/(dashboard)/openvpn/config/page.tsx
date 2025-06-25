import VpnConfigForm from './components/VpnConfigForm';

export default function VpnConfigPage() {
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">OpenVPN Configuration</h1>
      <VpnConfigForm />
    </div>
  );
}

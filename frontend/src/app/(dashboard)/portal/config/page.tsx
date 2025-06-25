import ConfigForm from './components/ConfigForm';

export default function PortalConfigPage() {
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">System Configuration</h1>
      <ConfigForm />
    </div>
  );
}

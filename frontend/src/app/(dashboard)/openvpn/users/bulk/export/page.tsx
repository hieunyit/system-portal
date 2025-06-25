import BulkExport from '../../components/BulkExport';

export default function BulkExportPage() {
  return (
    <div className="p-4">
      <h1 className="text-xl font-semibold mb-4">Export Users</h1>
      <BulkExport />
    </div>
  );
}

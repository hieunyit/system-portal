export default function AuditFilters() {
  return (
    <div className="space-x-2">
      <input type="text" placeholder="User" className="border p-1" />
      <input type="text" placeholder="IP" className="border p-1" />
      <button className="bg-blue-600 text-white px-2 py-1">Filter</button>
    </div>
  );
}

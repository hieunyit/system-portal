export default function StatusBadge({ status }: { status: string }) {
  const color = status === 'active' ? 'green' : 'red';
  return <span className={`px-2 py-1 text-xs text-white bg-${color}-600 rounded`}>{status}</span>;
}

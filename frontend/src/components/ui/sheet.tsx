export function Sheet({ children }: { children: React.ReactNode }) {
  return <div className="fixed inset-y-0 left-0 w-64 bg-white shadow">{children}</div>;
}

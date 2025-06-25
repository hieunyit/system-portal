export function Alert({ children }: { children: React.ReactNode }) {
  return <div className="border-l-4 border-red-500 bg-red-50 p-2">{children}</div>;
}

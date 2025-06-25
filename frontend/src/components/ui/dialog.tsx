export function Dialog({ children }: { children: React.ReactNode }) {
  return <div className="fixed inset-0 bg-black/50 flex items-center justify-center">{children}</div>;
}

export function Form({ children, onSubmit }: { children: React.ReactNode; onSubmit?: () => void }) {
  return (
    <form onSubmit={onSubmit} className="space-y-2">
      {children}
    </form>
  );
}

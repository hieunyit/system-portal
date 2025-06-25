import { ReactNode, useState } from 'react';

export function DropdownMenu({ label, children }: { label: string; children: ReactNode }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="relative inline-block">
      <button onClick={() => setOpen(!open)} className="px-2 py-1 bg-gray-200">{label}</button>
      {open && <div className="absolute bg-white border mt-1">{children}</div>}
    </div>
  );
}

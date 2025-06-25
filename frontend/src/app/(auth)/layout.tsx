import type { ReactNode } from 'react';

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-accent to-primary/80 dark:from-gray-800 dark:to-gray-900 p-4">
      {children}
    </div>
  );
}

import type { ReactNode } from 'react';

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-accent dark:from-gray-900 dark:to-gray-700 p-4">
      {children}
    </div>
  );
}

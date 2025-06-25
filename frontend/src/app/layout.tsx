import './globals.css';
import type { ReactNode } from 'react';
import { AuthProvider } from '@/lib/auth/context';

export const metadata = {
  title: 'System Portal',
  description: 'System Portal Frontend',
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="min-h-screen font-sans">
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}

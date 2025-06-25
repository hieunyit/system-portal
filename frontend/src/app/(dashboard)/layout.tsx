import type { ReactNode } from 'react';

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen flex">
      <aside className="w-64 bg-gray-900 text-white p-4 hidden sm:block">
        <nav className="space-y-2">
          <a href="/dashboard" className="block">Dashboard</a>
          <a href="/portal/users" className="block">Portal Users</a>
          <a href="/portal/groups" className="block">Portal Groups</a>
          <a href="/openvpn/users" className="block">VPN Users</a>
          <a href="/openvpn/groups" className="block">VPN Groups</a>
          <a href="/openvpn/status" className="block">VPN Status</a>
        </nav>
      </aside>
      <main className="flex-1 p-4">{children}</main>
    </div>
  );
}

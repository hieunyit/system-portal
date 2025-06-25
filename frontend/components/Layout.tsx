import Link from 'next/link'
import { ReactNode } from 'react'

export default function Layout({ children }: { children: ReactNode }) {
  const linkStyle = { marginRight: '1rem' }
  return (
    <div style={{fontFamily:'Arial'}}>
      <nav style={{background:'#333',padding:'0.5rem'}}>
        <Link href="/" style={{...linkStyle, color:'#fff'}}>Home</Link>
        <Link href="/login" style={{...linkStyle, color:'#fff'}}>Login</Link>
        <Link href="/openvpn/status" style={{...linkStyle, color:'#fff'}}>VPN Status</Link>
        <Link href="/openvpn/users" style={{...linkStyle, color:'#fff'}}>VPN Users</Link>
        <Link href="/openvpn/groups" style={{...linkStyle, color:'#fff'}}>VPN Groups</Link>
        <Link href="/openvpn/network" style={{...linkStyle, color:'#fff'}}>Network</Link>
        <Link href="/openvpn/server" style={{...linkStyle, color:'#fff'}}>Server Info</Link>
        <Link href="/portal/users" style={{...linkStyle, color:'#fff'}}>Portal Users</Link>
        <Link href="/portal/groups" style={{...linkStyle, color:'#fff'}}>Portal Groups</Link>
        <Link href="/portal/permissions" style={{...linkStyle, color:'#fff'}}>Permissions</Link>
        <Link href="/portal/ldap" style={{...linkStyle, color:'#fff'}}>LDAP</Link>
        <Link href="/portal/openvpn" style={{...linkStyle, color:'#fff'}}>OVPN Conn</Link>
        <Link href="/portal/smtp" style={{...linkStyle, color:'#fff'}}>SMTP</Link>
        <Link href="/portal/audit" style={{...linkStyle, color:'#fff'}}>Audit</Link>
        <Link href="/portal/templates" style={{...linkStyle, color:'#fff'}}>Templates</Link>
        <Link href="/dashboard/stats" style={{...linkStyle, color:'#fff'}}>Stats</Link>
        <Link href="/dashboard/activities" style={{...linkStyle, color:'#fff'}}>Activities</Link>
      </nav>
      <div style={{padding:'1rem'}}>{children}</div>
    </div>
  )
}

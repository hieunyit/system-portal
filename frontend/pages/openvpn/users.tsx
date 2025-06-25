import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface VpnUser {
  username: string
  email: string
  groupName: string
  isEnabled?: boolean
}

export default function VpnUsersPage() {
  const [users, setUsers] = useState<VpnUser[]>([])

  useEffect(() => {
    fetch('/api/openvpn/users')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setUsers(data.users || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>VPN Users</h1>
      <table style={{borderCollapse:'collapse'}}>
        <thead>
          <tr>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Username</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Email</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Group</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Enabled</th>
          </tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.username}>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.username}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.email}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.groupName}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.isEnabled ? 'Yes' : 'No'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </Layout>
  )
}

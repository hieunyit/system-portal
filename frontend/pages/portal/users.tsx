import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface PortalUser {
  id: string
  username: string
  email: string
  groupId?: string
  isActive?: boolean
}

export default function PortalUsersPage() {
  const [users, setUsers] = useState<PortalUser[]>([])

  useEffect(() => {
    fetch('/api/portal/users')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setUsers(data.data || data.users || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>Portal Users</h1>
      <table style={{borderCollapse:'collapse'}}>
        <thead>
          <tr>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Username</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Email</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Group</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Active</th>
          </tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id}>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.username}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.email}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.groupId}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{u.isActive ? 'Yes' : 'No'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </Layout>
  )
}

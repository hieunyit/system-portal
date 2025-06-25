import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface VpnGroup {
  groupName: string
  authMethod: string
  role?: string
}

export default function VpnGroupsPage() {
  const [groups, setGroups] = useState<VpnGroup[]>([])

  useEffect(() => {
    fetch('/api/openvpn/groups')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setGroups(data.groups || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>VPN Groups</h1>
      <table style={{borderCollapse:'collapse'}}>
        <thead>
          <tr>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Name</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Auth Method</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Role</th>
          </tr>
        </thead>
        <tbody>
          {groups.map(g => (
            <tr key={g.groupName}>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{g.groupName}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{g.authMethod}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{g.role}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </Layout>
  )
}

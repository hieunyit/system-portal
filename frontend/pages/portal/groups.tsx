import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface PortalGroup {
  id: string
  name: string
  displayName?: string
  isActive?: boolean
}

export default function PortalGroupsPage() {
  const [groups, setGroups] = useState<PortalGroup[]>([])

  useEffect(() => {
    fetch('/api/portal/groups')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setGroups(data.data || data.groups || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>Portal Groups</h1>
      <table style={{borderCollapse:'collapse'}}>
        <thead>
          <tr>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Name</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Display</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Active</th>
          </tr>
        </thead>
        <tbody>
          {groups.map(g => (
            <tr key={g.id}>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{g.name}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{g.displayName}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{g.isActive ? 'Yes' : 'No'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </Layout>
  )
}

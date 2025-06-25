import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface Permission {
  id: string
  resource: string
  action: string
}

export default function PermissionsPage() {
  const [perms, setPerms] = useState<Permission[]>([])

  useEffect(() => {
    fetch('/api/portal/permissions')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setPerms(data.data || data.permissions || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>Permissions</h1>
      <table style={{borderCollapse:'collapse'}}>
        <thead>
          <tr>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Resource</th>
            <th style={{border:'1px solid #ccc',padding:'4px'}}>Action</th>
          </tr>
        </thead>
        <tbody>
          {perms.map(p => (
            <tr key={p.id}>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{p.resource}</td>
              <td style={{border:'1px solid #ccc',padding:'4px'}}>{p.action}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </Layout>
  )
}

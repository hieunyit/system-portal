import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

export default function AuditStatsPage() {
  const [stats, setStats] = useState<any>(null)

  useEffect(() => {
    fetch('/api/portal/audit/stats')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setStats(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>Audit Stats</h1>
      {stats ? <pre>{JSON.stringify(stats, null, 2)}</pre> : <p>Loading...</p>}
    </Layout>
  )
}

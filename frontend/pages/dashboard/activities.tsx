import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

export default function DashboardActivitiesPage() {
  const [activities, setActivities] = useState<any[]>([])

  useEffect(() => {
    fetch('/api/portal/dashboard/activities')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setActivities(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>Recent Activities</h1>
      <ul>
        {activities.map((a, i) => (
          <li key={i}>{JSON.stringify(a)}</li>
        ))}
      </ul>
    </Layout>
  )
}

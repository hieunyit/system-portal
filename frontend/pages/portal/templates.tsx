import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

export default function EmailTemplatesPage() {
  const [templates, setTemplates] = useState<any[]>([])

  useEffect(() => {
    fetch('/api/portal/connections/templates/welcome')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setTemplates([data.data || data]))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>Email Templates</h1>
      {templates.map((t,i) => (
        <pre key={i}>{JSON.stringify(t, null, 2)}</pre>
      ))}
    </Layout>
  )
}

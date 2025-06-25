import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

export default function SMTPConfigPage() {
  const [config, setConfig] = useState<any>(null)

  useEffect(() => {
    fetch('/api/portal/connections/smtp')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setConfig(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>SMTP Configuration</h1>
      {config ? <pre>{JSON.stringify(config, null, 2)}</pre> : <p>Loading...</p>}
    </Layout>
  )
}

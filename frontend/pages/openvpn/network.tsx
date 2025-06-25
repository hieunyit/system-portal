import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

export default function NetworkConfigPage() {
  const [config, setConfig] = useState<any>(null)

  useEffect(() => {
    fetch('/api/openvpn/config/network')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setConfig(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>VPN Network Configuration</h1>
      {config ? <pre>{JSON.stringify(config, null, 2)}</pre> : <p>Loading...</p>}
    </Layout>
  )
}

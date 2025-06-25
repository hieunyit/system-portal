import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface ServerInfo {
  node_type?: string
  status?: string
  web_server_name?: string
}

export default function OpenVPNServerPage() {
  const [info, setInfo] = useState<ServerInfo | null>(null)

  useEffect(() => {
    fetch('/api/openvpn/config/server/info')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setInfo(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>OpenVPN Server Info</h1>
      {info ? (
        <pre>{JSON.stringify(info, null, 2)}</pre>
      ) : (
        <p>Loading...</p>
      )}
    </Layout>
  )
}

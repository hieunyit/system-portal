import { useEffect, useState } from 'react'
import Layout from '../../components/Layout'

interface VpnStatus {
  timestamp: string
  total_connected_users: number
  connected_users: any[]
}

export default function VpnStatusPage() {
  const [status, setStatus] = useState<VpnStatus | null>(null)

  useEffect(() => {
    fetch('/api/openvpn/vpn/status')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setStatus(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <h1>VPN Status</h1>
      {status ? (
        <div>
          <p>Last update: {status.timestamp}</p>
          <p>Total connected users: {status.total_connected_users}</p>
        </div>
      ) : (
        <p>Loading...</p>
      )}
    </Layout>
  )
}

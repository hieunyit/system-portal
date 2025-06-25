import Head from 'next/head'
import { useEffect, useState } from 'react'
import Layout from '../components/Layout'

interface VpnStatus {
  timestamp: string
  total_connected_users: number
}

export default function Home() {
  const [status, setStatus] = useState<VpnStatus | null>(null)

  useEffect(() => {
    fetch('/api/openvpn/vpn/status')
      .then(res => res.ok ? res.json() : Promise.reject('Failed to load'))
      .then(data => setStatus(data.data || data))
      .catch(err => console.error(err))
  }, [])

  return (
    <Layout>
      <Head>
        <title>System Portal</title>
      </Head>
      <h1>System Portal Dashboard</h1>
      {status ? (
        <div>
          <p>Last updated: {status.timestamp}</p>
          <p>Connected users: {status.total_connected_users}</p>
        </div>
      ) : (
        <p>Loading VPN status...</p>
      )}
    </Layout>
  )
}

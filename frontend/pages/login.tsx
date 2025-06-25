import { useState } from 'react'
import { useRouter } from 'next/router'
import Layout from '../components/Layout'

export default function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    try {
      const res = await fetch('/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      })
      if (!res.ok) throw new Error('Login failed')
      const data = await res.json()
      localStorage.setItem('accessToken', data.accessToken)
      localStorage.setItem('refreshToken', data.refreshToken)
      router.push('/')
    } catch (err) {
      setError((err as Error).message)
    }
  }

  return (
    <Layout>
      <h1>Login</h1>
      <form onSubmit={handleSubmit} style={{display:'flex',flexDirection:'column',maxWidth:'300px'}}>
        <input value={username} onChange={e=>setUsername(e.target.value)} placeholder="Username" required />
        <input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="Password" required />
        <button type="submit" style={{marginTop:'1rem'}}>Login</button>
        {error && <p style={{color:'red'}}>{error}</p>}
      </form>
    </Layout>
  )
}

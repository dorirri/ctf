import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import client from '../api/client'

export default function Register() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ username: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function handleChange(e) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    if (form.password.length < 8) {
      setError('Password must be at least 8 characters.')
      return
    }
    setLoading(true)
    try {
      await client.post('/auth/register', form)
      navigate('/login', { replace: true })
    } catch (err) {
      setError(err.response?.data?.error ?? 'Registration failed.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="form-card">
      <h1>$ register</h1>
      <p className="subtitle">Create an account to join the competition.</p>

      {error && <div className="alert alert-error">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            id="username"
            name="username"
            type="text"
            placeholder="h4ck3r"
            value={form.username}
            onChange={handleChange}
            required
            autoFocus
          />
        </div>
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            placeholder="you@example.com"
            value={form.email}
            onChange={handleChange}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            placeholder="min. 8 characters"
            value={form.password}
            onChange={handleChange}
            required
          />
        </div>
        <button className="btn btn-primary form-submit" type="submit" disabled={loading}>
          {loading ? 'Creating account...' : 'Register'}
        </button>
      </form>

      <p className="form-footer">
        Already have an account?{' '}
        <Link to="/login">Login here</Link>
      </p>
    </div>
  )
}

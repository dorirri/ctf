import { useState, useEffect } from 'react'
import client from '../../api/client'
import { useAuth } from '../../context/AuthContext'

function formatTime(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric',
  })
}

export default function AdminUsers() {
  const { user: self } = useAuth()
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [toggling, setToggling] = useState(null)

  useEffect(() => {
    client.get('/admin/users')
      .then(({ data }) => setUsers(data))
      .catch(() => setError('Failed to load users.'))
      .finally(() => setLoading(false))
  }, [])

  async function handleToggleDisable(u) {
    const action = u.is_disabled ? 'enable' : 'disable'
    if (!window.confirm(`${action} user "${u.username}"?`)) return
    setToggling(u.id)
    try {
      await client.patch(`/admin/users/${u.id}/disable`)
      setUsers((prev) =>
        prev.map((x) => (x.id === u.id ? { ...x, is_disabled: !x.is_disabled } : x))
      )
    } catch (err) {
      alert(err.response?.data?.error ?? 'Toggle failed.')
    } finally {
      setToggling(null)
    }
  }

  return (
    <>
      <div className="admin-page-header">
        <h1>// Users</h1>
        <span style={{ color: 'var(--text-dim)', fontSize: '0.8rem' }}>
          {users.length} registered
        </span>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      {loading ? (
        <div className="state-box">Loading...</div>
      ) : users.length === 0 ? (
        <div className="state-box">No users found.</div>
      ) : (
        <div className="admin-table-wrap">
          <table className="admin-table">
            <thead>
              <tr>
                <th>Username</th>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Joined</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => {
                const isSelf = u.id === self?.id
                return (
                  <tr key={u.id}>
                    <td>
                      {u.username}
                      {isSelf && (
                        <span style={{ color: 'var(--text-dim)', fontSize: '0.7rem', marginLeft: '0.4rem' }}>
                          (you)
                        </span>
                      )}
                    </td>
                    <td style={{ color: 'var(--text-dim)' }}>{u.email}</td>
                    <td>
                      <span
                        style={{
                          color: u.role === 'admin' ? 'var(--admin)' : 'var(--text-dim)',
                          fontSize: '0.78rem',
                        }}
                      >
                        {u.role}
                      </span>
                    </td>
                    <td>
                      <span className={`badge ${u.is_disabled ? 'badge-disabled' : 'badge-active'}`}>
                        {u.is_disabled ? 'disabled' : 'active'}
                      </span>
                    </td>
                    <td style={{ color: 'var(--text-dim)', fontSize: '0.78rem' }}>
                      {formatTime(u.created_at)}
                    </td>
                    <td>
                      <button
                        className={`btn btn-sm ${u.is_disabled ? 'btn-ghost' : 'btn-danger'}`}
                        onClick={() => handleToggleDisable(u)}
                        disabled={toggling === u.id || isSelf || u.role === 'admin'}
                        title={isSelf ? "Can't disable yourself" : u.role === 'admin' ? "Can't disable admin" : ''}
                      >
                        {toggling === u.id ? '...' : u.is_disabled ? 'Enable' : 'Disable'}
                      </button>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}

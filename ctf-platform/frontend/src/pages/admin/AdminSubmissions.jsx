import { useState, useEffect } from 'react'
import client from '../../api/client'

function formatTime(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString(undefined, {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
}

const FILTERS = ['all', 'correct', 'wrong']

export default function AdminSubmissions() {
  const [submissions, setSubmissions] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [filter, setFilter] = useState('all')

  useEffect(() => {
    client.get('/admin/submissions')
      .then(({ data }) => setSubmissions(data))
      .catch(() => setError('Failed to load submissions.'))
      .finally(() => setLoading(false))
  }, [])

  const visible = submissions.filter((s) => {
    if (filter === 'correct') return s.is_correct
    if (filter === 'wrong') return !s.is_correct
    return true
  })

  return (
    <>
      <div className="admin-page-header">
        <h1>// Submissions</h1>
        <span style={{ color: 'var(--text-dim)', fontSize: '0.8rem' }}>
          {submissions.length} total
        </span>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      <div className="admin-filter-bar">
        {FILTERS.map((f) => (
          <button
            key={f}
            className={`filter-btn${filter === f ? ' active' : ''}`}
            onClick={() => setFilter(f)}
          >
            {f}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="state-box">Loading...</div>
      ) : visible.length === 0 ? (
        <div className="state-box">No submissions match this filter.</div>
      ) : (
        <div className="admin-table-wrap">
          <table className="admin-table">
            <thead>
              <tr>
                <th>Player</th>
                <th>Challenge</th>
                <th>Result</th>
                <th>Submitted At</th>
              </tr>
            </thead>
            <tbody>
              {visible.map((s, i) => (
                <tr key={i}>
                  <td>{s.username}</td>
                  <td>{s.challenge_title}</td>
                  <td>
                    <span className={`badge ${s.is_correct ? 'badge-correct' : 'badge-wrong'}`}>
                      {s.is_correct ? 'correct' : 'wrong'}
                    </span>
                  </td>
                  <td style={{ color: 'var(--text-dim)', fontSize: '0.78rem' }}>
                    {formatTime(s.submitted_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}

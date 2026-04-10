import { useState, useEffect } from 'react'
import client from '../api/client'

function formatTime(iso) {
  if (!iso) return '—'
  const d = new Date(iso)
  return d.toLocaleString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function rankClass(rank) {
  if (rank === 1) return 'rank-gold'
  if (rank === 2) return 'rank-silver'
  if (rank === 3) return 'rank-bronze'
  return ''
}

export default function Scoreboard() {
  const [entries, setEntries] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    client.get('/scoreboard')
      .then(({ data }) => setEntries(data))
      .catch(() => setError('Failed to load scoreboard.'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="scoreboard-page">
      <h1>// Scoreboard</h1>
      <p className="page-subtitle">Top 20 players</p>

      {error && <div className="alert alert-error">{error}</div>}

      {loading && (
        <div className="state-box">Loading...</div>
      )}

      {!loading && !error && entries.length === 0 && (
        <div className="state-box">
          <span className="state-icon">[]</span>
          No solves yet. Be the first!
        </div>
      )}

      {!loading && entries.length > 0 && (
        <table className="scoreboard-table">
          <thead>
            <tr>
              <th>Rank</th>
              <th>Player</th>
              <th>Points</th>
              <th>Solves</th>
              <th>Last Solve</th>
            </tr>
          </thead>
          <tbody>
            {entries.map((entry) => (
              <tr key={entry.rank}>
                <td className={rankClass(entry.rank)}>#{entry.rank}</td>
                <td className="username-cell">{entry.username}</td>
                <td className="points-cell">{entry.total_points}</td>
                <td>{entry.solves_count}</td>
                <td className="time-cell">{formatTime(entry.last_solve_time)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

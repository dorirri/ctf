import { useState, useEffect, useCallback } from 'react'
import client from '../api/client'
import ChallengeCard from '../components/ChallengeCard'
import FlagSubmitModal from '../components/FlagSubmitModal'

export default function Challenges() {
  const [challenges, setChallenges] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selected, setSelected] = useState(null)

  const fetchChallenges = useCallback(async () => {
    try {
      const { data } = await client.get('/challenges')
      setChallenges(data)
    } catch {
      setError('Failed to load challenges.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetchChallenges() }, [fetchChallenges])

  function handleSolved(challengeId) {
    setChallenges((prev) =>
      prev.map((c) => (c.id === challengeId ? { ...c, is_solved: true } : c))
    )
  }

  function openModal(challenge) {
    // Fetch full challenge (with description) then open modal
    client.get(`/challenges/${challenge.id}`)
      .then(({ data }) => setSelected(data))
      .catch(() => setSelected(challenge))
  }

  const grouped = challenges.reduce((acc, ch) => {
    ;(acc[ch.category] ??= []).push(ch)
    return acc
  }, {})

  const categories = Object.keys(grouped).sort()
  const solvedCount = challenges.filter((c) => c.is_solved).length

  return (
    <div className="challenges-page">
      <h1>// Challenges</h1>
      <p className="page-subtitle">
        {loading
          ? 'Loading...'
          : `${challenges.length} challenges — ${solvedCount} solved`}
      </p>

      {error && <div className="alert alert-error">{error}</div>}

      {loading && (
        <div className="challenge-grid loading-grid" aria-label="Loading challenges">
          {Array.from({ length: 6 }).map((_, index) => (
            <div className="challenge-card challenge-card-skeleton" key={index}>
              <span />
              <span />
            </div>
          ))}
        </div>
      )}

      {!loading && !error && categories.length === 0 && (
        <div className="state-box">
          <span className="state-icon">[]</span>
          No challenges published yet.
        </div>
      )}

      {categories.map((cat) => (
        <div className="category-section" key={cat}>
          <div className="category-title">{cat}</div>
          <div className="challenge-grid">
            {grouped[cat].map((ch) => (
              <ChallengeCard key={ch.id} challenge={ch} onClick={openModal} />
            ))}
          </div>
        </div>
      ))}

      {selected && (
        <FlagSubmitModal
          challenge={selected}
          onClose={() => setSelected(null)}
          onSolved={(id) => {
            handleSolved(id)
            setSelected((prev) => prev ? { ...prev, is_solved: true } : null)
          }}
        />
      )}
    </div>
  )
}

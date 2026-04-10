import { useState, useEffect, useRef } from 'react'
import client from '../api/client'

export default function FlagSubmitModal({ challenge, onClose, onSolved }) {
  const [flag, setFlag] = useState('')
  const [status, setStatus] = useState(null) // { type: 'correct'|'wrong'|'already-solved'|'error', msg }
  const [loading, setLoading] = useState(false)
  const inputRef = useRef(null)

  useEffect(() => {
    inputRef.current?.focus()
    const onKey = (e) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  async function handleSubmit(e) {
    e.preventDefault()
    if (!flag.trim() || loading) return
    setLoading(true)
    setStatus(null)
    try {
      const { data } = await client.post('/submit', {
        challenge_id: challenge.id,
        flag: flag.trim(),
      })
      if (data.correct) {
        setStatus({ type: 'correct', msg: `Correct! +${data.points} points` })
        onSolved(challenge.id)
      } else {
        setStatus({ type: 'wrong', msg: 'Wrong flag. Try again.' })
      }
    } catch (err) {
      const msg = err.response?.data?.error ?? 'Submission failed.'
      if (err.response?.status === 409) {
        setStatus({ type: 'already-solved', msg: 'Already solved.' })
      } else {
        setStatus({ type: 'error', msg })
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={(e) => e.target === e.currentTarget && onClose()}>
      <div className="modal">
        <div className="modal-header">
          <div>
            <div className="modal-title">{challenge.title}</div>
            <div className="modal-meta">
              <span className="modal-category">{challenge.category}</span>
              <span className="modal-points">{challenge.points} pts</span>
            </div>
          </div>
          <button className="modal-close" onClick={onClose} aria-label="Close">×</button>
        </div>

        {challenge.description && (
          <div className="modal-description">{challenge.description}</div>
        )}

        {challenge.is_solved ? (
          <div className="modal-result already-solved">Challenge already solved.</div>
        ) : (
          <form className="modal-form" onSubmit={handleSubmit}>
            <input
              ref={inputRef}
              className="modal-flag-input"
              type="text"
              placeholder="CTF{...}"
              value={flag}
              onChange={(e) => setFlag(e.target.value)}
              disabled={loading}
            />
            <button className="btn btn-primary" type="submit" disabled={loading || !flag.trim()}>
              {loading ? '...' : 'Submit'}
            </button>
          </form>
        )}

        {status && (
          <div className={`modal-result ${status.type === 'error' ? 'wrong' : status.type}`}>
            {status.msg}
          </div>
        )}
      </div>
    </div>
  )
}

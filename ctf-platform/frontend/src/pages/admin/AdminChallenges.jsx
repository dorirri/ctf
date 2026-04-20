import { useState, useEffect, useCallback } from 'react'
import client from '../../api/client'

const CATEGORIES = ['web', 'crypto', 'forensics', 'misc', 'pwn', 'reverse']

const EMPTY_FORM = {
  title: '',
  description: '',
  category: 'web',
  points: '',
  flag: '',
  is_visible: true,
}

function ChallengeFormModal({ initial, onClose, onSaved }) {
  const isEdit = Boolean(initial)
  const [form, setForm] = useState(
    isEdit
      ? {
          title: initial.title,
          description: initial.description,
          category: initial.category.toLowerCase(),
          points: String(initial.points),
          flag: '',
          is_visible: initial.is_visible,
        }
      : { ...EMPTY_FORM }
  )
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function set(field, value) {
    setForm((f) => ({ ...f, [field]: value }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    const pts = parseInt(form.points, 10)
    if (!form.title.trim()) { setError('Title is required.'); return }
    if (!isEdit && !form.flag.trim()) { setError('Flag is required.'); return }
    if (isNaN(pts) || pts <= 0) { setError('Points must be a positive number.'); return }

    setLoading(true)
    try {
      const payload = {
        title: form.title.trim(),
        description: form.description.trim(),
        category: form.category,
        points: pts,
        is_visible: form.is_visible,
      }
      if (isEdit) {
        await client.put(`/admin/challenges/${initial.id}`, payload)
      } else {
        await client.post('/admin/challenges', { ...payload, flag: form.flag.trim() })
      }
      onSaved()
    } catch (err) {
      setError(err.response?.data?.error ?? 'Save failed.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    const onKey = (e) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  return (
    <div className="admin-modal-overlay" onClick={(e) => e.target === e.currentTarget && onClose()}>
      <div className="admin-modal">
        <h2>{isEdit ? '// Edit Challenge' : '// New Challenge'}</h2>

        {error && <div className="alert alert-error" style={{ marginBottom: '1rem' }}>{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="admin-form-grid">
            <div className="admin-form-group full">
              <label>Title</label>
              <input
                type="text"
                value={form.title}
                onChange={(e) => set('title', e.target.value)}
                placeholder="Challenge name"
                autoFocus
              />
            </div>

            <div className="admin-form-group full">
              <label>Description</label>
              <textarea
                value={form.description}
                onChange={(e) => set('description', e.target.value)}
                placeholder="Challenge description shown to players..."
              />
            </div>

            <div className="admin-form-group">
              <label>Category</label>
              <select value={form.category} onChange={(e) => set('category', e.target.value)}>
                {CATEGORIES.map((c) => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>

            <div className="admin-form-group">
              <label>Points</label>
              <input
                type="number"
                min="1"
                value={form.points}
                onChange={(e) => set('points', e.target.value)}
                placeholder="100"
              />
            </div>

            <div className="admin-form-group full">
              <label>{isEdit ? 'New Flag (leave blank to keep current)' : 'Flag'}</label>
              <input
                type="text"
                value={form.flag}
                onChange={(e) => set('flag', e.target.value)}
                placeholder="CTF{...}"
              />
            </div>

            <div className="admin-form-group full">
              <div className="admin-checkbox-row">
                <input
                  type="checkbox"
                  id="is_visible"
                  checked={form.is_visible}
                  onChange={(e) => set('is_visible', e.target.checked)}
                />
                <label htmlFor="is_visible">Visible to players</label>
              </div>
            </div>
          </div>

          <div className="admin-form-actions">
            <button type="button" className="btn btn-ghost btn-sm" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn btn-admin btn-sm" disabled={loading}>
              {loading ? 'Saving...' : isEdit ? 'Save Changes' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default function AdminChallenges() {
  const [challenges, setChallenges] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [modal, setModal] = useState(null) // null | 'new' | challenge object

  const fetch = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const { data } = await client.get('/admin/challenges')
      setChallenges(data)
    } catch {
      setError('Failed to load challenges.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetch() }, [fetch])

  async function handleDelete(ch) {
    if (!window.confirm(`Delete "${ch.title}"? This also removes all submissions for it.`)) return
    try {
      await client.delete(`/admin/challenges/${ch.id}`)
      fetch()
    } catch (err) {
      alert(err.response?.data?.error ?? 'Delete failed.')
    }
  }

  async function handleToggleVisibility(ch) {
    try {
      await client.patch(`/admin/challenges/${ch.id}/visibility`)
      setChallenges((prev) =>
        prev.map((c) => (c.id === ch.id ? { ...c, is_visible: !c.is_visible } : c))
      )
    } catch (err) {
      alert(err.response?.data?.error ?? 'Toggle failed.')
    }
  }

  return (
    <>
      <div className="admin-page-header">
        <h1>// Challenges</h1>
        <button className="btn btn-admin btn-sm" onClick={() => setModal('new')}>
          + New Challenge
        </button>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      {loading ? (
        <div className="state-box">Loading...</div>
      ) : challenges.length === 0 ? (
        <div className="state-box">No challenges yet.</div>
      ) : (
        <div className="admin-table-wrap">
          <table className="admin-table">
            <thead>
              <tr>
                <th>Title</th>
                <th>Category</th>
                <th>Points</th>
                <th>Visible</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {challenges.map((ch) => (
                <tr key={ch.id}>
                  <td>{ch.title}</td>
                  <td>{ch.category}</td>
                  <td>{ch.points}</td>
                  <td>
                    <span className={`badge ${ch.is_visible ? 'badge-visible' : 'badge-hidden'}`}>
                      {ch.is_visible ? 'visible' : 'hidden'}
                    </span>
                  </td>
                  <td>
                    <div className="actions-cell">
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => handleToggleVisibility(ch)}
                        title={ch.is_visible ? 'Hide' : 'Show'}
                      >
                        {ch.is_visible ? 'Hide' : 'Show'}
                      </button>
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => setModal(ch)}
                      >
                        Edit
                      </button>
                      <button
                        className="btn btn-danger btn-sm"
                        onClick={() => handleDelete(ch)}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {modal && (
        <ChallengeFormModal
          initial={modal === 'new' ? null : modal}
          onClose={() => setModal(null)}
          onSaved={() => { setModal(null); fetch() }}
        />
      )}
    </>
  )
}

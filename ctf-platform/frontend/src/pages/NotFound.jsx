import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="not-found">
      <div className="code">404</div>
      <div className="message">Page not found.</div>
      <Link to="/" className="btn btn-ghost" style={{ marginTop: '1rem' }}>
        ← Go home
      </Link>
    </div>
  )
}

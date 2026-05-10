import { NavLink, Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-logo">
        CTF<span>://</span>Platform
      </Link>

      <div className="navbar-links">
        <NavLink to="/challenges">Challenges</NavLink>
        <NavLink to="/scoreboard">Scoreboard</NavLink>
      </div>

      <div className="navbar-right">
        {user ? (
          <>
            <span className="navbar-username">{user.username}</span>
            <button className="btn btn-ghost" onClick={handleLogout}>
              Logout
            </button>
          </>
        ) : (
          <>
            <NavLink to="/login" className="btn btn-ghost">Login</NavLink>
            <NavLink to="/register" className="btn btn-primary">Register</NavLink>
          </>
        )}
      </div>
    </nav>
  )
}

import { NavLink, Link, Outlet } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

export default function AdminLayout() {
  const { user } = useAuth()

  return (
    <div className="admin-shell">
      <aside className="admin-sidebar">
        <div className="admin-sidebar-brand">
          Admin Panel
          <span>{user?.username}</span>
        </div>

        <nav className="admin-nav">
          <NavLink to="/admin/challenges">Challenges</NavLink>
          <NavLink to="/admin/submissions">Submissions</NavLink>
          <NavLink to="/admin/users">Users</NavLink>
          <div className="admin-nav-divider" />
          <Link to="/challenges" className="admin-back-link">← Back to Platform</Link>
        </nav>
      </aside>

      <div className="admin-content">
        <Outlet />
      </div>
    </div>
  )
}

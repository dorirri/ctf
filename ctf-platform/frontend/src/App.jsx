import { Routes, Route, Navigate, Outlet } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Navbar from './components/Navbar'
import ProtectedRoute from './components/ProtectedRoute'
import AdminRoute from './components/AdminRoute'
import Login from './pages/Login'
import Register from './pages/Register'
import Challenges from './pages/Challenges'
import Scoreboard from './pages/Scoreboard'
import NotFound from './pages/NotFound'
import AdminLayout from './pages/admin/AdminLayout'
import AdminChallenges from './pages/admin/AdminChallenges'
import AdminSubmissions from './pages/admin/AdminSubmissions'
import AdminUsers from './pages/admin/AdminUsers'

function PlayerShell() {
  return (
    <div className="page">
      <Navbar />
      <main className="content">
        <Outlet />
      </main>
    </div>
  )
}

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        {/* Player-facing shell with navbar */}
        <Route element={<PlayerShell />}>
          <Route path="/" element={<Navigate to="/challenges" replace />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route
            path="/challenges"
            element={
              <ProtectedRoute>
                <Challenges />
              </ProtectedRoute>
            }
          />
          <Route path="/scoreboard" element={<Scoreboard />} />
          <Route path="*" element={<NotFound />} />
        </Route>

        {/* Admin panel — full-page layout, no player navbar */}
        <Route
          path="/admin"
          element={
            <AdminRoute>
              <AdminLayout />
            </AdminRoute>
          }
        >
          <Route index element={<Navigate to="/admin/challenges" replace />} />
          <Route path="challenges" element={<AdminChallenges />} />
          <Route path="submissions" element={<AdminSubmissions />} />
          <Route path="users" element={<AdminUsers />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}

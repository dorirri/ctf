import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Navbar from './components/Navbar'
import ProtectedRoute from './components/ProtectedRoute'
import Login from './pages/Login'
import Register from './pages/Register'
import Challenges from './pages/Challenges'
import Scoreboard from './pages/Scoreboard'
import NotFound from './pages/NotFound'

export default function App() {
  return (
    <AuthProvider>
      <div className="page">
        <Navbar />
        <main className="content">
          <Routes>
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
          </Routes>
        </main>
      </div>
    </AuthProvider>
  )
}

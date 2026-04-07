import { createContext, useContext, useState, useCallback } from 'react'

const AuthContext = createContext(null)

function decodeJWT(token) {
  try {
    const payload = token.split('.')[1]
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return null
  }
}

function loadFromStorage() {
  const token = localStorage.getItem('ctf_token')
  if (!token) return { token: null, user: null }
  const payload = decodeJWT(token)
  if (!payload || payload.exp * 1000 < Date.now()) {
    localStorage.removeItem('ctf_token')
    localStorage.removeItem('ctf_user')
    return { token: null, user: null }
  }
  return {
    token,
    user: { id: payload.user_id, username: payload.username, role: payload.role },
  }
}

export function AuthProvider({ children }) {
  const [auth, setAuth] = useState(loadFromStorage)

  const login = useCallback((token) => {
    const payload = decodeJWT(token)
    if (!payload) return
    const user = { id: payload.user_id, username: payload.username, role: payload.role }
    localStorage.setItem('ctf_token', token)
    setAuth({ token, user })
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('ctf_token')
    localStorage.removeItem('ctf_user')
    setAuth({ token: null, user: null })
  }, [])

  return (
    <AuthContext.Provider value={{ ...auth, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}

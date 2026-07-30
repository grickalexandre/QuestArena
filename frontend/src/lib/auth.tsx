import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { signInWithEmailAndPassword, createUserWithEmailAndPassword } from 'firebase/auth'
import { api, type Teacher } from './api'
import { firebaseEnabled, getFirebaseAuth } from './firebase'

type AuthState = {
  teacher: Teacher | null
  token: string | null
  authMode: 'dev' | 'firebase' | null
  loading: boolean
  login: (email: string, password: string, name?: string) => Promise<void>
  register: (email: string, password: string, name?: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)
const TOKEN_KEY = 'qa_token'
const TEACHER_KEY = 'qa_teacher'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [teacher, setTeacher] = useState<Teacher | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [authMode, setAuthMode] = useState<'dev' | 'firebase' | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    ;(async () => {
      try {
        const mode = await api.authMode()
        setAuthMode(mode.mode as 'dev' | 'firebase')
        const saved = localStorage.getItem(TOKEN_KEY)
        if (saved) {
          const session = await api.session(saved)
          setToken(saved)
          setTeacher(session.teacher)
          localStorage.setItem(TEACHER_KEY, JSON.stringify(session.teacher))
        }
      } catch {
        localStorage.removeItem(TOKEN_KEY)
        localStorage.removeItem(TEACHER_KEY)
      } finally {
        setLoading(false)
      }
    })()
  }, [])

  const persist = (tok: string, t: Teacher) => {
    setToken(tok)
    setTeacher(t)
    localStorage.setItem(TOKEN_KEY, tok)
    localStorage.setItem(TEACHER_KEY, JSON.stringify(t))
  }

  const login = useCallback(async (email: string, password: string, name?: string) => {
    const mode = authMode ?? (await api.authMode()).mode
    if (mode === 'dev') {
      const res = await api.devLogin(email, password, name)
      persist(res.token, res.teacher)
      setAuthMode('dev')
      return
    }
    const auth = getFirebaseAuth()
    if (!auth || !firebaseEnabled) throw new Error('Firebase não configurado no frontend')
    const cred = await signInWithEmailAndPassword(auth, email, password)
    const idToken = await cred.user.getIdToken()
    const session = await api.session(idToken)
    persist(idToken, session.teacher)
    setAuthMode('firebase')
  }, [authMode])

  const register = useCallback(async (email: string, password: string, name?: string) => {
    const mode = authMode ?? (await api.authMode()).mode
    if (mode === 'dev') {
      const res = await api.devLogin(email, password, name || email.split('@')[0])
      persist(res.token, res.teacher)
      setAuthMode('dev')
      return
    }
    const auth = getFirebaseAuth()
    if (!auth || !firebaseEnabled) throw new Error('Firebase não configurado no frontend')
    const cred = await createUserWithEmailAndPassword(auth, email, password)
    const idToken = await cred.user.getIdToken()
    const session = await api.session(idToken)
    persist(idToken, { ...session.teacher, name: name || session.teacher.name })
    setAuthMode('firebase')
  }, [authMode])

  const logout = useCallback(() => {
    setToken(null)
    setTeacher(null)
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(TEACHER_KEY)
  }, [])

  const value = useMemo(
    () => ({ teacher, token, authMode, loading, login, register, logout }),
    [teacher, token, authMode, loading, login, register, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth outside provider')
  return ctx
}

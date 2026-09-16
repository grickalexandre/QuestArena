import { useState, type FormEvent } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import Credits from '../../components/Credits'
import { useAuth } from '../../lib/auth'

export default function TeacherLoginPage() {
  const { login, teacher, loading, authMode } = useAuth()
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  if (!loading && teacher) return <Navigate to="/teacher" replace />

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setBusy(true)
    try {
      await login(email, password)
      nav('/teacher')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha na autenticação')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="auth-shell">
      <form className="auth-card" onSubmit={onSubmit}>
        <Link to="/" className="brand-link">
          QuestArena
        </Link>
        <h1>Entrar como professor</h1>
        <p className="muted">
          Modo: <strong>{authMode ?? '...'}</strong>
          {authMode === 'dev' ? ' — login local sem Firebase' : ''}
        </p>
        <label>
          E-mail
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="oliveiraalexandre1972@gmail.com"
          />
        </label>
        <label>
          Senha
          <input
            type="password"
            required
            minLength={6}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
          />
        </label>
        {error && <p className="error">{error}</p>}
        <button className="btn btn-primary" disabled={busy} type="submit">
          {busy ? 'Aguarde...' : 'Entrar'}
        </button>
        <Credits compact />
      </form>
    </div>
  )
}

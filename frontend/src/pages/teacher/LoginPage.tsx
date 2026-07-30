import { useState, type FormEvent } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { useAuth } from '../../lib/auth'

export default function TeacherLoginPage() {
  const { login, register, teacher, loading, authMode } = useAuth()
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  if (!loading && teacher) return <Navigate to="/teacher" replace />

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setBusy(true)
    try {
      if (mode === 'login') await login(email, password)
      else await register(email, password, name)
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
        <h1>{mode === 'login' ? 'Entrar como professor' : 'Criar conta'}</h1>
        <p className="muted">
          Modo: <strong>{authMode ?? '...'}</strong>
          {authMode === 'dev' ? ' — login local sem Firebase' : ''}
        </p>
        {mode === 'register' && (
          <label>
            Nome
            <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Seu nome" />
          </label>
        )}
        <label>
          E-mail
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="prof@escola.com"
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
          {busy ? 'Aguarde...' : mode === 'login' ? 'Entrar' : 'Registrar'}
        </button>
        <button
          type="button"
          className="linkish"
          onClick={() => setMode(mode === 'login' ? 'register' : 'login')}
        >
          {mode === 'login' ? 'Criar nova conta' : 'Já tenho conta'}
        </button>
      </form>
    </div>
  )
}

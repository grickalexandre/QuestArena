import { useEffect, useState } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { api, type Quiz } from '../../lib/api'
import { useAuth } from '../../lib/auth'

export default function QuizzesPage() {
  const { teacher, token, logout, loading } = useAuth()
  const nav = useNavigate()
  const [quizzes, setQuizzes] = useState<Quiz[]>([])
  const [title, setTitle] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    if (!token) return
    api.listQuizzes(token).then(setQuizzes).catch((e) => setError(e.message))
  }, [token])

  if (!loading && !teacher) return <Navigate to="/teacher/login" replace />

  async function createQuiz() {
    if (!token || !title.trim()) return
    setBusy(true)
    setError('')
    try {
      const q = await api.createQuiz(token, title.trim(), '')
      nav(`/teacher/quizzes/${q.id}`)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Erro')
    } finally {
      setBusy(false)
    }
  }

  async function startLive(quizId: string) {
    if (!token) return
    setBusy(true)
    try {
      const session = await api.createSession(token, quizId)
      nav(`/teacher/host/${session.pin}`, { state: session })
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Erro ao iniciar')
    } finally {
      setBusy(false)
    }
  }

  async function removeQuiz(id: string) {
    if (!token || !confirm('Excluir este quiz?')) return
    await api.deleteQuiz(token, id)
    setQuizzes((prev) => prev.filter((q) => q.id !== id))
  }

  return (
    <div className="teacher-shell">
      <header className="teacher-top">
        <Link to="/" className="brand-link">
          QuestArena
        </Link>
        <div className="top-actions">
          <span className="pill">{teacher?.name}</span>
          <button className="btn btn-ghost" onClick={logout}>
            Sair
          </button>
        </div>
      </header>

      <section className="panel">
        <h1>Seus quizzes</h1>
        <div className="create-row">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Nome do novo quiz"
            onKeyDown={(e) => e.key === 'Enter' && createQuiz()}
          />
          <button className="btn btn-primary" disabled={busy} onClick={createQuiz}>
            Criar
          </button>
        </div>
        {error && <p className="error">{error}</p>}
      </section>

      <section className="quiz-list">
        {quizzes.length === 0 && <p className="muted">Nenhum quiz ainda. Crie o primeiro!</p>}
        {quizzes.map((q) => (
          <article key={q.id} className="quiz-row">
            <div>
              <h2>{q.title}</h2>
              <p className="muted">{q.description || 'Sem descrição'}</p>
            </div>
            <div className="row-actions">
              <Link className="btn btn-ghost" to={`/teacher/quizzes/${q.id}`}>
                Editar
              </Link>
              <button className="btn btn-accent" disabled={busy} onClick={() => startLive(q.id)}>
                Ao vivo
              </button>
              <button className="btn btn-danger" onClick={() => removeQuiz(q.id)}>
                Excluir
              </button>
            </div>
          </article>
        ))}
      </section>
    </div>
  )
}

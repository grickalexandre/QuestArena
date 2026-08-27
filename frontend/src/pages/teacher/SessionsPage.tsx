import { useEffect, useState } from 'react'
import { Link, Navigate } from 'react-router-dom'
import { api, type SessionRecord } from '../../lib/api'
import { useAuth } from '../../lib/auth'
import { downloadGradesCsv, formatGrade } from '../../lib/gradesCsv'

export default function SessionsPage() {
  const { teacher, token, logout, loading } = useAuth()
  const [sessions, setSessions] = useState<SessionRecord[]>([])
  const [openId, setOpenId] = useState<string | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!token) return
    api
      .listSessions(token)
      .then((list) => setSessions(Array.isArray(list) ? list : []))
      .catch((e) => setError(e instanceof Error ? e.message : 'Erro ao carregar notas'))
  }, [token])

  if (!loading && !teacher) return <Navigate to="/teacher/login" replace />

  const open = sessions.find((s) => s.id === openId)

  return (
    <div className="teacher-shell">
      <header className="teacher-top">
        <Link to="/teacher" className="brand-link">
          QuestArena
        </Link>
        <div className="top-actions">
          <Link className="btn btn-ghost" to="/teacher">
            Quizzes
          </Link>
          <span className="pill">{teacher?.name}</span>
          <button className="btn btn-ghost" onClick={logout}>
            Sair
          </button>
        </div>
      </header>

      <section className="panel">
        <h1>Notas das partidas</h1>
        <p className="muted">
          As notas ficam salvas depois que a partida encerra. Abra uma sessão ou baixe o CSV (abre no Excel) com RA,
          acertos e nota 0–10.
        </p>
        {error && <p className="error">{error}</p>}
      </section>

      <section className="quiz-list">
        {sessions.length === 0 && !error && (
          <p className="muted">Nenhuma partida encerrada ainda. Inicie um quiz ao vivo e termine para ver as notas.</p>
        )}
        {sessions.map((s) => (
          <article key={s.id} className="quiz-row">
            <div>
              <h2>{s.quizTitle || 'Quiz'}</h2>
              <p className="muted">
                PIN {s.pin} · {formatWhen(s.finishedAt)} · {s.ranking?.length || 0} aluno(s)
              </p>
            </div>
            <div className="row-actions">
              <button
                className="btn btn-ghost"
                type="button"
                disabled={!s.ranking?.length}
                onClick={() => downloadGradesCsv(s.quizTitle, s.pin, s.ranking || [])}
              >
                Baixar CSV
              </button>
              <button
                className="btn btn-accent"
                type="button"
                onClick={() => setOpenId((cur) => (cur === s.id ? null : s.id))}
              >
                {openId === s.id ? 'Fechar' : 'Ver notas'}
              </button>
            </div>
          </article>
        ))}
      </section>

      {open && (
        <section className="grades-card session-grades">
          <h3>{open.quizTitle} · PIN {open.pin}</h3>
          <p className="muted tiny">
            Nota 0–10 = 70% acertos + 30% XP. Encerrada em {formatWhen(open.finishedAt)}.
          </p>
          {(open.ranking || []).length === 0 ? (
            <p className="muted">Nenhum aluno nesta partida.</p>
          ) : (
            <div className="grades-table-wrap">
              <table className="grades-table">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>RA</th>
                    <th>Apelido</th>
                    <th>Acertos</th>
                    <th>XP</th>
                    <th>Nota</th>
                  </tr>
                </thead>
                <tbody>
                  {open.ranking.map((g) => (
                    <tr key={g.playerId}>
                      <td>{g.rank}</td>
                      <td>{g.ra || '—'}</td>
                      <td>{g.nickname}</td>
                      <td>
                        {g.correctCount}/{g.total}
                      </td>
                      <td>
                        {g.score}/{g.maxScore}
                      </td>
                      <td>
                        <strong>{formatGrade(g.grade)}</strong>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {(open.ranking || []).length > 0 && (
            <button
              className="btn btn-accent"
              type="button"
              onClick={() => downloadGradesCsv(open.quizTitle, open.pin, open.ranking)}
            >
              Baixar CSV das notas
            </button>
          )}
        </section>
      )}
    </div>
  )
}

function formatWhen(iso?: string) {
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleString('pt-BR', { dateStyle: 'short', timeStyle: 'short' })
}

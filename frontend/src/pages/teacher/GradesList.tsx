import { useState } from 'react'
import { downloadGradesCsv, formatGrade } from '../../lib/gradesCsv'
import type { SessionRecord } from '../../lib/api'

export function GradesList({
  sessions,
  emptyText,
}: {
  sessions: SessionRecord[]
  emptyText: string
}) {
  const [openId, setOpenId] = useState<string | null>(null)
  const open = sessions.find((s) => s.id === openId)

  if (sessions.length === 0) {
    return <p className="muted">{emptyText}</p>
  }

  return (
    <>
      <section className="quiz-list">
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
          <h3>
            {open.quizTitle} · PIN {open.pin}
          </h3>
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
    </>
  )
}

function formatWhen(iso?: string) {
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleString('pt-BR', { dateStyle: 'short', timeStyle: 'short' })
}

import { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'
import { api, type SessionRecord } from '../../lib/api'
import { useAuth } from '../../lib/auth'
import { GradesList } from './GradesList'
import { TeacherNav } from './TeacherNav'

export default function SessionsPage() {
  const { teacher, token, logout, loading } = useAuth()
  const [sessions, setSessions] = useState<SessionRecord[]>([])
  const [error, setError] = useState('')

  useEffect(() => {
    if (!token) return
    api
      .listSessions(token)
      .then((list) => setSessions(Array.isArray(list) ? list : []))
      .catch((e) => setError(e instanceof Error ? e.message : 'Erro ao carregar notas'))
  }, [token])

  if (!loading && !teacher) return <Navigate to="/teacher/login" replace />

  return (
    <div className="teacher-shell">
      <TeacherNav name={teacher?.name} current="grades" onLogout={logout} />

      <section className="panel">
        <h1>Notas da turma</h1>
        <p className="muted">
          Histórico de todas as partidas. Abra uma para ver RA, acertos e nota, ou baixe o CSV para o Excel.
        </p>
        {error && <p className="error">{error}</p>}
      </section>

      <GradesList
        sessions={sessions}
        emptyText="Nenhuma partida encerrada ainda. Inicie um quiz ao vivo e termine para as notas ficarem salvas aqui."
      />
    </div>
  )
}

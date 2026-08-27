import { Link } from 'react-router-dom'

export function TeacherNav({
  name,
  current,
  onLogout,
}: {
  name?: string
  current: 'quizzes' | 'grades'
  onLogout: () => void
}) {
  return (
    <header className="teacher-top">
      <Link to="/" className="brand-link">
        QuestArena
      </Link>
      <nav className="top-actions teacher-nav" aria-label="Área do professor">
        <Link className={current === 'quizzes' ? 'btn btn-primary' : 'btn btn-ghost'} to="/teacher">
          Quizzes
        </Link>
        <Link className={current === 'grades' ? 'btn btn-primary' : 'btn btn-ghost'} to="/teacher/sessions">
          Notas da turma
        </Link>
        <span className="pill">{name}</span>
        <button className="btn btn-ghost" type="button" onClick={onLogout}>
          Sair
        </button>
      </nav>
    </header>
  )
}

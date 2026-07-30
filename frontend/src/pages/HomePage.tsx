import { Link } from 'react-router-dom'

export default function HomePage() {
  return (
    <div className="home-shell">
      <div className="home-hero">
        <p className="brand">QuestArena</p>
        <h1>Arena de quizzes ao vivo</h1>
        <p className="lede">
          Professores criam desafios com peso e tempo. Alunos entram com PIN e competem por XP em tempo real.
        </p>
        <div className="cta-row">
          <Link className="btn btn-primary" to="/teacher/login">
            Área do professor
          </Link>
          <Link className="btn btn-accent" to="/play">
            Entrar na partida
          </Link>
        </div>
      </div>
      <div className="home-panel" aria-hidden="true">
        <div className="orb orb-a" />
        <div className="orb orb-b" />
        <div className="podium-preview">
          <span>2</span>
          <span className="gold">1</span>
          <span>3</span>
        </div>
      </div>
    </div>
  )
}

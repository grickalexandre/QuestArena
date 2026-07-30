import { useEffect, useMemo, useRef, useState } from 'react'
import { Link, Navigate, useParams } from 'react-router-dom'
import { avatarEmoji } from '../../lib/avatars'
import { useAuth } from '../../lib/auth'
import { useGameSocket } from '../../lib/useGameSocket'

type Player = { id: string; nickname: string; avatar: number; score: number }
type PublicQuestion = {
  id: string
  type?: 'multiple_choice' | 'essay'
  text: string
  options?: string[]
  weight: number
  timeLimitSec: number
  index: number
  total: number
}
type BoardEntry = { rank: number; playerId: string; nickname: string; avatar: number; score: number }
type ResultEntry = {
  playerId: string
  nickname: string
  text?: string
  similarity?: number
  points: number
  correct: boolean
}

const COLORS = ['#e21b3c', '#1368ce', '#d89e00', '#26890c']

export default function HostPage() {
  const { pin } = useParams()
  const { teacher, loading } = useAuth()
  const { connect, send, on, onOpen, connected } = useGameSocket()
  const [phase, setPhase] = useState('lobby')
  const [players, setPlayers] = useState<Player[]>([])
  const [question, setQuestion] = useState<PublicQuestion | null>(null)
  const [endsAt, setEndsAt] = useState<number | null>(null)
  const [correctIndex, setCorrectIndex] = useState<number | null>(null)
  const [expectedAnswer, setExpectedAnswer] = useState('')
  const [results, setResults] = useState<ResultEntry[]>([])
  const [board, setBoard] = useState<BoardEntry[]>([])
  const [answered, setAnswered] = useState(0)
  const [error, setError] = useState('')
  const [quizTitle, setQuizTitle] = useState('')
  const [hosted, setHosted] = useState(false)
  const [autoNextIn, setAutoNextIn] = useState<number | null>(null)
  const pendingAction = useRef<'start' | 'next' | null>(null)

  const remaining = useCountdown(endsAt)
  const teacherId = teacher?.id

  useEffect(() => {
    if (autoNextIn == null || autoNextIn <= 0) return
    const id = setTimeout(() => setAutoNextIn((v) => (v == null ? null : v - 1)), 1000)
    return () => clearTimeout(id)
  }, [autoNextIn])

  useEffect(() => {
    if (!teacherId || !pin) return

    connect()

    const offs = [
      on('hosted', (data) => {
        const d = data as {
          phase: string
          players: Player[]
          quizTitle: string
          question?: PublicQuestion
          endsAt?: string
          questionResult?: {
            correctIndex: number
            expectedAnswer?: string
            leaderboard: BoardEntry[]
            results: ResultEntry[]
            autoNextInSec?: number
          }
        }
        setHosted(true)
        setError('')
        setPhase(d.phase)
        setPlayers(d.players || [])
        setQuizTitle(d.quizTitle)
        if (d.question) setQuestion(d.question)
        if (d.endsAt && d.phase === 'question') setEndsAt(new Date(d.endsAt).getTime())
        if (d.questionResult) {
          setCorrectIndex(d.questionResult.correctIndex)
          setExpectedAnswer(d.questionResult.expectedAnswer || '')
          setResults(d.questionResult.results || [])
          setBoard(d.questionResult.leaderboard)
          setEndsAt(null)
          setAutoNextIn(d.questionResult.autoNextInSec ?? 5)
        }
        const action = pendingAction.current
        pendingAction.current = null
        if (action) send(action)
      }),
      on('player_joined', (data) => setPlayers((data as { players: Player[] }).players)),
      on('player_left', (data) => setPlayers((data as { players: Player[] }).players)),
      on('question', (data) => {
        const d = data as { question: PublicQuestion; endsAt: string }
        setPhase('question')
        setQuestion(d.question)
        setEndsAt(new Date(d.endsAt).getTime())
        setCorrectIndex(null)
        setExpectedAnswer('')
        setResults([])
        setAnswered(0)
        setAutoNextIn(null)
        setError('')
      }),
      on('answer_count', (data) => setAnswered((data as { answered: number }).answered)),
      on('question_result', (data) => {
        const d = data as {
          correctIndex: number
          expectedAnswer?: string
          leaderboard: BoardEntry[]
          results: ResultEntry[]
          autoNextInSec?: number
        }
        setPhase('reveal')
        setCorrectIndex(d.correctIndex)
        setExpectedAnswer(d.expectedAnswer || '')
        setResults(d.results || [])
        setBoard(d.leaderboard)
        setEndsAt(null)
        setAutoNextIn(d.autoNextInSec ?? 5)
      }),
      on('finished', (data) => {
        const d = data as { leaderboard: BoardEntry[] }
        setPhase('finished')
        setBoard(d.leaderboard)
        setAutoNextIn(null)
      }),
      on('error', (data) => {
        const msg = (data as { message: string }).message
        setError(msg)
        if (/host/i.test(msg)) {
          setHosted(false)
          if (!pendingAction.current) pendingAction.current = 'start'
          send('host_join', { pin, teacherId })
        }
      }),
      onOpen(() => {
        setHosted(false)
        send('host_join', { pin, teacherId })
      }),
    ]

    return () => offs.forEach((off) => off())
  }, [teacherId, pin, connect, on, onOpen, send])

  function runHostAction(action: 'start' | 'next') {
    setError('')
    if (!hosted || !connected) {
      pendingAction.current = action
      send('host_join', { pin, teacherId })
      return
    }
    pendingAction.current = null
    send(action)
  }

  if (!loading && !teacher) return <Navigate to="/teacher/login" replace />

  return (
    <div className="host-shell">
      <header className="host-top">
        <div>
          <Link to="/teacher" className="brand-link">
            QuestArena
          </Link>
          <h1>{quizTitle || 'Sessão ao vivo'}</h1>
        </div>
        <div className="pin-box">
          <span>PIN</span>
          <strong>{pin}</strong>
        </div>
      </header>

      {error && <p className="error banner">{error}</p>}
      {!hosted && phase === 'lobby' && (
        <p className="muted banner">Conectando como host...</p>
      )}

      {phase === 'lobby' && (
        <section className="host-lobby">
          <p className="lede">Alunos entram em /play com este PIN</p>
          <div className="player-grid">
            {players.map((p) => (
              <div key={p.id} className="player-chip">
                <span className="avatar emoji" aria-hidden="true">
                  {avatarEmoji(p.avatar)}
                </span>
                {p.nickname}
              </div>
            ))}
          </div>
          <button
            className="btn btn-primary btn-xl"
            onClick={() => runHostAction('start')}
            disabled={players.length === 0 || !hosted}
          >
            Iniciar partida ({players.length})
          </button>
        </section>
      )}

      {(phase === 'question' || phase === 'reveal') && question && (
        <section className="host-question">
          <div className="q-meta">
            <span>
              Questão {question.index + 1}/{question.total}
            </span>
            <span>{question.type === 'essay' ? 'Dissertativa' : 'Objetiva'}</span>
            <span>Peso {question.weight}</span>
            {phase === 'question' ? (
              <span className="timer">{formatClock(remaining)} · {answered}/{players.length}</span>
            ) : (
              <span>Resultado</span>
            )}
          </div>
          <h2>{question.text}</h2>

          {question.type === 'essay' ? (
            <div className="host-essay">
              {phase === 'reveal' && expectedAnswer && (
                <div className="expected-box">
                  <span className="label">Resposta esperada</span>
                  <p>{expectedAnswer}</p>
                </div>
              )}
              {phase === 'reveal' && (
                <ul className="essay-results">
                  {results.map((r) => (
                    <li key={r.playerId}>
                      <strong>{r.nickname}</strong>
                      <span className={r.correct ? 'ok' : 'ko'}>
                        {Math.round((r.similarity || 0) * 100)}% · +{r.points} XP
                      </span>
                      <p>{r.text || '—'}</p>
                    </li>
                  ))}
                </ul>
              )}
              {phase === 'question' && (
                <p className="muted">Aguardando respostas dissertativas...</p>
              )}
            </div>
          ) : (
            <div className="options-grid host host-opts">
              {(question.options || []).map((opt, i) => (
                <div
                  key={i}
                  className={`opt ${phase === 'reveal' && correctIndex === i ? 'correct' : ''} ${
                    phase === 'reveal' && correctIndex !== i ? 'dim' : ''
                  }`}
                  style={{ ['--opt' as string]: COLORS[i % COLORS.length] }}
                >
                  <span className="opt-letter">{String.fromCharCode(65 + i)}</span>
                  <span className="opt-text">{opt}</span>
                </div>
              ))}
            </div>
          )}

          {phase === 'reveal' && (
            <>
              <Leaderboard board={board} />
              <p className="muted">
                {autoNextIn != null && autoNextIn > 0
                  ? `Próxima questão em ${autoNextIn}s...`
                  : 'Avançando...'}
              </p>
              <button className="btn btn-primary btn-xl" onClick={() => runHostAction('next')}>
                {question.index + 1 >= question.total ? 'Ver pódio agora' : 'Próxima agora'}
              </button>
            </>
          )}
        </section>
      )}

      {phase === 'finished' && (
        <section className="host-question">
          <h2>Partida encerrada</h2>
          <Podium board={board} />
          <Leaderboard board={board} />
          <Link className="btn btn-primary" to="/teacher">
            Voltar
          </Link>
        </section>
      )}
    </div>
  )
}

export function Leaderboard({ board }: { board: BoardEntry[] }) {
  return (
    <ol className="leaderboard">
      {board.map((e) => (
        <li key={e.playerId}>
          <span className="rank">#{e.rank}</span>
          <span className="nick">
            <span className="avatar emoji sm" aria-hidden="true">
              {avatarEmoji(e.avatar)}
            </span>
            {e.nickname}
          </span>
          <span className="xp">{e.score} XP</span>
        </li>
      ))}
    </ol>
  )
}

export function Podium({ board }: { board: BoardEntry[] }) {
  const top = useMemo(() => {
    const [first, second, third] = board
    return [
      { entry: second, place: 2 },
      { entry: first, place: 1 },
      { entry: third, place: 3 },
    ].filter((x) => x.entry)
  }, [board])

  return (
    <div className="podium">
      {top.map(({ entry, place }) => (
        <div key={place} className={`pod place-${place}`}>
          <span className="avatar-big sm" aria-hidden="true">
            {avatarEmoji(entry!.avatar)}
          </span>
          <strong>{entry!.nickname}</strong>
          <span>{entry!.score} XP</span>
          <div className="block">{place}</div>
        </div>
      ))}
    </div>
  )
}

function useCountdown(endsAt: number | null) {
  const [left, setLeft] = useState(0)
  useEffect(() => {
    if (!endsAt) {
      setLeft(0)
      return
    }
    const tick = () => setLeft(Math.max(0, Math.ceil((endsAt - Date.now()) / 1000)))
    tick()
    const id = setInterval(tick, 200)
    return () => clearInterval(id)
  }, [endsAt])
  return left
}

function formatClock(totalSec: number) {
  const m = Math.floor(totalSec / 60)
  const s = totalSec % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

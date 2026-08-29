import { useEffect, useMemo, useRef, useState } from 'react'
import { Link, Navigate, useParams } from 'react-router-dom'
import Credits from '../../components/Credits'
import CodeBlock from '../../components/CodeBlock'
import { avatarEmoji } from '../../lib/avatars'
import { api } from '../../lib/api'
import { useAuth } from '../../lib/auth'
import { downloadGradesCsv, formatGrade } from '../../lib/gradesCsv'
import { useGameSocket } from '../../lib/useGameSocket'

type Player = {
  id: string
  nickname: string
  ra?: string
  avatar: number
  score: number
  connected?: boolean
  hidden?: boolean
  awayCount?: number
  awayTotal?: number
}
type PublicQuestion = {
  id: string
  type?: 'multiple_choice' | 'essay'
  text: string
  options?: string[]
  codeSnippet?: string
  codeLanguage?: string
  weight: number
  timeLimitSec: number
  index: number
  total: number
}
type BoardEntry = {
  rank: number
  playerId: string
  nickname: string
  ra?: string
  avatar: number
  score: number
  correctCount?: number
  grade?: number
}
type GradeRow = {
  rank: number
  playerId: string
  ra: string
  nickname: string
  correctCount: number
  total: number
  score: number
  maxScore: number
  grade: number
}
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
  const { connect, send, on, onOpen, connected, disconnect } = useGameSocket()
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
  const [grades, setGrades] = useState<GradeRow[]>([])
  const [roomMissing, setRoomMissing] = useState(false)
  const [playUrl, setPlayUrl] = useState(`${window.location.origin}/play`)
  const pendingAction = useRef<'start' | 'next' | 'end' | null>(null)
  const roomMissingRef = useRef(false)

  const remaining = useCountdown(endsAt)
  const teacherId = teacher?.id

  useEffect(() => {
    api
      .health()
      .then((h) => setPlayUrl(buildPlayUrl(h.lanIP)))
      .catch(() => setPlayUrl(buildPlayUrl()))
  }, [])

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
          leaderboard?: BoardEntry[]
          grades?: GradeRow[]
        }
        setHosted(true)
        setRoomMissing(false)
        roomMissingRef.current = false
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
        if (d.leaderboard) setBoard(d.leaderboard)
        if (d.grades) setGrades(d.grades)
        const action = pendingAction.current
        pendingAction.current = null
        if (action) send(action)
      }),
      on('player_joined', (data) => setPlayers((data as { players: Player[] }).players)),
      on('player_left', (data) => setPlayers((data as { players: Player[] }).players)),
      on('player_presence', (data) => setPlayers((data as { players: Player[] }).players)),
      on('question', (data) => {
        const d = data as { question: PublicQuestion; endsAt: string; players?: Player[] }
        setPhase('question')
        setQuestion(d.question)
        setEndsAt(new Date(d.endsAt).getTime())
        setCorrectIndex(null)
        setExpectedAnswer('')
        setResults([])
        setAnswered(0)
        setAutoNextIn(null)
        setError('')
        if (d.players) setPlayers(d.players)
        else setPlayers((prev) => prev.map((p) => ({ ...p, awayCount: 0 })))
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
        const d = data as { leaderboard: BoardEntry[]; grades?: GradeRow[] }
        setPhase('finished')
        setBoard(d.leaderboard)
        setGrades(d.grades || [])
        setAutoNextIn(null)
      }),
      on('error', (data) => {
        const msg = (data as { message: string }).message
        setError(msg)
        if (/sala não encontrada|room not found/i.test(msg)) {
          roomMissingRef.current = true
          setRoomMissing(true)
          setHosted(false)
          disconnect()
          return
        }
        if (/host/i.test(msg)) {
          setHosted(false)
          send('host_join', { pin, teacherId })
        }
      }),
      onOpen(() => {
        if (roomMissingRef.current) return
        setHosted(false)
        send('host_join', { pin, teacherId })
      }),
    ]

    return () => offs.forEach((off) => off())
  }, [teacherId, pin, connect, disconnect, on, onOpen, send])

  function runHostAction(action: 'start' | 'next' | 'end') {
    if (roomMissingRef.current) return
    setError('')
    if (!hosted || !connected) {
      pendingAction.current = action
      send('host_join', { pin, teacherId })
      return
    }
    pendingAction.current = null
    send(action)
  }

  const connectedCount = players.filter((p) => p.connected !== false).length
  const awayNow = players.filter((p) => p.hidden || p.connected === false)
  const flagged = players.filter((p) => (p.awayCount || 0) > 0)

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

      {roomMissing && (
        <section className="host-lobby">
          <p className="error banner">
            Esta sala (PIN {pin}) não existe mais. Isso acontece se o servidor reiniciar. Volte e clique em
            <strong> Ao vivo</strong> para gerar um PIN novo.
          </p>
          <Link className="btn btn-primary btn-xl" to="/teacher">
            Gerar novo PIN
          </Link>
        </section>
      )}

      {error && !roomMissing && <p className="error banner">{error}</p>}
      {!hosted && !roomMissing && phase === 'lobby' && (
        <p className="muted banner">Conectando como host...</p>
      )}

      {phase === 'lobby' && !roomMissing && (
        <section className="host-lobby">
          <p className="lede">Peça aos alunos para abrir este link no celular:</p>
          <div className="join-url-box">
            <strong>{playUrl}</strong>
            <button
              className="btn btn-ghost"
              type="button"
              onClick={() => navigator.clipboard.writeText(playUrl).catch(() => {})}
            >
              Copiar link
            </button>
          </div>
          <p className="muted">Depois entram com o PIN <strong>{pin}</strong>, RA e apelido.</p>
          <div className="player-grid">
            {players.map((p) => (
              <div key={p.id} className={playerChipClass(p)}>
                <span className="avatar emoji" aria-hidden="true">
                  {avatarEmoji(p.avatar)}
                </span>
                {p.nickname}
                {p.ra ? <span className="muted tiny"> · RA {p.ra}</span> : null}
                {playerChipNote(p)}
              </div>
            ))}
          </div>
          <button
            className="btn btn-primary btn-xl"
            onClick={() => runHostAction('start')}
            disabled={connectedCount === 0 || !hosted}
          >
            Iniciar partida ({connectedCount})
          </button>
          {hosted && connectedCount === 0 && (
            <p className="muted tiny">Quando o aluno entrar, o nome aparece aqui. Se o celular não abrir o link, use a rede Wi‑Fi da sala.</p>
          )}
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
              <span className="timer">{formatClock(remaining)} · {answered}/{connectedCount || players.length}</span>
            ) : (
              <span>Resultado</span>
            )}
          </div>
          <h2>{question.text}</h2>

          {question.codeSnippet && (
            <CodeBlock code={question.codeSnippet} language={question.codeLanguage} />
          )}

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

          {phase === 'question' && (
            <>
              <PresencePanel players={players} awayNow={awayNow} flagged={flagged} />
              <button className="btn btn-ghost" type="button" onClick={() => runHostAction('end')}>
                Encerrar e ver notas
              </button>
            </>
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
              <button className="btn btn-ghost" type="button" onClick={() => runHostAction('end')}>
                Encerrar e ver notas
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
          <div className="grades-card">
            <h3>Notas da turma</h3>
            <p className="muted tiny">
              Nota 0–10 = 70% acertos + 30% XP (quem responde mais rápido sobe a nota). Baixe o CSV para somar com
              outros quizzes pelo RA.
            </p>
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
                  {(grades.length ? grades : board).map((g) => {
                    const row = g as GradeRow & BoardEntry
                    return (
                      <tr key={row.playerId}>
                        <td>{row.rank}</td>
                        <td>{row.ra || '—'}</td>
                        <td>{row.nickname}</td>
                        <td>
                          {row.correctCount ?? '—'}
                          {row.total != null ? `/${row.total}` : ''}
                        </td>
                        <td>
                          {row.score}
                          {row.maxScore != null ? `/${row.maxScore}` : ''}
                        </td>
                        <td>
                          <strong>{formatGrade(row.grade)}</strong>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
            <button
              className="btn btn-accent"
              type="button"
              disabled={grades.length === 0 && board.length === 0}
              onClick={() => downloadGradesCsv(quizTitle, pin || '', grades.length ? grades : boardToGrades(board))}
            >
              Baixar CSV das notas
            </button>
          </div>
          <Link className="btn btn-primary" to="/teacher/sessions">
            Ver notas da turma
          </Link>
        </section>
      )}
      <Credits compact />
    </div>
  )
}

function playerChipClass(p: Player) {
  const parts = ['player-chip']
  if (p.connected === false) parts.push('offline')
  if (p.hidden || p.connected === false) parts.push('away')
  else if ((p.awayCount || 0) > 0) parts.push('flagged')
  return parts.join(' ')
}

function playerChipNote(p: Player) {
  if (p.connected === false) {
    return <span className="muted tiny"> · offline</span>
  }
  if (p.hidden) {
    return <span className="away-note"> · fora da tela</span>
  }
  if ((p.awayCount || 0) > 0) {
    return <span className="flag-note"> · saiu {p.awayCount}x</span>
  }
  return null
}

function PresencePanel({
  players,
  awayNow,
  flagged,
}: {
  players: Player[]
  awayNow: Player[]
  flagged: Player[]
}) {
  if (players.length === 0) return null
  return (
    <div className="presence-panel">
      {awayNow.length > 0 ? (
        <p className="away-banner">
          {awayNow.length === 1
            ? `${awayNow[0].nickname} saiu da tela agora`
            : `${awayNow.length} alunos saíram da tela agora`}
        </p>
      ) : flagged.length > 0 ? (
        <p className="flag-banner">
          {flagged.length === 1
            ? `${flagged[0].nickname} saiu da tela nesta questão`
            : `${flagged.length} alunos já saíram da tela nesta questão`}
        </p>
      ) : (
        <p className="muted tiny">Ninguém saiu da tela nesta questão.</p>
      )}
      <div className="player-grid compact">
        {players.map((p) => (
          <div key={p.id} className={playerChipClass(p)}>
            <span className="avatar emoji" aria-hidden="true">
              {avatarEmoji(p.avatar)}
            </span>
            {p.nickname}
            {playerChipNote(p)}
          </div>
        ))}
      </div>
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
            {e.ra ? <span className="muted tiny"> · {e.ra}</span> : null}
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
          {entry!.ra ? <span className="muted tiny">RA {entry!.ra}</span> : null}
          <span>{entry!.score} XP</span>
          {entry!.grade != null ? <span>Nota {formatGrade(entry!.grade)}</span> : null}
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

function buildPlayUrl(lanIP?: string) {
  const { protocol, hostname, port } = window.location
  const host =
    hostname === 'localhost' || hostname === '127.0.0.1' ? lanIP || hostname : hostname
  const suffix = port ? `:${port}` : ''
  return `${protocol}//${host}${suffix}/play`
}

function boardToGrades(board: BoardEntry[]): GradeRow[] {
  return board.map((e) => ({
    rank: e.rank,
    playerId: e.playerId,
    ra: e.ra || '',
    nickname: e.nickname,
    correctCount: e.correctCount ?? 0,
    total: 0,
    score: e.score,
    maxScore: 0,
    grade: e.grade ?? 0,
  }))
}

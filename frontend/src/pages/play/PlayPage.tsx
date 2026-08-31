import { useEffect, useMemo, useRef, useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import Credits from '../../components/Credits'
import CodeBlock from '../../components/CodeBlock'
import { AVATARS, avatarEmoji } from '../../lib/avatars'
import { useGameSocket } from '../../lib/useGameSocket'
import { usePlayPresence } from '../../lib/usePlayPresence'
import { Leaderboard, Podium } from '../teacher/HostPage'

type QuestionType = 'multiple_choice' | 'essay'

type PublicQuestion = {
  id: string
  type: QuestionType
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
  grade?: number
}

type ResultEntry = {
  playerId: string
  points: number
  correct: boolean
  similarity?: number
  text?: string
  choice?: number
}

type QuestionResult = {
  type?: QuestionType
  correctIndex: number
  expectedAnswer?: string
  expectedAnswers?: string[]
  leaderboard: BoardEntry[]
  results: ResultEntry[]
  autoNextInSec?: number
}

type JoinedPayload = {
  playerId: string
  nickname?: string
  ra?: string
  quizTitle: string
  avatar?: number
  phase?: 'lobby' | 'question' | 'reveal' | 'finished'
  answered?: boolean
  choice?: number
  answerText?: string
  question?: PublicQuestion
  endsAt?: string
  questionResult?: QuestionResult
  leaderboard?: BoardEntry[]
}

const COLORS = ['#e21b3c', '#1368ce', '#d89e00', '#26890c']
const PLAY_SESSION_KEY = 'qa_play_session'

type PlaySession = {
  pin: string
  nickname: string
  ra: string
  avatar: number
  playerId?: string
}

function loadPlaySession(): PlaySession | null {
  try {
    const raw = sessionStorage.getItem(PLAY_SESSION_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as PlaySession
    if (!parsed.pin || !parsed.ra || !parsed.nickname) return null
    return parsed
  } catch {
    return null
  }
}

function savePlaySession(s: PlaySession) {
  sessionStorage.setItem(PLAY_SESSION_KEY, JSON.stringify(s))
}

function clearPlaySession() {
  sessionStorage.removeItem(PLAY_SESSION_KEY)
}

export default function PlayPage() {
  const { connect, send, on, onOpen, onClose, connected } = useGameSocket()
  const saved = useMemo(() => loadPlaySession(), [])
  const [pin, setPin] = useState(saved?.pin ?? '')
  const [nickname, setNickname] = useState(saved?.nickname ?? '')
  const [ra, setRa] = useState(saved?.ra ?? '')
  const [avatar, setAvatar] = useState(() => saved?.avatar ?? Math.floor(Math.random() * AVATARS.length))
  const [phase, setPhase] = useState<'join' | 'lobby' | 'question' | 'reveal' | 'finished'>('join')
  const [playerId, setPlayerId] = useState(saved?.playerId ?? '')
  const [quizTitle, setQuizTitle] = useState('')
  const [question, setQuestion] = useState<PublicQuestion | null>(null)
  const [endsAt, setEndsAt] = useState<number | null>(null)
  const [selected, setSelected] = useState<number | null>(null)
  const [essayDraft, setEssayDraft] = useState('')
  const [submittedText, setSubmittedText] = useState('')
  const [locked, setLocked] = useState(false)
  const [correctIndex, setCorrectIndex] = useState<number | null>(null)
  const [expectedAnswer, setExpectedAnswer] = useState('')
  const [expectedAnswers, setExpectedAnswers] = useState<string[]>([])
  const [lastPoints, setLastPoints] = useState<number | null>(null)
  const [lastSimilarity, setLastSimilarity] = useState<number | null>(null)
  const [lastCorrect, setLastCorrect] = useState<boolean | null>(null)
  const [board, setBoard] = useState<BoardEntry[]>([])
  const [error, setError] = useState('')
  const [autoNextIn, setAutoNextIn] = useState<number | null>(null)
  const [reconnecting, setReconnecting] = useState(false)
  const [leftThisQuestion, setLeftThisQuestion] = useState(false)
  const away = usePlayPresence(send, phase !== 'join')
  const remaining = useCountdown(endsAt)
  const timerPct = useMemo(() => {
    if (!question || !endsAt) return 0
    const total = question.timeLimitSec * 1000
    const left = Math.max(0, endsAt - Date.now())
    return Math.min(100, (left / total) * 100)
  }, [question, endsAt, remaining])

  const credsRef = useRef({ pin, nickname, ra, avatar, playerId, phase })
  credsRef.current = { pin, nickname, ra, avatar, playerId, phase }

  useEffect(() => {
    if (autoNextIn == null || autoNextIn <= 0) return
    const id = setTimeout(() => setAutoNextIn((v) => (v == null ? null : v - 1)), 1000)
    return () => clearTimeout(id)
  }, [autoNextIn])

  useEffect(() => {
    if (phase === 'question' && away) setLeftThisQuestion(true)
  }, [phase, away])

  function applyQuestion(
    q: PublicQuestion,
    endsAtIso?: string,
    alreadyAnswered?: boolean,
    choice?: number,
    answerText?: string,
  ) {
    setQuestion({
      ...q,
      type: q.type || 'multiple_choice',
    })
    setEndsAt(endsAtIso ? new Date(endsAtIso).getTime() : null)
    setSelected(typeof choice === 'number' && choice >= 0 ? choice : null)
    setEssayDraft(answerText || '')
    setSubmittedText(answerText || '')
    setLocked(!!alreadyAnswered)
    setCorrectIndex(null)
    setExpectedAnswer('')
    setExpectedAnswers([])
    setLastPoints(null)
    setLastSimilarity(null)
    setLastCorrect(null)
    setAutoNextIn(null)
    setLeftThisQuestion(false)
  }

  function applyResult(d: QuestionResult, id: string) {
    setPhase('reveal')
    setCorrectIndex(d.correctIndex)
    setExpectedAnswer(d.expectedAnswer || '')
    setExpectedAnswers(d.expectedAnswers || [])
    setBoard(d.leaderboard)
    setEndsAt(null)
    setAutoNextIn(d.autoNextInSec ?? 5)
    const mine = d.results.find((r) => r.playerId === id)
    if (mine) {
      setLastPoints(mine.points)
      setLastCorrect(mine.correct)
      if (typeof mine.similarity === 'number') setLastSimilarity(mine.similarity)
      if (mine.text) setSubmittedText(mine.text)
      if (typeof mine.choice === 'number' && mine.choice >= 0) setSelected(mine.choice)
    }
  }

  function sendJoin(next = credsRef.current) {
    if (!next.pin.trim() || !next.ra.trim() || !next.nickname.trim()) return
    send('join', {
      pin: next.pin.trim(),
      nickname: next.nickname.trim(),
      ra: next.ra.trim(),
      avatar: next.avatar,
    })
  }

  useEffect(() => {
    const offs = [
      on('joined', (data) => {
        const d = data as JoinedPayload
        setPlayerId(d.playerId)
        setQuizTitle(d.quizTitle)
        if (d.nickname) setNickname(d.nickname)
        if (d.ra) setRa(d.ra)
        if (typeof d.avatar === 'number') setAvatar(d.avatar)
        setError('')
        setReconnecting(false)
        const creds = credsRef.current
        savePlaySession({
          pin: creds.pin.trim(),
          nickname: (d.nickname || creds.nickname).trim(),
          ra: (d.ra || creds.ra).trim(),
          avatar: typeof d.avatar === 'number' ? d.avatar : creds.avatar,
          playerId: d.playerId,
        })
        const nextPhase = d.phase || 'lobby'
        if (nextPhase === 'question' && d.question) {
          setPhase('question')
          applyQuestion(d.question, d.endsAt, d.answered, d.choice, d.answerText)
        } else if (nextPhase === 'reveal' && d.question) {
          applyQuestion(d.question, d.endsAt, true, d.choice)
          if (d.questionResult) applyResult(d.questionResult, d.playerId)
          else setPhase('reveal')
        } else if (nextPhase === 'finished') {
          setPhase('finished')
          if (d.leaderboard) setBoard(d.leaderboard)
        } else {
          setPhase('lobby')
        }
      }),
      on('question', (data) => {
        const d = data as { question: PublicQuestion; endsAt: string }
        setPhase('question')
        applyQuestion(d.question, d.endsAt, false)
        setError('')
        setReconnecting(false)
      }),
      on('answer_ack', (data) => {
        const d = data as { points: number; similarity?: number }
        setLocked(true)
        setLastPoints(d.points)
        if (typeof d.similarity === 'number') setLastSimilarity(d.similarity)
      }),
      on('question_result', (data) => {
        applyResult(data as QuestionResult, credsRef.current.playerId)
      }),
      on('finished', (data) => {
        setPhase('finished')
        setBoard((data as { leaderboard: BoardEntry[] }).leaderboard)
        setAutoNextIn(null)
        setReconnecting(false)
      }),
      on('error', (data) => {
        const msg = (data as { message: string }).message
        setError(msg)
        if (/PIN inválido/i.test(msg)) {
          clearPlaySession()
          setPhase('join')
          setReconnecting(false)
        }
      }),
      onOpen(() => {
        const creds = credsRef.current
        const stored = loadPlaySession()
        if (!stored) return
        setReconnecting(false)
        sendJoin({
          pin: stored.pin,
          nickname: stored.nickname,
          ra: stored.ra,
          avatar: stored.avatar,
          playerId: stored.playerId || creds.playerId,
          phase: creds.phase,
        })
      }),
      onClose(() => {
        if (credsRef.current.phase !== 'join') setReconnecting(true)
      }),
    ]
    const stored = loadPlaySession()
    if (stored) {
      connect()
    }
    return () => offs.forEach((off) => off())
  }, [on, onOpen, onClose, send, connect])

  function join(e: FormEvent) {
    e.preventDefault()
    setError('')
    savePlaySession({ pin: pin.trim(), nickname: nickname.trim(), ra: ra.trim(), avatar })
    connect()
    sendJoin({ pin, nickname, ra, avatar, playerId, phase })
  }

  function resetJoin() {
    clearPlaySession()
    setPhase('join')
    setPin('')
    setNickname('')
    setRa('')
    setPlayerId('')
    setQuizTitle('')
    setAvatar(Math.floor(Math.random() * AVATARS.length))
    setReconnecting(false)
    setError('')
  }

  function answerChoice(choice: number) {
    if (phase !== 'question') return
    if (locked && selected === choice) return
    setSelected(choice)
    setError('')
    send('answer', { choice })
  }

  function submitEssay(e?: FormEvent) {
    e?.preventDefault()
    if (phase !== 'question') return
    const text = essayDraft.trim()
    if (!text) {
      setError('Escreva sua resposta antes de enviar')
      return
    }
    if (locked && text === submittedText.trim()) return
    setError('')
    send('answer', { text })
    setSubmittedText(text)
  }

  const isEssay = question?.type === 'essay'
  const myRank = board.find((b) => b.playerId === playerId)

  return (
    <div className={`play-shell ${phase === 'question' || phase === 'reveal' ? 'play-focus' : ''}`}>
      {phase === 'join' && (
        <form className="join-card" onSubmit={join}>
          <Link to="/" className="brand-link">
            QuestArena
          </Link>
          <h1>Entrar na arena</h1>

          <div className="avatar-preview">
            <span className="avatar-big" aria-hidden="true">
              {avatarEmoji(avatar)}
            </span>
            <span className="muted tiny">Escolha seu avatar</span>
          </div>

          <div className="avatar-grid" role="listbox" aria-label="Avatares">
            {AVATARS.map((a) => (
              <button
                key={a.id}
                type="button"
                role="option"
                aria-selected={avatar === a.id}
                className={`avatar-pick ${avatar === a.id ? 'selected' : ''}`}
                title={a.label}
                onClick={() => setAvatar(a.id)}
              >
                <span aria-hidden="true">{a.emoji}</span>
              </button>
            ))}
          </div>

          <label>
            PIN da partida
            <input
              value={pin}
              onChange={(e) => setPin(e.target.value.replace(/\D/g, '').slice(0, 6))}
              inputMode="numeric"
              placeholder="000000"
              required
              minLength={6}
              autoComplete="one-time-code"
            />
          </label>
          <label>
            RA
            <input
              value={ra}
              onChange={(e) => setRa(e.target.value.replace(/\s/g, '').slice(0, 20))}
              placeholder="Seu RA"
              required
              minLength={3}
              autoComplete="off"
            />
          </label>
          <label>
            Nickname
            <input
              value={nickname}
              onChange={(e) => setNickname(e.target.value.slice(0, 20))}
              placeholder="Seu apelido"
              required
              minLength={2}
              autoComplete="nickname"
            />
          </label>
          {error && <p className="error">{error}</p>}
          <button className="btn btn-accent btn-xl" type="submit">
            Entrar como {avatarEmoji(avatar)} {nickname.trim() || '...'}
          </button>
          <Credits compact />
        </form>
      )}

      {reconnecting && phase !== 'join' && (
        <p className="muted banner">Reconectando à partida...</p>
      )}
      {phase !== 'join' && !connected && !reconnecting && (
        <p className="muted banner">Sem conexão — tentando de novo...</p>
      )}

      {phase === 'lobby' && (
        <section className="play-wait">
          <p className="brand">QuestArena</p>
          <div className="avatar-big lobby">{avatarEmoji(avatar)}</div>
          <h1>Bem-vindo, {nickname}</h1>
          <p className="lede">{quizTitle}</p>
          <p className="muted">RA {ra}</p>
          <div className="pulse-ring" />
          <p>Aguardando o professor iniciar...</p>
          {error && <p className="error">{error}</p>}
        </section>
      )}

      {(phase === 'question' || phase === 'reveal') && question && (
        <section className="play-question">
          <div className="timer-track" aria-hidden="true">
            <div
              className={`timer-fill ${remaining <= 5 && phase === 'question' ? 'urgent' : ''}`}
              style={{ width: `${phase === 'question' ? timerPct : 0}%` }}
            />
          </div>

          <div className="q-meta">
            <span className="pill soft">
              {question.index + 1}/{question.total}
            </span>
            <span className="pill soft">{isEssay ? 'Dissertativa' : 'Objetiva'}</span>
            {phase === 'question' ? (
              <span className={`timer ${remaining <= 5 ? 'urgent-text' : ''}`}>{formatClock(remaining)}</span>
            ) : (
              <span>Resultado</span>
            )}
          </div>

          {phase === 'question' && leftThisQuestion && (
            <p className="error banner away-warn">
              Você saiu da tela. O professor foi avisado.
            </p>
          )}

          <h2>{question.text}</h2>

          {question.codeSnippet && (
            <CodeBlock code={question.codeSnippet} language={question.codeLanguage} />
          )}

          {phase === 'question' && !isEssay && (
            <>
              <div className="options-grid play-opts">
                {(question.options || []).map((opt, i) => {
                  let cls = 'opt'
                  if (selected === i) cls += ' selected'
                  return (
                    <button
                      key={i}
                      className={cls}
                      style={{ ['--opt' as string]: COLORS[i % COLORS.length] }}
                      onClick={() => answerChoice(i)}
                    >
                      <span className="opt-letter">{String.fromCharCode(65 + i)}</span>
                      <span className="opt-text">{opt}</span>
                    </button>
                  )
                })}
              </div>
              {locked ? (
                <p className="muted tiny change-hint">Resposta salva. Clique em outra opção para trocar até o tempo acabar.</p>
              ) : (
                <p className="muted tiny change-hint">Pode trocar a escolha enquanto o tempo estiver rodando.</p>
              )}
            </>
          )}

          {phase === 'question' && isEssay && (
            <form className="essay-box" onSubmit={submitEssay}>
              <textarea
                value={essayDraft}
                onChange={(e) => setEssayDraft(e.target.value.slice(0, 2000))}
                placeholder="Explique com suas palavras. A correção compara com a resposta de referência."
                rows={6}
                autoFocus
              />
              <div className="essay-actions">
                <span className="muted tiny">{essayDraft.trim().length}/2000</span>
                <button
                  className="btn btn-accent btn-xl"
                  type="submit"
                  disabled={!essayDraft.trim() || (locked && essayDraft.trim() === submittedText.trim())}
                >
                  {locked ? 'Atualizar resposta' : 'Enviar resposta'}
                </button>
              </div>
              {locked && (
                <p className="muted tiny change-hint">Salva. Pode editar e atualizar até o tempo acabar.</p>
              )}
            </form>
          )}

          {phase === 'reveal' && (
            <div className="reveal-card">
              {isEssay ? (
                <>
                  <div className="sim-meter">
                    <div className="sim-bar">
                      <div
                        className="sim-fill"
                        style={{ width: `${Math.round((lastSimilarity || 0) * 100)}%` }}
                      />
                    </div>
                    <strong>{Math.round((lastSimilarity || 0) * 100)}% similar</strong>
                  </div>
                  <p className={`result-banner ${lastCorrect ? 'ok' : 'ko'}`}>
                    {lastCorrect ? 'Resposta aceita' : 'Não atingiu o mínimo'}
                    {lastPoints != null ? ` · +${lastPoints} XP` : ''}
                  </p>
                  <div className="answer-compare">
                    <div>
                      <span className="label">Sua resposta</span>
                      <p>{submittedText || essayDraft || '—'}</p>
                    </div>
                    <div>
                      <span className="label">Esperada</span>
                      <p>{expectedAnswer || '—'}</p>
                      {expectedAnswers.length > 0 && (
                        <ul className="alt-accepted">
                          {expectedAnswers.map((a) => (
                            <li key={a}>{a}</li>
                          ))}
                        </ul>
                      )}
                    </div>
                  </div>
                </>
              ) : (
                <>
                  <div className="options-grid play-opts">
                    {(question.options || []).map((opt, i) => {
                      let cls = 'opt'
                      if (selected === i) cls += ' selected'
                      if (correctIndex === i) cls += ' correct'
                      if (selected === i && correctIndex !== i) cls += ' wrong'
                      return (
                        <div
                          key={i}
                          className={cls}
                          style={{ ['--opt' as string]: COLORS[i % COLORS.length] }}
                        >
                          <span className="opt-letter">{String.fromCharCode(65 + i)}</span>
                          <span className="opt-text">{opt}</span>
                        </div>
                      )
                    })}
                  </div>
                  <p className={`result-banner ${lastPoints && lastPoints > 0 ? 'ok' : 'ko'}`}>
                    {lastPoints && lastPoints > 0 ? `+${lastPoints} XP` : 'Sem XP nesta rodada'}
                  </p>
                </>
              )}
              {autoNextIn != null && autoNextIn > 0 && (
                <p className="muted">Próxima em {autoNextIn}s...</p>
              )}
            </div>
          )}

          {error && <p className="error">{error}</p>}
        </section>
      )}

      {phase === 'finished' && (
        <section className="play-question">
          <h2>Fim da quest!</h2>
          <Podium board={board} />
          <Leaderboard board={board} />
          <p className="muted">
            Sua posição: {myRank ? `#${myRank.rank}` : '—'}
            {myRank?.grade != null ? ` · Nota ${myRank.grade.toFixed(1).replace('.', ',')}` : ''}
          </p>
          <button
            className="btn btn-primary"
            onClick={resetJoin}
          >
            Nova partida
          </button>
        </section>
      )}
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

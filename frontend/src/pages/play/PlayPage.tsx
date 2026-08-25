import { useEffect, useMemo, useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { AVATARS, avatarEmoji } from '../../lib/avatars'
import { useGameSocket } from '../../lib/useGameSocket'
import { Leaderboard, Podium } from '../teacher/HostPage'

type QuestionType = 'multiple_choice' | 'essay'

type PublicQuestion = {
  id: string
  type: QuestionType
  text: string
  options?: string[]
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

const COLORS = ['#e21b3c', '#1368ce', '#d89e00', '#26890c']

export default function PlayPage() {
  const { connect, send, on } = useGameSocket()
  const [pin, setPin] = useState('')
  const [nickname, setNickname] = useState('')
  const [ra, setRa] = useState('')
  const [avatar, setAvatar] = useState(() => Math.floor(Math.random() * AVATARS.length))
  const [phase, setPhase] = useState<'join' | 'lobby' | 'question' | 'reveal' | 'finished'>('join')
  const [playerId, setPlayerId] = useState('')
  const [quizTitle, setQuizTitle] = useState('')
  const [question, setQuestion] = useState<PublicQuestion | null>(null)
  const [endsAt, setEndsAt] = useState<number | null>(null)
  const [selected, setSelected] = useState<number | null>(null)
  const [essayDraft, setEssayDraft] = useState('')
  const [submittedText, setSubmittedText] = useState('')
  const [locked, setLocked] = useState(false)
  const [correctIndex, setCorrectIndex] = useState<number | null>(null)
  const [expectedAnswer, setExpectedAnswer] = useState('')
  const [lastPoints, setLastPoints] = useState<number | null>(null)
  const [lastSimilarity, setLastSimilarity] = useState<number | null>(null)
  const [lastCorrect, setLastCorrect] = useState<boolean | null>(null)
  const [board, setBoard] = useState<BoardEntry[]>([])
  const [error, setError] = useState('')
  const [autoNextIn, setAutoNextIn] = useState<number | null>(null)
  const remaining = useCountdown(endsAt)
  const timerPct = useMemo(() => {
    if (!question || !endsAt) return 0
    const total = question.timeLimitSec * 1000
    const left = Math.max(0, endsAt - Date.now())
    return Math.min(100, (left / total) * 100)
  }, [question, endsAt, remaining])

  useEffect(() => {
    if (autoNextIn == null || autoNextIn <= 0) return
    const id = setTimeout(() => setAutoNextIn((v) => (v == null ? null : v - 1)), 1000)
    return () => clearTimeout(id)
  }, [autoNextIn])

  useEffect(() => {
    const offs = [
      on('joined', (data) => {
        const d = data as { playerId: string; quizTitle: string; avatar?: number }
        setPlayerId(d.playerId)
        setQuizTitle(d.quizTitle)
        if (typeof d.avatar === 'number') setAvatar(d.avatar)
        setPhase('lobby')
        setError('')
      }),
      on('question', (data) => {
        const d = data as { question: PublicQuestion; endsAt: string }
        setPhase('question')
        setQuestion({
          ...d.question,
          type: d.question.type || 'multiple_choice',
        })
        setEndsAt(new Date(d.endsAt).getTime())
        setSelected(null)
        setEssayDraft('')
        setSubmittedText('')
        setLocked(false)
        setCorrectIndex(null)
        setExpectedAnswer('')
        setLastPoints(null)
        setLastSimilarity(null)
        setLastCorrect(null)
        setAutoNextIn(null)
        setError('')
      }),
      on('answer_ack', (data) => {
        const d = data as { points: number; similarity?: number }
        setLocked(true)
        setLastPoints(d.points)
        if (typeof d.similarity === 'number') setLastSimilarity(d.similarity)
      }),
      on('question_result', (data) => {
        const d = data as {
          type?: QuestionType
          correctIndex: number
          expectedAnswer?: string
          leaderboard: BoardEntry[]
          results: ResultEntry[]
          autoNextInSec?: number
        }
        setPhase('reveal')
        setCorrectIndex(d.correctIndex)
        setExpectedAnswer(d.expectedAnswer || '')
        setBoard(d.leaderboard)
        setEndsAt(null)
        setAutoNextIn(d.autoNextInSec ?? 5)
        const mine = d.results.find((r) => r.playerId === playerId)
        if (mine) {
          setLastPoints(mine.points)
          setLastCorrect(mine.correct)
          if (typeof mine.similarity === 'number') setLastSimilarity(mine.similarity)
          if (mine.text) setSubmittedText(mine.text)
        }
      }),
      on('finished', (data) => {
        setPhase('finished')
        setBoard((data as { leaderboard: BoardEntry[] }).leaderboard)
        setAutoNextIn(null)
      }),
      on('error', (data) => setError((data as { message: string }).message)),
    ]
    return () => offs.forEach((off) => off())
  }, [on, playerId])

  function join(e: FormEvent) {
    e.preventDefault()
    setError('')
    connect()
    send('join', { pin: pin.trim(), nickname: nickname.trim(), ra: ra.trim(), avatar })
  }

  function answerChoice(choice: number) {
    if (locked || phase !== 'question') return
    setSelected(choice)
    setError('')
    send('answer', { choice })
  }

  function submitEssay(e?: FormEvent) {
    e?.preventDefault()
    if (locked || phase !== 'question') return
    const text = essayDraft.trim()
    if (!text) {
      setError('Escreva sua resposta antes de enviar')
      return
    }
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
        </form>
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

          <h2>{question.text}</h2>

          {phase === 'question' && !isEssay && (
            <div className="options-grid play-opts">
              {(question.options || []).map((opt, i) => {
                let cls = 'opt'
                if (selected === i) cls += ' selected'
                return (
                  <button
                    key={i}
                    className={cls}
                    style={{ ['--opt' as string]: COLORS[i % COLORS.length] }}
                    disabled={locked}
                    onClick={() => answerChoice(i)}
                  >
                    <span className="opt-letter">{String.fromCharCode(65 + i)}</span>
                    <span className="opt-text">{opt}</span>
                  </button>
                )
              })}
            </div>
          )}

          {phase === 'question' && isEssay && !locked && (
            <form className="essay-box" onSubmit={submitEssay}>
              <textarea
                value={essayDraft}
                onChange={(e) => setEssayDraft(e.target.value.slice(0, 2000))}
                placeholder="Digite sua resposta com suas palavras..."
                rows={6}
                autoFocus
              />
              <div className="essay-actions">
                <span className="muted tiny">{essayDraft.trim().length}/2000</span>
                <button className="btn btn-accent btn-xl" type="submit" disabled={!essayDraft.trim()}>
                  Enviar resposta
                </button>
              </div>
            </form>
          )}

          {phase === 'question' && isEssay && locked && (
            <div className="waiting-card">
              <div className="pulse-ring small" />
              <p>Resposta enviada!</p>
              <p className="muted">Aguarde o tempo ou os demais jogadores...</p>
              <blockquote className="sent-answer">{submittedText || essayDraft}</blockquote>
            </div>
          )}

          {phase === 'question' && !isEssay && locked && (
            <div className="waiting-card">
              <div className="pulse-ring small" />
              <p>Resposta enviada! Aguarde o resultado...</p>
            </div>
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
                    {lastCorrect ? 'Resposta aceita' : 'Abaixo do limiar'}
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
            onClick={() => {
              setPhase('join')
              setPin('')
              setNickname('')
              setRa('')
              setPlayerId('')
              setAvatar(Math.floor(Math.random() * AVATARS.length))
            }}
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

import { useEffect, useState, type FormEvent } from 'react'
import { Link, Navigate, useNavigate, useParams } from 'react-router-dom'
import { api, type Question, type QuestionType, type Quiz } from '../../lib/api'
import { useAuth } from '../../lib/auth'

const emptyForm = {
  type: 'multiple_choice' as QuestionType,
  text: '',
  options: ['', '', '', ''],
  correctIndex: 0,
  expectedAnswer: '',
  similarityThreshold: 0.55,
  weight: 1,
  timeLimitMin: 1,
}

const OPTION_COLORS = ['#e21b3c', '#1368ce', '#d89e00', '#26890c']

function secToMin(sec: number) {
  if (!sec || sec <= 0) return 1
  return Math.max(1, Math.round(sec / 60))
}

function formatDuration(sec: number) {
  if (sec < 60) return `${sec}s`
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return s === 0 ? `${m} min` : `${m} min ${s}s`
}

export default function QuizEditorPage() {
  const { id } = useParams()
  const { teacher, token, loading } = useAuth()
  const nav = useNavigate()
  const [quiz, setQuiz] = useState<Quiz | null>(null)
  const [questions, setQuestions] = useState<Question[]>([])
  const [form, setForm] = useState(emptyForm)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [error, setError] = useState('')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')

  useEffect(() => {
    if (!token || !id) return
    ;(async () => {
      try {
        const [qList, quizData] = await Promise.all([
          api.listQuestions(token, id),
          api.getQuiz(token, id),
        ])
        setQuiz(quizData)
        setTitle(quizData.title)
        setDescription(quizData.description)
        setQuestions(qList)
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Erro')
      }
    })()
  }, [token, id])

  if (!loading && !teacher) return <Navigate to="/teacher/login" replace />

  async function saveMeta() {
    if (!token || !id) return
    const updated = await api.updateQuiz(token, id, title, description)
    setQuiz(updated)
  }

  function setType(type: QuestionType) {
    setForm((prev) => ({
      ...prev,
      type,
      timeLimitMin: type === 'essay' ? Math.max(prev.timeLimitMin, 2) : prev.timeLimitMin,
    }))
  }

  function startEdit(q: Question) {
    setEditingId(q.id)
    setForm({
      type: q.type || 'multiple_choice',
      text: q.text,
      options: [...(q.options || []), '', '', '', ''].slice(0, 4),
      correctIndex: q.correctIndex ?? 0,
      expectedAnswer: q.expectedAnswer || '',
      similarityThreshold: q.similarityThreshold || 0.55,
      weight: q.weight,
      timeLimitMin: secToMin(q.timeLimitSec),
    })
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    if (!token || !id) return
    setError('')

    const minutes = Math.min(10, Math.max(1, Number(form.timeLimitMin) || 1))
    const payload: Partial<Question> = {
      type: form.type,
      text: form.text,
      weight: Number(form.weight) || 1,
      timeLimitSec: minutes * 60,
      order: editingId ? questions.find((q) => q.id === editingId)?.order ?? questions.length : questions.length,
    }

    if (form.type === 'essay') {
      if (!form.expectedAnswer.trim()) {
        setError('Informe a resposta esperada')
        return
      }
      payload.expectedAnswer = form.expectedAnswer.trim()
      payload.similarityThreshold = Number(form.similarityThreshold) || 0.55
      payload.options = []
      payload.correctIndex = -1
    } else {
      const options = form.options.map((o) => o.trim()).filter(Boolean)
      if (options.length < 2) {
        setError('Informe pelo menos 2 opções')
        return
      }
      payload.options = options
      payload.correctIndex = Math.min(form.correctIndex, options.length - 1)
      payload.expectedAnswer = ''
    }

    try {
      if (editingId) {
        const updated = await api.updateQuestion(token, id, editingId, payload)
        setQuestions((prev) => prev.map((q) => (q.id === editingId ? updated : q)))
      } else {
        const created = await api.createQuestion(token, id, payload)
        setQuestions((prev) => [...prev, created])
      }
      setForm(emptyForm)
      setEditingId(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao salvar')
    }
  }

  async function removeQuestion(qid: string) {
    if (!token || !id || !confirm('Excluir questão?')) return
    await api.deleteQuestion(token, id, qid)
    setQuestions((prev) => prev.filter((q) => q.id !== qid))
  }

  async function goLive() {
    if (!token || !id) return
    try {
      const session = await api.createSession(token, id)
      nav(`/teacher/host/${session.pin}`, { state: session })
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Erro')
    }
  }

  return (
    <div className="teacher-shell">
      <header className="teacher-top">
        <Link to="/teacher" className="brand-link">
          ← Quizzes
        </Link>
        <button className="btn btn-accent" onClick={goLive} disabled={questions.length === 0}>
          Iniciar ao vivo
        </button>
      </header>

      <section className="panel">
        <h1>Editor de quiz</h1>
        <div className="create-row">
          <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Título" />
          <button className="btn btn-ghost" onClick={saveMeta}>
            Salvar título
          </button>
        </div>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Descrição (opcional)"
          rows={2}
        />
      </section>

      <div className="editor-grid">
        <section className="panel">
          <h2>{editingId ? 'Editar questão' : 'Nova questão'}</h2>
          <form className="question-form" onSubmit={onSubmit}>
            <div className="type-toggle" role="group" aria-label="Tipo de questão">
              <button
                type="button"
                className={`chip ${form.type === 'multiple_choice' ? 'active' : ''}`}
                onClick={() => setType('multiple_choice')}
              >
                Objetiva
              </button>
              <button
                type="button"
                className={`chip ${form.type === 'essay' ? 'active' : ''}`}
                onClick={() => setType('essay')}
              >
                Dissertativa
              </button>
            </div>

            <label>
              Enunciado
              <textarea
                required
                value={form.text}
                onChange={(e) => setForm({ ...form, text: e.target.value })}
                rows={3}
              />
            </label>

            {form.type === 'multiple_choice' ? (
              <div className="option-editor-list">
                {form.options.map((opt, i) => (
                  <button
                    key={i}
                    type="button"
                    className={`option-editor ${form.correctIndex === i ? 'is-correct' : ''}`}
                    style={{ ['--opt' as string]: OPTION_COLORS[i] }}
                    onClick={() => setForm({ ...form, correctIndex: i })}
                  >
                    <span className="option-editor-letter">{String.fromCharCode(65 + i)}</span>
                    <textarea
                      value={opt}
                      rows={3}
                      placeholder={`Alternativa ${String.fromCharCode(65 + i)}`}
                      onClick={(e) => e.stopPropagation()}
                      onChange={(e) => {
                        const options = [...form.options]
                        options[i] = e.target.value
                        setForm({ ...form, options })
                      }}
                    />
                    <span className="option-editor-mark">
                      {form.correctIndex === i ? 'Correta' : 'Marcar'}
                    </span>
                  </button>
                ))}
              </div>
            ) : (
              <>
                <label>
                  Resposta esperada
                  <textarea
                    required
                    value={form.expectedAnswer}
                    onChange={(e) => setForm({ ...form, expectedAnswer: e.target.value })}
                    rows={4}
                    placeholder="Texto de referência para correção por similaridade"
                  />
                </label>
                <label>
                  Limiar de similaridade ({Math.round(form.similarityThreshold * 100)}%)
                  <input
                    type="range"
                    min={0.3}
                    max={0.95}
                    step={0.05}
                    value={form.similarityThreshold}
                    onChange={(e) => setForm({ ...form, similarityThreshold: Number(e.target.value) })}
                  />
                </label>
                <p className="muted tiny">
                  XP proporcional à similaridade. Acima do limiar conta como “acertou”.
                </p>
              </>
            )}

            <div className="meta-row">
              <label>
                Peso
                <input
                  type="number"
                  min={0.5}
                  step={0.5}
                  value={form.weight}
                  onChange={(e) => setForm({ ...form, weight: Number(e.target.value) })}
                />
              </label>
              <label>
                Tempo (minutos)
                <input
                  type="number"
                  min={1}
                  max={10}
                  step={1}
                  value={form.timeLimitMin}
                  onChange={(e) => setForm({ ...form, timeLimitMin: Number(e.target.value) })}
                />
              </label>
            </div>
            {error && <p className="error">{error}</p>}
            <div className="row-actions">
              <button className="btn btn-primary" type="submit">
                {editingId ? 'Atualizar' : 'Adicionar'}
              </button>
              {editingId && (
                <button
                  type="button"
                  className="btn btn-ghost"
                  onClick={() => {
                    setEditingId(null)
                    setForm(emptyForm)
                  }}
                >
                  Cancelar
                </button>
              )}
            </div>
          </form>
        </section>

        <section className="panel">
          <h2>Questões ({questions.length})</h2>
          {questions.length === 0 && <p className="muted">Nenhuma questão ainda.</p>}
          <ul className="q-list">
            {questions.map((q, idx) => (
              <li key={q.id}>
                <div>
                  <strong>
                    {idx + 1}. {q.text}
                  </strong>
                  <p className="muted">
                    {q.type === 'essay' ? 'Dissertativa' : 'Objetiva'} · Peso {q.weight} ·{' '}
                    {formatDuration(q.timeLimitSec)}
                    {q.type === 'essay'
                      ? ` · limiar ${Math.round((q.similarityThreshold || 0.55) * 100)}%`
                      : ` · correta: ${q.options?.[q.correctIndex] ?? '—'}`}
                  </p>
                </div>
                <div className="row-actions">
                  <button className="btn btn-ghost" onClick={() => startEdit(q)}>
                    Editar
                  </button>
                  <button className="btn btn-danger" onClick={() => removeQuestion(q.id)}>
                    Excluir
                  </button>
                </div>
              </li>
            ))}
          </ul>
          {quiz && <p className="muted tiny">ID: {quiz.id}</p>}
        </section>
      </div>
    </div>
  )
}

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export type Teacher = { id: string; email: string; name: string }

export type Quiz = {
  id: string
  teacherId: string
  title: string
  description: string
  createdAt: string
  updatedAt: string
}

export type QuestionType = 'multiple_choice' | 'essay'

export type Question = {
  id: string
  quizId: string
  type: QuestionType
  text: string
  options: string[]
  correctIndex: number
  expectedAnswer?: string
  expectedAnswers?: string[]
  similarityThreshold?: number
  codeSnippet?: string
  codeLanguage?: string
  weight: number
  timeLimitSec: number
  order: number
}

export type RankingEntry = {
  playerId: string
  nickname: string
  ra: string
  score: number
  correctCount: number
  total: number
  maxScore: number
  grade: number
  rank: number
}

export type SessionRecord = {
  id: string
  quizId: string
  quizTitle: string
  teacherId: string
  pin: string
  ranking: RankingEntry[]
  startedAt: string
  finishedAt: string
}

function authHeaders(token?: string | null): HeadersInit {
  const h: HeadersInit = { 'Content-Type': 'application/json' }
  if (token) h.Authorization = `Bearer ${token}`
  return h
}

async function request<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, opts)
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || res.statusText)
  return data as T
}

export const api = {
  health: () => request<{ status: string; authMode: string; lanIP?: string }>('/api/health'),
  authMode: () => request<{ mode: string }>('/api/auth/mode'),
  devLogin: (email: string, password: string, name?: string) =>
    request<{ token: string; teacher: Teacher }>('/api/auth/dev-login', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ email, password, name }),
    }),
  session: (token: string) =>
    request<{ teacher: Teacher; authMode: string }>('/api/auth/session', {
      method: 'POST',
      headers: authHeaders(token),
      body: JSON.stringify({ token }),
    }),
  listQuizzes: (token: string) =>
    request<Quiz[]>('/api/quizzes', { headers: authHeaders(token) }),
  getQuiz: (token: string, id: string) =>
    request<Quiz>(`/api/quizzes/${id}`, { headers: authHeaders(token) }),
  createQuiz: (token: string, title: string, description: string) =>
    request<Quiz>('/api/quizzes', {
      method: 'POST',
      headers: authHeaders(token),
      body: JSON.stringify({ title, description }),
    }),
  updateQuiz: (token: string, id: string, title: string, description: string) =>
    request<Quiz>(`/api/quizzes/${id}`, {
      method: 'PUT',
      headers: authHeaders(token),
      body: JSON.stringify({ title, description }),
    }),
  deleteQuiz: (token: string, id: string) =>
    request<{ status: string }>(`/api/quizzes/${id}`, {
      method: 'DELETE',
      headers: authHeaders(token),
    }),
  listQuestions: (token: string, quizId: string) =>
    request<Question[]>(`/api/quizzes/${quizId}/questions`, { headers: authHeaders(token) }),
  createQuestion: (token: string, quizId: string, q: Partial<Question>) =>
    request<Question>(`/api/quizzes/${quizId}/questions`, {
      method: 'POST',
      headers: authHeaders(token),
      body: JSON.stringify(q),
    }),
  updateQuestion: (token: string, quizId: string, qid: string, q: Partial<Question>) =>
    request<Question>(`/api/quizzes/${quizId}/questions/${qid}`, {
      method: 'PUT',
      headers: authHeaders(token),
      body: JSON.stringify(q),
    }),
  deleteQuestion: (token: string, quizId: string, qid: string) =>
    request<{ status: string }>(`/api/quizzes/${quizId}/questions/${qid}`, {
      method: 'DELETE',
      headers: authHeaders(token),
    }),
  previewSimilarity: (
    token: string,
    body: {
      text: string
      expectedAnswer: string
      expectedAnswers?: string[]
      threshold?: number
    },
  ) =>
    request<{
      similarity: number
      passed: boolean
      threshold: number
      matches: { answer: string; similarity: number }[]
    }>('/api/similarity', {
      method: 'POST',
      headers: authHeaders(token),
      body: JSON.stringify(body),
    }),
  createSession: (token: string, quizId: string) =>
    request<{ pin: string; sessionId: string; quizId: string; quizTitle: string }>(
      '/api/sessions',
      {
        method: 'POST',
        headers: authHeaders(token),
        body: JSON.stringify({ quizId }),
      },
    ),
  listSessions: (token: string) =>
    request<SessionRecord[]>('/api/sessions', { headers: authHeaders(token) }),
  getSession: (token: string, id: string) =>
    request<SessionRecord>(`/api/sessions/${id}`, { headers: authHeaders(token) }),
}

export function wsUrl(): string {
  const env = import.meta.env.VITE_WS_URL as string | undefined
  if (env) return env
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${window.location.host}/ws`
}

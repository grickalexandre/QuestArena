import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AuthProvider } from './lib/auth'
import HomePage from './pages/HomePage'
import TeacherLoginPage from './pages/teacher/LoginPage'
import QuizzesPage from './pages/teacher/QuizzesPage'
import QuizEditorPage from './pages/teacher/QuizEditorPage'
import HostPage from './pages/teacher/HostPage'
import PlayPage from './pages/play/PlayPage'

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/teacher/login" element={<TeacherLoginPage />} />
          <Route path="/teacher" element={<QuizzesPage />} />
          <Route path="/teacher/quizzes/:id" element={<QuizEditorPage />} />
          <Route path="/teacher/host/:pin" element={<HostPage />} />
          <Route path="/play" element={<PlayPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

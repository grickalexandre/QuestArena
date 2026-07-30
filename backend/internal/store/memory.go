package store

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/questarena/questarena/internal/models"
)

// MemoryStore is used when Firebase is not configured (local/dev).
type MemoryStore struct {
	mu        sync.RWMutex
	teachers  map[string]*models.Teacher
	quizzes   map[string]*models.Quiz
	questions map[string][]*models.Question // quizID -> questions
	sessions  map[string]*models.SessionRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		teachers:  make(map[string]*models.Teacher),
		quizzes:   make(map[string]*models.Quiz),
		questions: make(map[string][]*models.Question),
		sessions:  make(map[string]*models.SessionRecord),
	}
}

func (m *MemoryStore) UpsertTeacher(_ context.Context, t *models.Teacher) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *t
	m.teachers[t.ID] = &cp
	return nil
}

func (m *MemoryStore) GetTeacher(_ context.Context, id string) (*models.Teacher, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.teachers[id]
	if !ok {
		return nil, fmt.Errorf("teacher not found")
	}
	cp := *t
	return &cp, nil
}

func (m *MemoryStore) CreateQuiz(_ context.Context, q *models.Quiz) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if q.ID == "" {
		q.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	q.CreatedAt = now
	q.UpdatedAt = now
	cp := *q
	m.quizzes[q.ID] = &cp
	return nil
}

func (m *MemoryStore) UpdateQuiz(_ context.Context, q *models.Quiz) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	existing, ok := m.quizzes[q.ID]
	if !ok {
		return fmt.Errorf("quiz not found")
	}
	if existing.TeacherID != q.TeacherID {
		return fmt.Errorf("forbidden")
	}
	q.CreatedAt = existing.CreatedAt
	q.UpdatedAt = time.Now().UTC()
	cp := *q
	m.quizzes[q.ID] = &cp
	return nil
}

func (m *MemoryStore) DeleteQuiz(_ context.Context, teacherID, quizID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	q, ok := m.quizzes[quizID]
	if !ok {
		return fmt.Errorf("quiz not found")
	}
	if q.TeacherID != teacherID {
		return fmt.Errorf("forbidden")
	}
	delete(m.quizzes, quizID)
	delete(m.questions, quizID)
	return nil
}

func (m *MemoryStore) GetQuiz(_ context.Context, quizID string) (*models.Quiz, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	q, ok := m.quizzes[quizID]
	if !ok {
		return nil, fmt.Errorf("quiz not found")
	}
	cp := *q
	return &cp, nil
}

func (m *MemoryStore) ListQuizzes(_ context.Context, teacherID string) ([]models.Quiz, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.Quiz, 0)
	for _, q := range m.quizzes {
		if q.TeacherID == teacherID {
			out = append(out, *q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].UpdatedAt.After(out[j].UpdatedAt)
	})
	return out, nil
}

func (m *MemoryStore) CreateQuestion(_ context.Context, q *models.Question) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.quizzes[q.QuizID]; !ok {
		return fmt.Errorf("quiz not found")
	}
	if q.ID == "" {
		q.ID = uuid.NewString()
	}
	q.CreatedAt = time.Now().UTC()
	cp := *q
	opts := make([]string, len(q.Options))
	copy(opts, q.Options)
	cp.Options = opts
	m.questions[q.QuizID] = append(m.questions[q.QuizID], &cp)
	return nil
}

func (m *MemoryStore) UpdateQuestion(_ context.Context, q *models.Question) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.questions[q.QuizID]
	for i, existing := range list {
		if existing.ID == q.ID {
			q.CreatedAt = existing.CreatedAt
			cp := *q
			opts := make([]string, len(q.Options))
			copy(opts, q.Options)
			cp.Options = opts
			list[i] = &cp
			return nil
		}
	}
	return fmt.Errorf("question not found")
}

func (m *MemoryStore) DeleteQuestion(_ context.Context, quizID, questionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.questions[quizID]
	for i, q := range list {
		if q.ID == questionID {
			m.questions[quizID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("question not found")
}

func (m *MemoryStore) ListQuestions(_ context.Context, quizID string) ([]models.Question, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := m.questions[quizID]
	out := make([]models.Question, 0, len(list))
	for _, q := range list {
		cp := *q
		opts := make([]string, len(q.Options))
		copy(opts, q.Options)
		cp.Options = opts
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Order < out[j].Order
	})
	return out, nil
}

func (m *MemoryStore) SaveSession(_ context.Context, s *models.SessionRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	cp := *s
	m.sessions[s.ID] = &cp
	return nil
}

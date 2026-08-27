package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/questarena/questarena/internal/models"
)

// MemoryStore is used when Firebase is not configured (local/dev).
type MemoryStore struct {
	mu        sync.RWMutex
	path      string
	teachers  map[string]*models.Teacher
	quizzes   map[string]*models.Quiz
	questions map[string][]*models.Question // quizID -> questions
	sessions  map[string]*models.SessionRecord
}

type memorySnapshot struct {
	Teachers  []models.Teacher                 `json:"teachers"`
	Quizzes   []models.Quiz                    `json:"quizzes"`
	Questions map[string][]models.Question     `json:"questions"`
	Sessions  []models.SessionRecord           `json:"sessions"`
}

func NewMemoryStore() *MemoryStore {
	return newMemoryStore("")
}

func NewPersistentMemoryStore(path string) *MemoryStore {
	m := newMemoryStore(path)
	_ = m.load()
	return m
}

func newMemoryStore(path string) *MemoryStore {
	return &MemoryStore{
		path:      path,
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
	return m.persistLocked()
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
	return m.persistLocked()
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
	return m.persistLocked()
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
	return m.persistLocked()
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
	return m.persistLocked()
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
			return m.persistLocked()
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
			return m.persistLocked()
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
	if s.Ranking != nil {
		cp.Ranking = make([]models.RankingEntry, len(s.Ranking))
		copy(cp.Ranking, s.Ranking)
	}
	m.sessions[s.ID] = &cp
	return m.persistLocked()
}

func (m *MemoryStore) ListSessions(_ context.Context, teacherID string) ([]models.SessionRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.SessionRecord, 0)
	for _, s := range m.sessions {
		if s.TeacherID == teacherID {
			cp := *s
			if s.Ranking != nil {
				cp.Ranking = make([]models.RankingEntry, len(s.Ranking))
				copy(cp.Ranking, s.Ranking)
			}
			out = append(out, cp)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].FinishedAt.After(out[j].FinishedAt)
	})
	return out, nil
}

func (m *MemoryStore) GetSession(_ context.Context, teacherID, sessionID string) (*models.SessionRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	if s.TeacherID != teacherID {
		return nil, fmt.Errorf("forbidden")
	}
	cp := *s
	if s.Ranking != nil {
		cp.Ranking = make([]models.RankingEntry, len(s.Ranking))
		copy(cp.Ranking, s.Ranking)
	}
	return &cp, nil
}

func (m *MemoryStore) load() error {
	if m.path == "" {
		return nil
	}
	raw, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var snap memorySnapshot
	if err := json.Unmarshal(raw, &snap); err != nil {
		return err
	}
	for _, t := range snap.Teachers {
		cp := t
		m.teachers[cp.ID] = &cp
	}
	for _, q := range snap.Quizzes {
		cp := q
		m.quizzes[cp.ID] = &cp
	}
	for quizID, list := range snap.Questions {
		copied := make([]*models.Question, 0, len(list))
		for _, q := range list {
			cp := q
			opts := make([]string, len(q.Options))
			copy(opts, q.Options)
			cp.Options = opts
			copied = append(copied, &cp)
		}
		m.questions[quizID] = copied
	}
	for _, s := range snap.Sessions {
		cp := s
		if s.Ranking != nil {
			cp.Ranking = make([]models.RankingEntry, len(s.Ranking))
			copy(cp.Ranking, s.Ranking)
		}
		m.sessions[cp.ID] = &cp
	}
	return nil
}

func (m *MemoryStore) persistLocked() error {
	if m.path == "" {
		return nil
	}
	snap := memorySnapshot{
		Teachers:  make([]models.Teacher, 0, len(m.teachers)),
		Quizzes:   make([]models.Quiz, 0, len(m.quizzes)),
		Questions: make(map[string][]models.Question, len(m.questions)),
		Sessions:  make([]models.SessionRecord, 0, len(m.sessions)),
	}
	for _, t := range m.teachers {
		snap.Teachers = append(snap.Teachers, *t)
	}
	for _, q := range m.quizzes {
		snap.Quizzes = append(snap.Quizzes, *q)
	}
	for quizID, list := range m.questions {
		out := make([]models.Question, 0, len(list))
		for _, q := range list {
			cp := *q
			opts := make([]string, len(q.Options))
			copy(opts, q.Options)
			cp.Options = opts
			out = append(out, cp)
		}
		snap.Questions[quizID] = out
	}
	for _, s := range m.sessions {
		cp := *s
		if s.Ranking != nil {
			cp.Ranking = make([]models.RankingEntry, len(s.Ranking))
			copy(cp.Ranking, s.Ranking)
		}
		snap.Sessions = append(snap.Sessions, cp)
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp := m.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.path)
}

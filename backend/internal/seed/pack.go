package seed

import (
	"context"
	"slices"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

const timeLimitOneMin = 60

type draftQuestion struct {
	text         string
	options      []string
	correctIndex int
	code         string
	codeLanguage string
}

func (d draftQuestion) toQuestion(quizID string, order int) *models.Question {
	q := &models.Question{
		QuizID:       quizID,
		Type:         models.QuestionMultipleChoice,
		Text:         d.text,
		Options:      append([]string{}, d.options...),
		CorrectIndex: d.correctIndex,
		CodeSnippet:  d.code,
		Weight:       1,
		TimeLimitSec: timeLimitOneMin,
		Order:        order,
	}
	if d.code != "" {
		q.CodeLanguage = models.NormalizeCodeLanguage(d.codeLanguage)
	}
	return q
}

func sameContent(a, b *models.Question) bool {
	return a.Type == b.Type &&
		a.Text == b.Text &&
		a.CorrectIndex == b.CorrectIndex &&
		a.CodeSnippet == b.CodeSnippet &&
		a.CodeLanguage == b.CodeLanguage &&
		a.Weight == b.Weight &&
		a.TimeLimitSec == b.TimeLimitSec &&
		slices.Equal(a.Options, b.Options)
}

// pack is a fixed quiz that every teacher receives automatically on first access.
type pack struct {
	idPrefix  string
	title     string
	desc      string
	questions []draftQuestion
}

func (p pack) quizID(teacherID string) string {
	return p.idPrefix + teacherID
}

func (p pack) ensure(ctx context.Context, st store.Store, teacherID string) error {
	if teacherID == "" {
		return nil
	}
	id := p.quizID(teacherID)
	existing, err := st.GetQuiz(ctx, id)
	if err != nil || existing == nil {
		q := &models.Quiz{
			ID:          id,
			TeacherID:   teacherID,
			Title:       p.title,
			Description: p.desc,
		}
		if err := st.CreateQuiz(ctx, q); err != nil {
			return err
		}
	} else if existing.Title != p.title || existing.Description != p.desc {
		existing.Title = p.title
		existing.Description = p.desc
		if err := st.UpdateQuiz(ctx, existing); err != nil {
			return err
		}
	}
	qs, err := st.ListQuestions(ctx, id)
	if err != nil {
		return err
	}
	current := make(map[int]models.Question, len(qs))
	for _, q := range qs {
		if _, taken := current[q.Order]; !taken {
			current[q.Order] = q
		}
	}
	// Questions the teacher added after the seeded ones keep their own order and
	// are never touched here.
	for i, d := range p.questions {
		want := d.toQuestion(id, i)
		saved, ok := current[i]
		if !ok {
			if err := st.CreateQuestion(ctx, want); err != nil {
				return err
			}
			continue
		}
		if sameContent(&saved, want) {
			continue
		}
		want.ID = saved.ID
		if err := st.UpdateQuestion(ctx, want); err != nil {
			return err
		}
	}
	return nil
}

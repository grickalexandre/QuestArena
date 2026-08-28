package seed

import (
	"context"

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
	if len(qs) >= len(p.questions) {
		return nil
	}
	for i, d := range p.questions {
		already := false
		for _, q := range qs {
			if q.Order == i {
				already = true
				break
			}
		}
		if already {
			continue
		}
		item := &models.Question{
			QuizID:       id,
			Type:         models.QuestionMultipleChoice,
			Text:         d.text,
			Options:      append([]string{}, d.options...),
			CorrectIndex: d.correctIndex,
			CodeSnippet:  d.code,
			Weight:       1,
			TimeLimitSec: timeLimitOneMin,
			Order:        i,
		}
		if d.code != "" {
			item.CodeLanguage = models.NormalizeCodeLanguage(d.codeLanguage)
		}
		if err := st.CreateQuestion(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

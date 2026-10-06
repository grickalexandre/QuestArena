package seed

import (
	"context"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestRequisitosQuestions(t *testing.T) {
	qs := requisitosQuestions()
	if len(qs) != 20 {
		t.Fatalf("want 20 questions, got %d", len(qs))
	}
	for i, d := range qs {
		if d.text == "" {
			t.Errorf("q%d: empty text", i+1)
		}
		if len(d.options) != 4 {
			t.Errorf("q%d: want 4 options, got %d", i+1, len(d.options))
		}
		if d.correctIndex < 0 || d.correctIndex >= len(d.options) {
			t.Errorf("q%d: correctIndex %d out of range", i+1, d.correctIndex)
		}
		if d.timeLimitSec != timeLimitOneMin {
			t.Errorf("q%d: want %ds, got %ds", i+1, timeLimitOneMin, d.timeLimitSec)
		}
		q := d.toQuestion("quiz", i)
		if q.Type != models.QuestionMultipleChoice {
			t.Errorf("q%d: want multiple choice", i+1)
		}
		if q.TimeLimitSec != timeLimitOneMin {
			t.Errorf("q%d: stored time %d", i+1, q.TimeLimitSec)
		}
	}
}

func TestEnsureRequisitosQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-req"

	for i := 0; i < 2; i++ {
		if err := EnsureRequisitosQuiz(ctx, st, teacherID); err != nil {
			t.Fatalf("ensure run %d: %v", i+1, err)
		}
	}
	quizzes, err := st.ListQuizzes(ctx, teacherID)
	if err != nil {
		t.Fatalf("list quizzes: %v", err)
	}
	if len(quizzes) != 1 || !IsRequisitosSeedQuiz(quizzes[0].ID) {
		t.Fatalf("quiz not seeded: %+v", quizzes)
	}
	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 20 {
		t.Fatalf("want 20 questions, got %d", len(qs))
	}
}

package seed

import (
	"context"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestSorveteriaQuestions(t *testing.T) {
	qs := sorveteriaQuestions()
	if len(qs) != 10 {
		t.Fatalf("want 10 questions, got %d", len(qs))
	}
	mc, essay := 0, 0
	for i, d := range qs {
		if d.text == "" {
			t.Errorf("q%d: empty text", i+1)
		}
		q := d.toQuestion("quiz", i)
		if d.expectedAnswer != "" {
			essay++
			if q.Type != models.QuestionEssay {
				t.Errorf("q%d: want essay", i+1)
			}
			if q.TimeLimitSec != timeLimitFiveMin {
				t.Errorf("q%d: want %ds, got %ds", i+1, timeLimitFiveMin, q.TimeLimitSec)
			}
			if q.GradingCriteria == "" {
				t.Errorf("q%d: missing grading criteria", i+1)
			}
			if len(q.KeyTerms) == 0 {
				t.Errorf("q%d: missing key terms", i+1)
			}
			continue
		}
		mc++
		if len(d.options) != 4 {
			t.Errorf("q%d: want 4 options, got %d", i+1, len(d.options))
		}
		if d.correctIndex < 0 || d.correctIndex >= len(d.options) {
			t.Errorf("q%d: correctIndex %d out of range", i+1, d.correctIndex)
		}
		if q.TimeLimitSec != timeLimitOneMin {
			t.Errorf("q%d: want %ds, got %ds", i+1, timeLimitOneMin, q.TimeLimitSec)
		}
	}
	if mc != 7 || essay != 3 {
		t.Fatalf("want 7 MC + 3 essay, got %d MC + %d essay", mc, essay)
	}
}

func TestEnsureSorveteriaQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-sorveteria"

	for i := 0; i < 2; i++ {
		if err := EnsureSorveteriaQuiz(ctx, st, teacherID); err != nil {
			t.Fatalf("ensure run %d: %v", i+1, err)
		}
	}

	quizzes, err := st.ListQuizzes(ctx, teacherID)
	if err != nil {
		t.Fatalf("list quizzes: %v", err)
	}
	if len(quizzes) != 1 {
		t.Fatalf("want 1 quiz, got %d", len(quizzes))
	}
	if !IsSorveteriaSeedQuiz(quizzes[0].ID) {
		t.Errorf("quiz %q not recognized as seed", quizzes[0].ID)
	}

	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 10 {
		t.Fatalf("want 10 questions, got %d", len(qs))
	}
}

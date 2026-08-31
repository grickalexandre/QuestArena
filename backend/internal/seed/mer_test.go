package seed

import (
	"context"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestMerQuestions(t *testing.T) {
	qs := merQuestions()
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
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
			if q.ExpectedAnswer == "" {
				t.Errorf("q%d: missing expectedAnswer", i+1)
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
	if mc != 10 || essay != 5 {
		t.Fatalf("want 10 MC + 5 essay, got %d MC + %d essay", mc, essay)
	}
}

func TestEnsureMerQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-mer"

	for i := 0; i < 2; i++ {
		if err := EnsureMerQuiz(ctx, st, teacherID); err != nil {
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
	if !IsMerSeedQuiz(quizzes[0].ID) {
		t.Errorf("quiz %q not recognized as seed", quizzes[0].ID)
	}

	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}
}

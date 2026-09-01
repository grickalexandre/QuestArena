package seed

import (
	"context"
	"strings"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestArquiteturaQuestions(t *testing.T) {
	qs := arquiteturaQuestions()
	if len(qs) != 30 {
		t.Fatalf("want 30 questions, got %d", len(qs))
	}
	official := 0
	for i, d := range qs {
		if d.text == "" {
			t.Errorf("q%d: empty text", i+1)
		}
		if d.expectedAnswer != "" {
			t.Errorf("q%d: all questions must be multiple choice", i+1)
		}
		if len(d.options) != 4 {
			t.Errorf("q%d: want 4 options, got %d", i+1, len(d.options))
		}
		if d.correctIndex < 0 || d.correctIndex >= len(d.options) {
			t.Errorf("q%d: correctIndex %d out of range", i+1, d.correctIndex)
		}
		for j, opt := range d.options {
			if opt == "" {
				t.Errorf("q%d option %d empty", i+1, j)
			}
		}
		q := d.toQuestion("quiz", i)
		if q.Type != models.QuestionMultipleChoice {
			t.Errorf("q%d: want multiple_choice", i+1)
		}
		if d.timeLimitSec == 90 {
			if q.TimeLimitSec != 90 {
				t.Errorf("q%d: want 90s, got %ds", i+1, q.TimeLimitSec)
			}
		} else if q.TimeLimitSec != timeLimitOneMin {
			t.Errorf("q%d: want %ds, got %ds", i+1, timeLimitOneMin, q.TimeLimitSec)
		}
		if strings.HasPrefix(d.text, "[ENADE") {
			official++
		}
	}
	if official != 6 {
		t.Errorf("want 6 official ENADE items, got %d", official)
	}
	for i, d := range qs {
		if strings.Contains(strings.ToLower(d.text), "dissert") || strings.Contains(strings.ToLower(d.text), "discursiv") {
			t.Errorf("q%d: quiz must be alternatives only, found dissertativa/discursiva wording", i+1)
		}
	}
}

func TestEnsureArquiteturaQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-arq"

	for i := 0; i < 2; i++ {
		if err := EnsureArquiteturaQuiz(ctx, st, teacherID); err != nil {
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
	if !IsArquiteturaSeedQuiz(quizzes[0].ID) {
		t.Errorf("quiz %q not recognized as seed", quizzes[0].ID)
	}

	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 30 {
		t.Fatalf("want 30 questions, got %d", len(qs))
	}
}

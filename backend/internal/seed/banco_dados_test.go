package seed

import (
	"context"
	"strings"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestBancoDadosQuestions(t *testing.T) {
	qs := bancoDadosQuestions()
	if len(qs) != 30 {
		t.Fatalf("want 30 questions, got %d", len(qs))
	}
	official := 0
	joined := strings.Builder{}
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
		if d.timeLimitSec != 120 {
			t.Errorf("q%d: want 120s, got %ds", i+1, d.timeLimitSec)
		}
		for j, opt := range d.options {
			if opt == "" {
				t.Errorf("q%d option %d empty", i+1, j)
			}
		}
		if d.code != "" && models.NormalizeCodeLanguage(d.codeLanguage) != d.codeLanguage {
			t.Errorf("q%d: unsupported language %q", i+1, d.codeLanguage)
		}
		q := d.toQuestion("quiz", i)
		if q.Type != models.QuestionMultipleChoice {
			t.Errorf("q%d: want multiple_choice", i+1)
		}
		if q.TimeLimitSec != 120 {
			t.Errorf("q%d: toQuestion time %d", i+1, q.TimeLimitSec)
		}
		if strings.HasPrefix(d.text, "[ENADE") {
			official++
		}
		joined.WriteString(strings.ToLower(d.text))
		joined.WriteByte(' ')
		for _, opt := range d.options {
			joined.WriteString(strings.ToLower(opt))
			joined.WriteByte(' ')
		}
	}
	if official < 5 {
		t.Errorf("want at least 5 official/reconstructed ENADE items, got %d", official)
	}
	blob := joined.String()
	for _, need := range []string{"join", "group by", "isolament", "entidade", "normal", "acid"} {
		if !strings.Contains(blob, need) {
			t.Errorf("quiz must cover %q from the banco-de-dados material", need)
		}
	}
}

func TestEnsureBancoDadosQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-bd"

	for i := 0; i < 2; i++ {
		if err := EnsureBancoDadosQuiz(ctx, st, teacherID); err != nil {
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
	if !IsBancoDadosSeedQuiz(quizzes[0].ID) {
		t.Errorf("quiz %q not recognized as seed", quizzes[0].ID)
	}
	if quizzes[0].Title != bancoDadosTitle {
		t.Errorf("title: got %q", quizzes[0].Title)
	}

	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 30 {
		t.Fatalf("want 30 questions, got %d", len(qs))
	}
}

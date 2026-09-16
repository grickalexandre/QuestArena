package seed

import (
	"context"
	"strings"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestHerancaQuestions(t *testing.T) {
	qs := herancaQuestions()
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}
	joined := strings.Builder{}
	for i, q := range qs {
		if q.text == "" {
			t.Errorf("q%d: empty text", i+1)
		}
		if len(q.options) != 2 {
			t.Errorf("q%d: want 2 options (certo/errado), got %d", i+1, len(q.options))
		}
		if q.correctIndex < 0 || q.correctIndex >= len(q.options) {
			t.Errorf("q%d: correctIndex %d out of range", i+1, q.correctIndex)
		}
		if q.timeLimitSec != timeLimitHeranca {
			t.Errorf("q%d: want %ds (sem corrida de XP), got %ds", i+1, timeLimitHeranca, q.timeLimitSec)
		}
		seen := map[string]bool{}
		for j, opt := range q.options {
			if opt != "Certo" && opt != "Errado" {
				t.Errorf("q%d option %d: want Certo or Errado, got %q", i+1, j, opt)
			}
			seen[opt] = true
		}
		if !seen["Certo"] || !seen["Errado"] {
			t.Errorf("q%d: must offer both Certo and Errado", i+1)
		}
		low := strings.ToLower(q.text)
		for _, ban := range []string{"public", "private", "protected", "tempo de execução", "tostring", " invalidcast"} {
			if strings.Contains(low, ban) {
				t.Errorf("q%d: prova não cobre %q", i+1, ban)
			}
		}
		joined.WriteString(low)
		joined.WriteByte(' ')
		got := q.toQuestion("quiz", i)
		if got.Type != models.QuestionMultipleChoice {
			t.Errorf("q%d: want multiple_choice", i+1)
		}
		if got.TimeLimitSec != timeLimitHeranca {
			t.Errorf("q%d: toQuestion time %d", i+1, got.TimeLimitSec)
		}
	}
	blob := joined.String()
	for _, need := range []string{"heran", "abstract", "animal", "conta", ": base", "construtor", "polimorf"} {
		if !strings.Contains(blob, need) {
			t.Errorf("quiz must cover %q through oficina 6.1", need)
		}
	}
}

func TestEnsureHerancaQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-heranca"

	for i := 0; i < 2; i++ {
		if err := EnsureHerancaQuiz(ctx, st, teacherID); err != nil {
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
	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}
}

func TestEnsureHerancaRewritesOutdatedQuestions(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-heranca-2"

	if err := EnsureHerancaQuiz(ctx, st, teacherID); err != nil {
		t.Fatalf("first ensure: %v", err)
	}
	quizzes, err := st.ListQuizzes(ctx, teacherID)
	if err != nil {
		t.Fatalf("list quizzes: %v", err)
	}
	if len(quizzes) != 1 {
		t.Fatalf("want 1 quiz, got %d", len(quizzes))
	}
	if !IsHerancaSeedQuiz(quizzes[0].ID) {
		t.Errorf("quiz %q not recognized as seed", quizzes[0].ID)
	}
	if quizzes[0].Title != herancaTitle {
		t.Errorf("title: got %q", quizzes[0].Title)
	}

	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}

	stale := qs[0]
	stale.Text = "pergunta antiga com public e private"
	stale.Options = []string{"A", "B", "C", "D"}
	stale.CorrectIndex = 0
	if err := st.UpdateQuestion(ctx, &stale); err != nil {
		t.Fatalf("update: %v", err)
	}

	if err := EnsureHerancaQuiz(ctx, st, teacherID); err != nil {
		t.Fatalf("second ensure: %v", err)
	}
	qs, err = st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list after ensure: %v", err)
	}
	want := herancaQuestions()[0]
	for _, q := range qs {
		if q.Order != 0 {
			continue
		}
		if q.Text != want.text {
			t.Errorf("question 1 not restored: got %q", q.Text)
		}
		if q.ID != stale.ID {
			t.Errorf("should keep id %q, got %q", stale.ID, q.ID)
		}
		if len(q.Options) != 2 {
			t.Errorf("want 2 options after rewrite, got %d", len(q.Options))
		}
	}
}

package seed

import (
	"context"
	"testing"

	"github.com/questarena/questarena/internal/store"
)

func TestEnsureNodeSupabaseQuizIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-1"

	for i := 0; i < 2; i++ {
		if err := EnsureNodeSupabaseQuiz(ctx, st, teacherID); err != nil {
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
	if !IsNodeSupabaseSeedQuiz(quizzes[0].ID) {
		t.Errorf("quiz %q not recognized as seed", quizzes[0].ID)
	}

	qs, err := st.ListQuestions(ctx, quizzes[0].ID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}
	for _, q := range qs {
		if q.CodeSnippet != "" && q.CodeLanguage == "" {
			t.Errorf("question %q has code without language", q.Text)
		}
		if q.TimeLimitSec != timeLimitOneMin {
			t.Errorf("question %q: want %ds, got %ds", q.Text, timeLimitOneMin, q.TimeLimitSec)
		}
	}
}

func TestEnsureRewritesOutdatedQuestions(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	const teacherID = "teacher-2"

	if err := EnsureNodeSupabaseQuiz(ctx, st, teacherID); err != nil {
		t.Fatalf("first ensure: %v", err)
	}
	quizID := nodeSupabasePack().quizID(teacherID)
	qs, err := st.ListQuestions(ctx, quizID)
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}

	stale := qs[6]
	stale.Text = "pergunta antiga"
	stale.Options = []string{"a", "b", "c", "d"}
	if err := st.UpdateQuestion(ctx, &stale); err != nil {
		t.Fatalf("update question: %v", err)
	}

	if err := EnsureNodeSupabaseQuiz(ctx, st, teacherID); err != nil {
		t.Fatalf("second ensure: %v", err)
	}

	qs, err = st.ListQuestions(ctx, quizID)
	if err != nil {
		t.Fatalf("list questions after ensure: %v", err)
	}
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}
	want := nodeSupabaseQuestions()[6]
	for _, q := range qs {
		if q.Order != 6 {
			continue
		}
		if q.Text != want.text {
			t.Errorf("question 7 not restored: got %q", q.Text)
		}
		if q.ID != stale.ID {
			t.Errorf("question 7 should keep its id %q, got %q", stale.ID, q.ID)
		}
	}
}

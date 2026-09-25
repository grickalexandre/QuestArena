package seed

import (
	"testing"

	"github.com/questarena/questarena/internal/models"
)

func TestNodeSupabaseQuestions(t *testing.T) {
	qs := nodeSupabaseQuestions()
	if len(qs) != 30 {
		t.Fatalf("want 30 questions, got %d", len(qs))
	}
	withCode := 0
	for i, q := range qs {
		if q.text == "" {
			t.Errorf("q%d: empty text", i+1)
		}
		if len(q.options) != 4 {
			t.Errorf("q%d: want 4 options, got %d", i+1, len(q.options))
		}
		if q.correctIndex < 0 || q.correctIndex >= len(q.options) {
			t.Errorf("q%d: correctIndex %d out of range", i+1, q.correctIndex)
		}
		if q.timeLimitSec != timeLimitNodeSupabase {
			t.Errorf("q%d: want %ds, got %ds", i+1, timeLimitNodeSupabase, q.timeLimitSec)
		}
		for j, opt := range q.options {
			if opt == "" {
				t.Errorf("q%d option %d empty", i+1, j)
			}
		}
		if q.code == "" {
			if q.codeLanguage != "" {
				t.Errorf("q%d: language %q without code", i+1, q.codeLanguage)
			}
			continue
		}
		withCode++
		if models.NormalizeCodeLanguage(q.codeLanguage) != q.codeLanguage {
			t.Errorf("q%d: unsupported code language %q", i+1, q.codeLanguage)
		}
	}
	if withCode == 0 {
		t.Error("expected questions with code snippets")
	}
}

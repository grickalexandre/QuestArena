package seed

import "testing"

func TestHerancaQuestions(t *testing.T) {
	qs := herancaQuestions()
	if len(qs) != 15 {
		t.Fatalf("want 15 questions, got %d", len(qs))
	}
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
		for j, opt := range q.options {
			if opt == "" {
				t.Errorf("q%d option %d empty", i+1, j)
			}
		}
	}
}

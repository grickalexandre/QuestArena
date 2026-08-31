package game

import (
	"testing"

	"github.com/questarena/questarena/internal/models"
)

func TestSimilarityExact(t *testing.T) {
	if got := Similarity("Herança", "heranca"); got < 0.99 {
		t.Fatalf("expected near 1, got %v", got)
	}
}

func TestSimilarityPartial(t *testing.T) {
	got := Similarity("reaproveitamento de codigo com classes pai e filha", "reaproveitamento de código usando classe pai")
	if got < 0.4 {
		t.Fatalf("expected decent similarity, got %v", got)
	}
}

func TestSimilarityDifferent(t *testing.T) {
	got := Similarity("banana", "heranca em csharp")
	if got > 0.25 {
		t.Fatalf("expected low similarity, got %v", got)
	}
}

func TestSimilarityLongerParaphrase(t *testing.T) {
	got := Similarity(
		"Herança é o reaproveitamento de código de uma classe pai para uma classe filha",
		"reaproveitamento de código usando classe pai",
	)
	if got < 0.55 {
		t.Fatalf("longer correct explanation should score high, got %v", got)
	}
}

func TestBestSimilarityPicksClosest(t *testing.T) {
	got := BestSimilarity("polimorfismo", []string{
		"reaproveitamento de código",
		"capacidade de um objeto se comportar de várias formas",
		"polimorfismo",
	})
	if got < 0.99 {
		t.Fatalf("expected exact alternative to win, got %v", got)
	}
}

func TestGradeEssayUsesAlternatives(t *testing.T) {
	q := models.Question{
		ExpectedAnswer:  "reaproveitamento de código",
		ExpectedAnswers: []string{"herança", "classe pai e classe filha"},
	}
	got := GradeEssay("Herança", q)
	if got < 0.99 {
		t.Fatalf("alternative should match, got %v", got)
	}
}

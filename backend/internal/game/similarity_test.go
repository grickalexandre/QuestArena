package game

import "testing"

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

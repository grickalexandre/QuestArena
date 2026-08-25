package game

import "testing"

func TestCalcGradeAllCorrectInstant(t *testing.T) {
	got := calcGrade(15, 15, 15000, 15000)
	if got != 10 {
		t.Fatalf("want 10, got %v", got)
	}
}

func TestCalcGradeAllCorrectSlow(t *testing.T) {
	// XP mínimo com 15 acertos no limite do tempo: 20% de 15000 = 3000
	got := calcGrade(15, 15, 3000, 15000)
	if got != 7.6 {
		t.Fatalf("want 7.6, got %v", got)
	}
}

func TestCalcGradeTwelveInstant(t *testing.T) {
	got := calcGrade(12, 15, 12000, 15000)
	if got != 8 {
		t.Fatalf("want 8.0, got %v", got)
	}
}

func TestValidRA(t *testing.T) {
	if !validRA("123456") {
		t.Fatal("digits should be valid")
	}
	if validRA("12") {
		t.Fatal("too short")
	}
	if validRA("12 34") {
		t.Fatal("space should be invalid after normalize-only digits/letters")
	}
}

func TestNormalizeRA(t *testing.T) {
	if got := normalizeRA("  ab-12 "); got != "AB-12" {
		t.Fatalf("got %q", got)
	}
}

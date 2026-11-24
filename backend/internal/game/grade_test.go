package game

import "testing"

func TestQuizAwardsXP(t *testing.T) {
	if quizAwardsXP("seed-enade-revisao-teacher1") {
		t.Fatal("Quest 12 ENADE must not award XP")
	}
	if !quizAwardsXP("seed-gcs-teacher1") {
		t.Fatal("other quizzes should award XP")
	}
	if maxQuizXP("seed-enade-revisao-x", nil) != 0 {
		t.Fatal("ENADE max XP must be 0")
	}
}

func TestCalcGradeAllCorrect(t *testing.T) {
	got := calcGrade(15, 15, 15000, 15000)
	if got != 10 {
		t.Fatalf("want 10, got %v", got)
	}
}

func TestCalcGradeAllCorrectIgnoresXP(t *testing.T) {
	// Mesmo com XP baixo, a nota depende só dos acertos.
	got := calcGrade(15, 15, 3000, 15000)
	if got != 10 {
		t.Fatalf("want 10 (sem fator XP), got %v", got)
	}
}

func TestCalcGradeTwelveOfFifteen(t *testing.T) {
	got := calcGrade(12, 15, 12000, 15000)
	if got != 8 {
		t.Fatalf("want 8.0, got %v", got)
	}
}

func TestCalcGradeZero(t *testing.T) {
	got := calcGrade(0, 15, 0, 15000)
	if got != 0 {
		t.Fatalf("want 0, got %v", got)
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

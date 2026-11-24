package game

import (
	"math"
	"strings"
	"unicode"

	"github.com/questarena/questarena/internal/models"
)

const (
	gradeCorrectWeight = 10.0
)

func normalizeRA(raw string) string {
	return strings.ToUpper(strings.TrimSpace(raw))
}

func validRA(ra string) bool {
	n := len([]rune(ra))
	if n < 3 || n > 20 {
		return false
	}
	for _, r := range ra {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}
	return true
}

// quizAwardsXP: Quest 12 (ENADE) não dá XP — só acertos/nota.
func quizAwardsXP(quizID string) bool {
	return !strings.HasPrefix(quizID, "seed-enade-revisao-")
}

func maxQuizXP(quizID string, questions []models.Question) int {
	if !quizAwardsXP(quizID) {
		return 0
	}
	sum := 0.0
	for _, q := range questions {
		w := q.Weight
		if w <= 0 {
			w = 1
		}
		sum += w * 1000
	}
	return int(math.Round(sum))
}

// calcGrade: nota 0–10 só por acertos (sem fator de XP/velocidade).
func calcGrade(correctCount, total, score, maxScore int) float64 {
	_ = score
	_ = maxScore
	if total <= 0 {
		return 0
	}
	n := (float64(correctCount) / float64(total)) * gradeCorrectWeight
	if n < 0 {
		n = 0
	}
	if n > 10 {
		n = 10
	}
	return math.Round(n*10) / 10
}

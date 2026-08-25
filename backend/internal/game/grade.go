package game

import (
	"math"
	"strings"
	"unicode"

	"github.com/questarena/questarena/internal/models"
)

const (
	gradeCorrectWeight = 7.0
	gradeSpeedWeight   = 3.0
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

func maxQuizXP(questions []models.Question) int {
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

// calcGrade: 70% acertos + 30% XP (velocidade). Resultado 0–10 com 1 casa.
func calcGrade(correctCount, total, score, maxScore int) float64 {
	if total <= 0 {
		return 0
	}
	acerto := (float64(correctCount) / float64(total)) * gradeCorrectWeight
	vel := 0.0
	if maxScore > 0 {
		vel = (float64(score) / float64(maxScore)) * gradeSpeedWeight
	}
	n := acerto + vel
	if n < 0 {
		n = 0
	}
	if n > 10 {
		n = 10
	}
	return math.Round(n*10) / 10
}

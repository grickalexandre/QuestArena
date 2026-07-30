package game

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Similarity returns a score in [0,1] comparing student text to expected answer.
// Combines token Jaccard and character bigram Dice after normalization.
func Similarity(a, b string) float64 {
	na := normalizeText(a)
	nb := normalizeText(b)
	if na == "" && nb == "" {
		return 1
	}
	if na == "" || nb == "" {
		return 0
	}
	if na == nb {
		return 1
	}
	j := jaccardTokens(na, nb)
	d := diceBigrams(na, nb)
	// weigh tokens a bit more for short answers
	score := 0.55*j + 0.45*d
	if score > 1 {
		score = 1
	}
	return score
}

func normalizeText(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, s)
	if err != nil {
		out = s
	}
	var b strings.Builder
	prevSpace := false
	for _, r := range out {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevSpace = false
			continue
		}
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if !prevSpace && b.Len() > 0 {
				b.WriteByte(' ')
				prevSpace = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func jaccardTokens(a, b string) float64 {
	ta := uniqueTokens(a)
	tb := uniqueTokens(b)
	if len(ta) == 0 && len(tb) == 0 {
		return 1
	}
	inter := 0
	for t := range ta {
		if tb[t] {
			inter++
		}
	}
	union := len(ta) + len(tb) - inter
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

func uniqueTokens(s string) map[string]bool {
	out := make(map[string]bool)
	for _, t := range strings.Fields(s) {
		if t != "" {
			out[t] = true
		}
	}
	return out
}

func diceBigrams(a, b string) float64 {
	ba := bigrams(a)
	bb := bigrams(b)
	if len(ba) == 0 && len(bb) == 0 {
		return 1
	}
	if len(ba) == 0 || len(bb) == 0 {
		return 0
	}
	inter := 0
	for g, ca := range ba {
		if cb, ok := bb[g]; ok {
			if ca < cb {
				inter += ca
			} else {
				inter += cb
			}
		}
	}
	return 2 * float64(inter) / float64(lenCount(ba)+lenCount(bb))
}

func bigrams(s string) map[string]int {
	// remove spaces for character-level fuzzy match
	s = strings.ReplaceAll(s, " ", "")
	out := make(map[string]int)
	runes := []rune(s)
	if len(runes) < 2 {
		if len(runes) == 1 {
			out[string(runes)] = 1
		}
		return out
	}
	for i := 0; i < len(runes)-1; i++ {
		g := string(runes[i : i+2])
		out[g]++
	}
	return out
}

func lenCount(m map[string]int) int {
	n := 0
	for _, c := range m {
		n += c
	}
	return n
}

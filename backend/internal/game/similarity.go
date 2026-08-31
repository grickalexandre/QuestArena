package game

import (
	"strings"
	"unicode"

	"github.com/questarena/questarena/internal/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Similarity compares a student answer (a) to a reference (b) in [0,1].
// Token recall is weighted so a longer explanation that covers the key
// ideas scores higher than raw Jaccard (which punishes extra words).
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

	ts := contentTokens(na)
	te := contentTokens(nb)
	if len(te) == 0 {
		te = uniqueTokens(nb)
	}
	if len(ts) == 0 {
		ts = uniqueTokens(na)
	}

	_, recall, prec := tokenStats(ts, te)
	f1 := 0.0
	if prec+recall > 0 {
		f1 = 2 * prec * recall / (prec + recall)
	}

	dRecall := bigramRecall(na, nb)
	d := diceBigrams(na, nb)
	score := 0.40*recall + 0.25*f1 + 0.20*dRecall + 0.15*d
	if score > 1 {
		score = 1
	}
	if score < 0 {
		score = 0
	}
	return score
}

// BestSimilarity returns the highest Similarity against any reference.
func BestSimilarity(student string, refs []string) float64 {
	best := 0.0
	for _, ref := range refs {
		if s := Similarity(student, ref); s > best {
			best = s
		}
	}
	return best
}

// GradeEssay scores a student text against all accepted references on the question.
func GradeEssay(text string, q models.Question) float64 {
	return BestSimilarity(text, q.EssayReferences())
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

func contentTokens(s string) map[string]bool {
	out := make(map[string]bool)
	for _, t := range strings.Fields(s) {
		if t == "" || isStopword(t) {
			continue
		}
		out[stemToken(t)] = true
	}
	return out
}

func tokenStats(student, expected map[string]bool) (inter int, recall, prec float64) {
	for t := range expected {
		if student[t] {
			inter++
		}
	}
	if len(expected) > 0 {
		recall = float64(inter) / float64(len(expected))
	}
	if len(student) > 0 {
		prec = float64(inter) / float64(len(student))
	}
	return inter, recall, prec
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

func stemToken(t string) string {
	r := []rune(t)
	n := len(r)
	if n <= 4 {
		return t
	}
	if r[n-1] == 's' {
		return string(r[:n-1])
	}
	return t
}

func isStopword(t string) bool {
	_, ok := ptStopwords[t]
	return ok
}

var ptStopwords = map[string]bool{
	"a": true, "o": true, "as": true, "os": true,
	"um": true, "uma": true, "uns": true, "umas": true,
	"de": true, "da": true, "do": true, "das": true, "dos": true,
	"em": true, "na": true, "no": true, "nas": true, "nos": true,
	"para": true, "pra": true, "por": true, "com": true, "sem": true,
	"que": true, "se": true, "e": true, "ou": true, "mas": true,
	"como": true, "ao": true, "aos": true,
	"pelo": true, "pela": true, "pelos": true, "pelas": true,
	"ser": true, "foi": true, "sao": true, "tem": true, "ha": true,
	"este": true, "esta": true, "esse": true, "essa": true,
	"isto": true, "isso": true, "seu": true, "sua": true, "seus": true, "suas": true,
	"mais": true, "muito": true, "ja": true, "nao": true, "sim": true,
	"tambem": true, "entre": true, "sobre": true, "ate": true,
	"quando": true, "onde": true, "qual": true, "quais": true,
	"porque": true,
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

func bigramRecall(student, expected string) float64 {
	be := bigrams(expected)
	bs := bigrams(student)
	den := lenCount(be)
	if den == 0 {
		return 1
	}
	if lenCount(bs) == 0 {
		return 0
	}
	inter := 0
	for g, ce := range be {
		if cs, ok := bs[g]; ok {
			if cs < ce {
				inter += cs
			} else {
				inter += ce
			}
		}
	}
	return float64(inter) / float64(den)
}

func bigrams(s string) map[string]int {
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

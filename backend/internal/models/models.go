package models

import (
	"strings"
	"time"
)

type QuestionType string

const (
	QuestionMultipleChoice QuestionType = "multiple_choice"
	QuestionEssay          QuestionType = "essay"
)

type Teacher struct {
	ID        string    `json:"id" firestore:"id"`
	Email     string    `json:"email" firestore:"email"`
	Name      string    `json:"name" firestore:"name"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type Quiz struct {
	ID          string    `json:"id" firestore:"id"`
	TeacherID   string    `json:"teacherId" firestore:"teacherId"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type Question struct {
	ID                  string       `json:"id" firestore:"id"`
	QuizID              string       `json:"quizId" firestore:"quizId"`
	Type                QuestionType `json:"type" firestore:"type"`
	Text                string       `json:"text" firestore:"text"`
	Options             []string     `json:"options" firestore:"options"`
	CorrectIndex        int          `json:"correctIndex" firestore:"correctIndex"`
	ExpectedAnswer      string       `json:"expectedAnswer" firestore:"expectedAnswer"`
	ExpectedAnswers     []string     `json:"expectedAnswers,omitempty" firestore:"expectedAnswers,omitempty"`
	SimilarityThreshold float64      `json:"similarityThreshold" firestore:"similarityThreshold"`
	CodeSnippet         string       `json:"codeSnippet" firestore:"codeSnippet"`
	CodeLanguage        string       `json:"codeLanguage" firestore:"codeLanguage"`
	Weight              float64      `json:"weight" firestore:"weight"`
	TimeLimitSec        int          `json:"timeLimitSec" firestore:"timeLimitSec"`
	Order               int          `json:"order" firestore:"order"`
	CreatedAt           time.Time    `json:"createdAt" firestore:"createdAt"`
}

// PublicQuestion is sent to players (no correct / expected answer).
type PublicQuestion struct {
	ID           string       `json:"id"`
	Type         QuestionType `json:"type"`
	Text         string       `json:"text"`
	Options      []string     `json:"options,omitempty"`
	CodeSnippet  string       `json:"codeSnippet,omitempty"`
	CodeLanguage string       `json:"codeLanguage,omitempty"`
	Weight       float64      `json:"weight"`
	TimeLimitSec int          `json:"timeLimitSec"`
	Index        int          `json:"index"`
	Total        int          `json:"total"`
}

// CodeLanguages lists the syntax highlighting modes accepted for question snippets.
var CodeLanguages = []string{
	"plain", "java", "python", "javascript", "typescript", "csharp",
	"c", "cpp", "go", "php", "sql", "kotlin", "html", "css", "json", "shell",
}

// NormalizeCodeLanguage falls back to "plain" for unknown identifiers.
func NormalizeCodeLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "js":
		lang = "javascript"
	case "ts":
		lang = "typescript"
	case "c#", "cs":
		lang = "csharp"
	case "c++":
		lang = "cpp"
	case "py":
		lang = "python"
	case "bash", "sh":
		lang = "shell"
	}
	for _, l := range CodeLanguages {
		if l == lang {
			return lang
		}
	}
	return "plain"
}

type SessionRecord struct {
	ID         string         `json:"id" firestore:"id"`
	QuizID     string         `json:"quizId" firestore:"quizId"`
	QuizTitle  string         `json:"quizTitle" firestore:"quizTitle"`
	TeacherID  string         `json:"teacherId" firestore:"teacherId"`
	PIN        string         `json:"pin" firestore:"pin"`
	Ranking    []RankingEntry `json:"ranking" firestore:"ranking"`
	StartedAt  time.Time      `json:"startedAt" firestore:"startedAt"`
	FinishedAt time.Time      `json:"finishedAt" firestore:"finishedAt"`
}

type RankingEntry struct {
	PlayerID     string  `json:"playerId" firestore:"playerId"`
	Nickname     string  `json:"nickname" firestore:"nickname"`
	RA           string  `json:"ra" firestore:"ra"`
	Score        int     `json:"score" firestore:"score"`
	CorrectCount int     `json:"correctCount" firestore:"correctCount"`
	Total        int     `json:"total" firestore:"total"`
	MaxScore     int     `json:"maxScore" firestore:"maxScore"`
	Grade        float64 `json:"grade" firestore:"grade"`
	Rank         int     `json:"rank" firestore:"rank"`
}

// EssayReferences returns the primary expected answer plus unique alternatives.
func (q Question) EssayReferences() []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, 1+len(q.ExpectedAnswers))
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		key := strings.ToLower(s)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		out = append(out, s)
	}
	add(q.ExpectedAnswer)
	for _, a := range q.ExpectedAnswers {
		add(a)
	}
	return out
}

// NormalizeEssayRefs puts the first unique reference in ExpectedAnswer and the rest in ExpectedAnswers.
func (q *Question) NormalizeEssayRefs() {
	refs := q.EssayReferences()
	if len(refs) == 0 {
		q.ExpectedAnswer = ""
		q.ExpectedAnswers = nil
		return
	}
	q.ExpectedAnswer = refs[0]
	if len(refs) > 1 {
		q.ExpectedAnswers = refs[1:]
		return
	}
	q.ExpectedAnswers = nil
}

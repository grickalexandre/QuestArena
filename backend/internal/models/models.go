package models

import "time"

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
	SimilarityThreshold float64      `json:"similarityThreshold" firestore:"similarityThreshold"`
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
	Weight       float64      `json:"weight"`
	TimeLimitSec int          `json:"timeLimitSec"`
	Index        int          `json:"index"`
	Total        int          `json:"total"`
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
	PlayerID string `json:"playerId" firestore:"playerId"`
	Nickname string `json:"nickname" firestore:"nickname"`
	Score    int    `json:"score" firestore:"score"`
	Rank     int    `json:"rank" firestore:"rank"`
}

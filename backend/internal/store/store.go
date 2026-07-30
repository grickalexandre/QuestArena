package store

import (
	"context"

	"github.com/questarena/questarena/internal/models"
)

type Store interface {
	UpsertTeacher(ctx context.Context, t *models.Teacher) error
	GetTeacher(ctx context.Context, id string) (*models.Teacher, error)

	CreateQuiz(ctx context.Context, q *models.Quiz) error
	UpdateQuiz(ctx context.Context, q *models.Quiz) error
	DeleteQuiz(ctx context.Context, teacherID, quizID string) error
	GetQuiz(ctx context.Context, quizID string) (*models.Quiz, error)
	ListQuizzes(ctx context.Context, teacherID string) ([]models.Quiz, error)

	CreateQuestion(ctx context.Context, q *models.Question) error
	UpdateQuestion(ctx context.Context, q *models.Question) error
	DeleteQuestion(ctx context.Context, quizID, questionID string) error
	ListQuestions(ctx context.Context, quizID string) ([]models.Question, error)

	SaveSession(ctx context.Context, s *models.SessionRecord) error
}

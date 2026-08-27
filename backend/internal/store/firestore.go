package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"github.com/questarena/questarena/internal/models"
	"google.golang.org/api/iterator"
)

type FirestoreStore struct {
	client *firestore.Client
}

func NewFirestoreStore(ctx context.Context, app *firebase.App) (*FirestoreStore, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return &FirestoreStore{client: client}, nil
}

func (f *FirestoreStore) UpsertTeacher(ctx context.Context, t *models.Teacher) error {
	_, err := f.client.Collection("teachers").Doc(t.ID).Set(ctx, t)
	return err
}

func (f *FirestoreStore) GetTeacher(ctx context.Context, id string) (*models.Teacher, error) {
	doc, err := f.client.Collection("teachers").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var t models.Teacher
	if err := doc.DataTo(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (f *FirestoreStore) CreateQuiz(ctx context.Context, q *models.Quiz) error {
	if q.ID == "" {
		q.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	q.CreatedAt = now
	q.UpdatedAt = now
	_, err := f.client.Collection("quizzes").Doc(q.ID).Set(ctx, q)
	return err
}

func (f *FirestoreStore) UpdateQuiz(ctx context.Context, q *models.Quiz) error {
	doc, err := f.client.Collection("quizzes").Doc(q.ID).Get(ctx)
	if err != nil {
		return err
	}
	var existing models.Quiz
	if err := doc.DataTo(&existing); err != nil {
		return err
	}
	if existing.TeacherID != q.TeacherID {
		return fmt.Errorf("forbidden")
	}
	q.CreatedAt = existing.CreatedAt
	q.UpdatedAt = time.Now().UTC()
	_, err = f.client.Collection("quizzes").Doc(q.ID).Set(ctx, q)
	return err
}

func (f *FirestoreStore) DeleteQuiz(ctx context.Context, teacherID, quizID string) error {
	doc, err := f.client.Collection("quizzes").Doc(quizID).Get(ctx)
	if err != nil {
		return err
	}
	var q models.Quiz
	if err := doc.DataTo(&q); err != nil {
		return err
	}
	if q.TeacherID != teacherID {
		return fmt.Errorf("forbidden")
	}
	qs, err := f.client.Collection("quizzes").Doc(quizID).Collection("questions").Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	batch := f.client.Batch()
	for _, d := range qs {
		batch.Delete(d.Ref)
	}
	batch.Delete(f.client.Collection("quizzes").Doc(quizID))
	_, err = batch.Commit(ctx)
	return err
}

func (f *FirestoreStore) GetQuiz(ctx context.Context, quizID string) (*models.Quiz, error) {
	doc, err := f.client.Collection("quizzes").Doc(quizID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var q models.Quiz
	if err := doc.DataTo(&q); err != nil {
		return nil, err
	}
	return &q, nil
}

func (f *FirestoreStore) ListQuizzes(ctx context.Context, teacherID string) ([]models.Quiz, error) {
	iter := f.client.Collection("quizzes").Where("teacherId", "==", teacherID).Documents(ctx)
	out := make([]models.Quiz, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var q models.Quiz
		if err := doc.DataTo(&q); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].UpdatedAt.After(out[j].UpdatedAt)
	})
	return out, nil
}

func (f *FirestoreStore) CreateQuestion(ctx context.Context, q *models.Question) error {
	if q.ID == "" {
		q.ID = uuid.NewString()
	}
	q.CreatedAt = time.Now().UTC()
	_, err := f.client.Collection("quizzes").Doc(q.QuizID).Collection("questions").Doc(q.ID).Set(ctx, q)
	return err
}

func (f *FirestoreStore) UpdateQuestion(ctx context.Context, q *models.Question) error {
	ref := f.client.Collection("quizzes").Doc(q.QuizID).Collection("questions").Doc(q.ID)
	doc, err := ref.Get(ctx)
	if err != nil {
		return err
	}
	var existing models.Question
	if err := doc.DataTo(&existing); err != nil {
		return err
	}
	q.CreatedAt = existing.CreatedAt
	_, err = ref.Set(ctx, q)
	return err
}

func (f *FirestoreStore) DeleteQuestion(ctx context.Context, quizID, questionID string) error {
	_, err := f.client.Collection("quizzes").Doc(quizID).Collection("questions").Doc(questionID).Delete(ctx)
	return err
}

func (f *FirestoreStore) ListQuestions(ctx context.Context, quizID string) ([]models.Question, error) {
	iter := f.client.Collection("quizzes").Doc(quizID).Collection("questions").Documents(ctx)
	out := make([]models.Question, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var q models.Question
		if err := doc.DataTo(&q); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Order < out[j].Order
	})
	return out, nil
}

func (f *FirestoreStore) SaveSession(ctx context.Context, s *models.SessionRecord) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	_, err := f.client.Collection("sessions").Doc(s.ID).Set(ctx, s)
	return err
}

func (f *FirestoreStore) ListSessions(ctx context.Context, teacherID string) ([]models.SessionRecord, error) {
	iter := f.client.Collection("sessions").Where("teacherId", "==", teacherID).Documents(ctx)
	out := make([]models.SessionRecord, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var s models.SessionRecord
		if err := doc.DataTo(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].FinishedAt.After(out[j].FinishedAt)
	})
	return out, nil
}

func (f *FirestoreStore) GetSession(ctx context.Context, teacherID, sessionID string) (*models.SessionRecord, error) {
	doc, err := f.client.Collection("sessions").Doc(sessionID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var s models.SessionRecord
	if err := doc.DataTo(&s); err != nil {
		return nil, err
	}
	if s.TeacherID != teacherID {
		return nil, fmt.Errorf("forbidden")
	}
	return &s, nil
}

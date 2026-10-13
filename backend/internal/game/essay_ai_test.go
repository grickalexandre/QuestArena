package game

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/questarena/questarena/internal/grader"
	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

type fakeGrader struct {
	score float64
	err   error
	got   chan grader.Batch
}

func (f *fakeGrader) Grade(_ context.Context, b grader.Batch) ([]grader.Verdict, error) {
	f.got <- b
	if f.err != nil {
		return nil, f.err
	}
	out := make([]grader.Verdict, 0, len(b.Answers))
	for _, a := range b.Answers {
		out = append(out, grader.Verdict{ID: a.ID, Score: f.score, Feedback: "Explicou a ideia certa."})
	}
	return out, nil
}

func essayRoom(t *testing.T, g EssayGrader) (*Hub, *Room, *Client) {
	t.Helper()
	h := NewHub(store.NewMemoryStore())
	if g != nil {
		h.SetGrader(g)
	}
	room, err := h.CreateRoom(models.Quiz{ID: "q1", Title: "Quiz"}, []models.Question{{
		ID: "1", Type: models.QuestionEssay, Text: "O que é especialização?",
		ExpectedAnswer:      "Uma entidade filha herda a chave da entidade pai; é um tipo de.",
		GradingCriteria:     "Citar que é relação 'é um tipo de'.",
		SimilarityThreshold: 0.55, Weight: 1, TimeLimitSec: 60,
	}}, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{hub: h, send: make(chan []byte, 64)}
	join, _ := json.Marshal(map[string]any{"pin": room.PIN, "nickname": "Ana", "ra": "RA123", "avatar": 1})
	c.handlePlayerJoin(join)
	room.mu.Lock()
	room.startQuestionLocked()
	room.mu.Unlock()
	ans, _ := json.Marshal(map[string]any{"text": "Aluno é um caso particular de pessoa, herda tudo dela."})
	c.handleAnswer(ans)
	return h, room, c
}

func waitPending(t *testing.T, room *Room, pid string) answerRecord {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		room.mu.Lock()
		rec := room.History[0][pid]
		room.mu.Unlock()
		if !rec.pending {
			return rec
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("ai grading did not finish")
	return answerRecord{}
}

func TestAIGradingReplacesLexicalScore(t *testing.T) {
	g := &fakeGrader{score: 0.9, got: make(chan grader.Batch, 1)}
	_, room, c := essayRoom(t, g)
	pid := c.playerID

	room.mu.Lock()
	room.revealLocked()
	room.mu.Unlock()

	b := <-g.got
	if len(b.Answers) != 1 || b.Criteria == "" || strings.Contains(b.Answers[0].ID, pid) {
		t.Fatalf("batch must carry criteria and anonymous ids: %+v", b)
	}
	rec := waitPending(t, room, pid)
	if rec.method != gradeMethodAI || rec.similarity != 0.9 || !rec.correct {
		t.Fatalf("ai verdict not applied: %+v", rec)
	}
	want := int(float64(rec.basePoints)*0.9 + 0.5)
	room.mu.Lock()
	score := room.Players[pid].Score
	correct := room.Players[pid].CorrectCount
	room.mu.Unlock()
	if rec.points != want || score != want || correct != 1 {
		t.Fatalf("points=%d score=%d correct=%d want %d", rec.points, score, correct, want)
	}
}

func TestAIGradingFailureKeepsLexicalScore(t *testing.T) {
	g := &fakeGrader{err: errors.New("quota"), got: make(chan grader.Batch, 1)}
	_, room, c := essayRoom(t, g)
	pid := c.playerID
	room.mu.Lock()
	before := room.History[0][pid]
	room.revealLocked()
	room.mu.Unlock()
	<-g.got
	rec := waitPending(t, room, pid)
	if rec.method != gradeMethodText || rec.points != before.points {
		t.Fatalf("fallback changed score: before=%+v after=%+v", before, rec)
	}
}

func TestTeacherOverrideAdjustsTotals(t *testing.T) {
	h, room, c := essayRoom(t, nil)
	pid := c.playerID
	room.mu.Lock()
	room.revealLocked()
	room.mu.Unlock()

	host := &Client{hub: h, send: make(chan []byte, 64)}
	hj, _ := json.Marshal(map[string]any{"pin": room.PIN, "teacherId": "teacher-1"})
	host.handleHostJoin(hj)

	ov, _ := json.Marshal(map[string]any{"questionIndex": 0, "playerId": pid, "score": 100})
	host.handleGradeOverride(ov)
	room.mu.Lock()
	rec := room.History[0][pid]
	p := room.Players[pid]
	paused := room.autoPaused
	room.mu.Unlock()
	if rec.method != gradeMethodTeacher || rec.points != rec.basePoints || p.Score != rec.basePoints || p.CorrectCount != 1 || !paused {
		t.Fatalf("override 100: rec=%+v score=%d correct=%d paused=%v", rec, p.Score, p.CorrectCount, paused)
	}

	ov, _ = json.Marshal(map[string]any{"questionIndex": 0, "playerId": pid, "score": 0})
	host.handleGradeOverride(ov)
	if p.Score != 0 || p.CorrectCount != 0 {
		t.Fatalf("override 0: score=%d correct=%d", p.Score, p.CorrectCount)
	}

	player, _ := json.Marshal(map[string]any{"questionIndex": 0, "playerId": pid, "score": 100})
	c.handleGradeOverride(player)
	if p.Score != 0 {
		t.Fatal("player must not override grades")
	}
}

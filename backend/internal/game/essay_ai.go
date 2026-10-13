package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/questarena/questarena/internal/grader"
	"github.com/questarena/questarena/internal/models"
)

const (
	gradeMethodText    = "texto"
	gradeMethodAI      = "ia"
	gradeMethodTeacher = "professor"

	aiGradingTimeout = 90 * time.Second
)

// EssayGrader scores all answers to one essay question in a single call.
type EssayGrader interface {
	Grade(ctx context.Context, b grader.Batch) ([]grader.Verdict, error)
}

// SetGrader enables context-based essay grading. With nil, essays use lexical grading only.
func (h *Hub) SetGrader(g EssayGrader) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.grader = g
}

func (h *Hub) essayGrader() EssayGrader {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.grader
}

// AIGradingEnabled reports whether essays are graded by context.
func (h *Hub) AIGradingEnabled() bool {
	return h.essayGrader() != nil
}

func essayThreshold(q models.Question) float64 {
	if q.SimilarityThreshold <= 0 {
		return 0.55
	}
	return q.SimilarityThreshold
}

// startAIGradingLocked marks the current essay answers as pending and grades them in the background.
// It returns false when there is nothing to grade by context.
func (r *Room) startAIGradingLocked() bool {
	q := r.Questions[r.CurrentIndex]
	if q.Type != models.QuestionEssay || r.hub == nil {
		return false
	}
	g := r.hub.essayGrader()
	if g == nil {
		return false
	}
	idx := r.CurrentIndex
	answers := r.History[idx]
	ids := map[string]string{}
	batch := grader.Batch{
		Question: q.Text,
		Expected: q.EssayReferences(),
		Criteria: q.GradingCriteria,
		KeyTerms: q.KeyTerms,
	}
	n := 0
	for pid, rec := range answers {
		if rec.text == "" || rec.method == gradeMethodTeacher {
			continue
		}
		n++
		anon := fmt.Sprintf("r%d", n)
		ids[anon] = pid
		batch.Answers = append(batch.Answers, grader.Answer{ID: anon, Text: rec.text})
		rec.pending = true
		answers[pid] = rec
	}
	if n == 0 {
		return false
	}
	go r.runAIGrading(g, idx, batch, ids)
	return true
}

func (r *Room) runAIGrading(g EssayGrader, idx int, batch grader.Batch, ids map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), aiGradingTimeout)
	defer cancel()
	started := time.Now()
	verdicts, err := g.Grade(ctx, batch)
	if err != nil {
		log.Printf("ai grading pin=%s q=%d: %v (mantida a correção por texto)", r.PIN, idx+1, err)
	} else {
		log.Printf("ai grading pin=%s q=%d answers=%d took=%s", r.PIN, idx+1, len(verdicts), time.Since(started).Round(time.Millisecond))
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, v := range verdicts {
		pid, ok := ids[v.ID]
		if !ok {
			continue
		}
		rec, ok := r.History[idx][pid]
		if !ok || !rec.pending {
			continue
		}
		r.applyEssayScoreLocked(idx, pid, v.Score, v.Feedback, gradeMethodAI)
	}
	for _, pid := range ids {
		if rec, ok := r.History[idx][pid]; ok && rec.pending {
			rec.pending = false
			r.History[idx][pid] = rec
		}
	}
	r.afterRegradeLocked(idx, true)
}

// applyEssayScoreLocked replaces the score of one essay answer and adjusts the player's totals.
func (r *Room) applyEssayScoreLocked(idx int, playerID string, score float64, feedback, method string) bool {
	answers := r.History[idx]
	rec, ok := answers[playerID]
	if !ok || rec.text == "" {
		return false
	}
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	newPoints := int(math.Round(float64(rec.basePoints) * score))
	newCorrect := score >= essayThreshold(r.Questions[idx])
	if p, ok := r.Players[playerID]; ok {
		p.Score += newPoints - rec.points
		if newCorrect && !rec.correct {
			p.CorrectCount++
		} else if !newCorrect && rec.correct {
			p.CorrectCount--
		}
	}
	rec.points = newPoints
	rec.correct = newCorrect
	rec.similarity = score
	rec.feedback = feedback
	rec.method = method
	rec.pending = false
	answers[playerID] = rec
	return true
}

// afterRegradeLocked pushes updated scores to everyone. When the AI finishes on the
// question still being shown, it also starts the countdown to the next question.
func (r *Room) afterRegradeLocked(idx int, fromAI bool) {
	switch {
	case r.Phase == PhaseReveal && r.CurrentIndex == idx:
		payload := r.buildResultPayloadLocked()
		hold := 0
		if fromAI && !r.autoPaused {
			hold = revealHoldSec(models.QuestionEssay)
		}
		payload["autoNextInSec"] = hold
		r.broadcastLocked("question_result", payload)
		if hold > 0 {
			r.scheduleAutoNextLocked(hold)
		}
	case r.Phase == PhaseFinished:
		r.broadcastFinishedLocked()
		r.saveSessionLocked()
	default:
		r.broadcastLocked("scores_updated", map[string]any{
			"leaderboard": leaderboard(r),
			"players":     publicPlayers(r),
		})
	}
}

func (c *Client) handleGradeOverride(data json.RawMessage) {
	var req struct {
		QuestionIndex int     `json:"questionIndex"`
		PlayerID      string  `json:"playerId"`
		Score         float64 `json:"score"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		c.sendError("invalid grade_override")
		return
	}
	room := c.hub.GetRoom(c.roomPIN)
	if room == nil || c.role != RoleHost || room.Host != c {
		c.sendError("você não é o host desta sala — reconectando...")
		return
	}
	if req.Score < 0 || req.Score > 100 {
		c.sendError("a nota deve ficar entre 0 e 100")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if req.QuestionIndex < 0 || req.QuestionIndex >= len(room.Questions) {
		c.sendError("questão inválida")
		return
	}
	if room.Questions[req.QuestionIndex].Type != models.QuestionEssay {
		c.sendError("só dá para ajustar nota de dissertativa")
		return
	}
	if room.Phase == PhaseQuestion && room.CurrentIndex == req.QuestionIndex {
		c.sendError("aguarde o fim do tempo para ajustar a nota")
		return
	}
	rec, ok := room.History[req.QuestionIndex][req.PlayerID]
	if !ok || rec.text == "" {
		c.sendError("esse aluno não respondeu esta questão")
		return
	}
	if !room.applyEssayScoreLocked(req.QuestionIndex, req.PlayerID, req.Score/100, rec.feedback, gradeMethodTeacher) {
		c.sendError("não foi possível ajustar a nota")
		return
	}
	if room.Phase == PhaseReveal && room.CurrentIndex == req.QuestionIndex {
		room.autoPaused = true
		room.cancelAdvanceLocked()
	}
	room.afterRegradeLocked(req.QuestionIndex, false)
}

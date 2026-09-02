package game

import (
	"encoding/json"
	"testing"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

func TestPlayerRejoinKeepsSeatAndScore(t *testing.T) {
	h := NewHub(store.NewMemoryStore())
	room, err := h.CreateRoom(models.Quiz{ID: "q1", Title: "Quiz"}, nil, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}

	c1 := &Client{hub: h, send: make(chan []byte, 8)}
	joinPayload, _ := json.Marshal(map[string]any{
		"pin": room.PIN, "nickname": "Ana", "ra": "RA123", "avatar": 2,
	})
	c1.handlePlayerJoin(joinPayload)
	if len(room.Players) != 1 {
		t.Fatalf("want 1 player, got %d", len(room.Players))
	}
	var player *Player
	for _, p := range room.Players {
		player = p
	}
	player.Score = 900
	player.CorrectCount = 2

	c1.cleanup()
	if _, ok := room.Players[player.ID]; !ok {
		t.Fatal("disconnected player was removed")
	}
	if player.conn != nil {
		t.Fatal("expected player.conn nil after disconnect")
	}

	c2 := &Client{hub: h, send: make(chan []byte, 8)}
	c2.handlePlayerJoin(joinPayload)
	if len(room.Players) != 1 {
		t.Fatalf("rejoin created a new seat, got %d players", len(room.Players))
	}
	if c2.playerID != player.ID {
		t.Fatalf("rejoin assigned new id: %s vs %s", c2.playerID, player.ID)
	}
	if player.Score != 900 || player.CorrectCount != 2 {
		t.Fatalf("score lost on rejoin: %+v", player)
	}
	if player.conn != c2 {
		t.Fatal("rejoin did not reattach connection")
	}
}

func TestMemoryStoreListsSessionsByTeacher(t *testing.T) {
	mem := store.NewMemoryStore()
	rec := &models.SessionRecord{
		ID:        "s1",
		TeacherID: "t1",
		QuizTitle: "Herança",
		PIN:       "123456",
		Ranking: []models.RankingEntry{
			{PlayerID: "p1", Nickname: "Ana", RA: "RA123", Grade: 8.5, Rank: 1},
		},
	}
	if err := mem.SaveSession(nil, rec); err != nil {
		t.Fatal(err)
	}
	list, err := mem.ListSessions(nil, "t1")
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v n=%d", err, len(list))
	}
	got, err := mem.GetSession(nil, "t1", "s1")
	if err != nil || got.Ranking[0].Grade != 8.5 {
		t.Fatalf("get: %v %+v", err, got)
	}
	if _, err := mem.GetSession(nil, "other", "s1"); err == nil {
		t.Fatal("expected forbidden")
	}
}

func TestPresenceCountsAwayOnlyOnQuestion(t *testing.T) {
	h := NewHub(store.NewMemoryStore())
	room, err := h.CreateRoom(models.Quiz{ID: "q1", Title: "Quiz"}, []models.Question{{
		ID: "1", Text: "Q", Options: []string{"A", "B"}, CorrectIndex: 0, Weight: 1000, TimeLimitSec: 30,
	}}, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}

	c := &Client{hub: h, send: make(chan []byte, 16)}
	joinPayload, _ := json.Marshal(map[string]any{
		"pin": room.PIN, "nickname": "Ana", "ra": "RA123", "avatar": 1,
	})
	c.handlePlayerJoin(joinPayload)

	hidden, _ := json.Marshal(map[string]any{"hidden": true})
	visible, _ := json.Marshal(map[string]any{"hidden": false})
	c.handlePresence(hidden)

	var player *Player
	for _, p := range room.Players {
		player = p
	}
	if !player.Hidden || player.AwayCount != 0 || player.AwayTotal != 0 {
		t.Fatalf("lobby should not count away: hidden=%v count=%d total=%d", player.Hidden, player.AwayCount, player.AwayTotal)
	}

	room.mu.Lock()
	room.startQuestionLocked()
	room.mu.Unlock()
	if player.AwayCount != 1 || player.AwayTotal != 1 {
		t.Fatalf("already hidden at start: count=%d total=%d", player.AwayCount, player.AwayTotal)
	}

	c.handlePresence(visible)
	c.handlePresence(hidden)
	if player.AwayCount != 2 || player.AwayTotal != 2 {
		t.Fatalf("leave during question: count=%d total=%d", player.AwayCount, player.AwayTotal)
	}
	c.handlePresence(hidden)
	if player.AwayCount != 2 {
		t.Fatalf("duplicate hidden incremented count: %d", player.AwayCount)
	}
	c.handlePresence(visible)
	c.handlePresence(hidden)
	if player.AwayCount != 3 || player.AwayTotal != 3 || !player.Hidden {
		t.Fatalf("second leave: hidden=%v count=%d total=%d", player.Hidden, player.AwayCount, player.AwayTotal)
	}

	room.mu.Lock()
	room.startQuestionLocked()
	room.mu.Unlock()
	if player.AwayCount != 1 || player.AwayTotal != 4 {
		t.Fatalf("new question keeps hidden flag: count=%d total=%d", player.AwayCount, player.AwayTotal)
	}
}

func TestPlayerCannotChangeAnswerAfterSubmit(t *testing.T) {
	h := NewHub(store.NewMemoryStore())
	room, err := h.CreateRoom(models.Quiz{ID: "q1", Title: "Quiz"}, []models.Question{{
		ID: "1", Text: "Q", Options: []string{"A", "B", "C"}, CorrectIndex: 1, Weight: 1000, TimeLimitSec: 30,
	}}, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}

	c := &Client{hub: h, send: make(chan []byte, 32)}
	joinPayload, _ := json.Marshal(map[string]any{
		"pin": room.PIN, "nickname": "Ana", "ra": "RA123", "avatar": 1,
	})
	c.handlePlayerJoin(joinPayload)

	room.mu.Lock()
	room.startQuestionLocked()
	room.mu.Unlock()

	wrong, _ := json.Marshal(map[string]any{"choice": 0})
	c.handleAnswer(wrong)
	p := room.Players[c.playerID]
	if !p.answered || p.choice != 0 || p.CorrectCount != 0 || p.Score != 0 {
		t.Fatalf("first wrong answer: %+v ans=%+v", p, room.Answers[p.ID])
	}

	right, _ := json.Marshal(map[string]any{"choice": 1})
	c.handleAnswer(right)
	if p.choice != 0 || p.CorrectCount != 0 || p.Score != 0 || room.Answers[p.ID].correct {
		t.Fatalf("allowed change after lock: choice=%d score=%d ans=%+v", p.choice, p.Score, room.Answers[p.ID])
	}

	ack := lastEnvelopeOfType(t, c, "answer_ack")
	if _, ok := ack["points"]; ok {
		t.Fatalf("points leaked in answer_ack: %+v", ack)
	}
	if _, ok := ack["similarity"]; ok {
		t.Fatalf("similarity leaked in answer_ack: %+v", ack)
	}
	if ack["locked"] != true {
		t.Fatalf("expected locked ack: %+v", ack)
	}
}

func TestPresenceReportsInspect(t *testing.T) {
	h := NewHub(store.NewMemoryStore())
	room, err := h.CreateRoom(models.Quiz{ID: "q1", Title: "Quiz"}, []models.Question{{
		ID: "1", Text: "Q", Options: []string{"A", "B"}, CorrectIndex: 0, Weight: 1000, TimeLimitSec: 30,
	}}, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}

	c := &Client{hub: h, send: make(chan []byte, 16)}
	joinPayload, _ := json.Marshal(map[string]any{
		"pin": room.PIN, "nickname": "Ana", "ra": "RA123", "avatar": 1,
	})
	c.handlePlayerJoin(joinPayload)

	inspectOn, _ := json.Marshal(map[string]any{"inspect": true})
	c.handlePresence(inspectOn)

	var player *Player
	for _, p := range room.Players {
		player = p
	}
	if !player.Inspecting || player.InspectCount != 0 {
		t.Fatalf("lobby should not count inspect: inspecting=%v count=%d", player.Inspecting, player.InspectCount)
	}

	room.mu.Lock()
	room.startQuestionLocked()
	room.mu.Unlock()
	if player.InspectCount != 1 {
		t.Fatalf("already inspecting at start: count=%d", player.InspectCount)
	}

	inspectOff, _ := json.Marshal(map[string]any{"inspect": false})
	c.handlePresence(inspectOff)
	c.handlePresence(inspectOn)
	if !player.Inspecting || player.InspectCount != 2 {
		t.Fatalf("inspect during question: inspecting=%v count=%d", player.Inspecting, player.InspectCount)
	}
}

func lastEnvelopeOfType(t *testing.T, c *Client, typ string) map[string]any {
	t.Helper()
	var last map[string]any
	for {
		select {
		case raw := <-c.send:
			var env Envelope
			if err := json.Unmarshal(raw, &env); err != nil {
				t.Fatal(err)
			}
			if env.Type != typ {
				continue
			}
			var data map[string]any
			if err := json.Unmarshal(env.Data, &data); err != nil {
				t.Fatal(err)
			}
			last = data
		default:
			if last == nil {
				t.Fatalf("no %s message", typ)
			}
			return last
		}
	}
}

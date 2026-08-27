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

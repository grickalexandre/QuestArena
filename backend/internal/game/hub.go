package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Role string

const (
	RoleHost   Role = "host"
	RolePlayer Role = "player"
)

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhaseQuestion Phase = "question"
	PhaseReveal   Phase = "reveal"
	PhaseFinished Phase = "finished"
)

type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	role     Role
	roomPIN  string
	playerID string
	nickname string
	hostID   string
}

type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Score    int    `json:"score"`
	Avatar   int    `json:"avatar"`
	conn     *Client
	answered bool
	choice   int
}

type answerRecord struct {
	playerID   string
	choice     int
	text       string
	elapsed    time.Duration
	correct    bool
	points     int
	similarity float64
}

type Room struct {
	PIN          string
	SessionID    string
	Quiz         models.Quiz
	Questions    []models.Question
	TeacherID    string
	Phase        Phase
	CurrentIndex int
	QuestionAt   time.Time
	Deadline     time.Time
	Host         *Client
	Players      map[string]*Player
	Answers      map[string]answerRecord
	timerCancel  context.CancelFunc
	advanceCancel context.CancelFunc
	mu           sync.Mutex
	hub          *Hub
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	store store.Store
}

func NewHub(s store.Store) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		store: s,
	}
}

func (h *Hub) CreateRoom(quiz models.Quiz, questions []models.Question, teacherID string) (*Room, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var pin string
	for i := 0; i < 30; i++ {
		pin = fmt.Sprintf("%06d", rand.Intn(900000)+100000)
		if _, exists := h.rooms[pin]; !exists {
			break
		}
	}
	room := &Room{
		PIN:       pin,
		SessionID: uuid.NewString(),
		Quiz:      quiz,
		Questions: questions,
		TeacherID: teacherID,
		Phase:     PhaseLobby,
		Players:   make(map[string]*Player),
		Answers:   make(map[string]answerRecord),
		hub:       h,
	}
	h.rooms[pin] = room
	return room, nil
}

func (h *Hub) GetRoom(pin string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[pin]
}

func (h *Hub) RemoveRoom(pin string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.rooms[pin]; ok {
		room.mu.Lock()
		room.cancelTimersLocked()
		room.mu.Unlock()
		delete(h.rooms, pin)
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade: %v", err)
		return
	}
	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.cleanup()
		c.conn.Close()
	}()
	c.conn.SetReadLimit(8192)
	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var env Envelope
		if err := json.Unmarshal(msg, &env); err != nil {
			c.sendError("invalid message")
			continue
		}
		c.handle(env)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(20 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) emit(typ string, data any) {
	payload, _ := json.Marshal(data)
	env, _ := json.Marshal(Envelope{Type: typ, Data: payload})
	select {
	case c.send <- env:
	default:
	}
}

func (c *Client) sendError(msg string) {
	c.emit("error", map[string]string{"message": msg})
}

func (c *Client) handle(env Envelope) {
	switch env.Type {
	case "host_join":
		c.handleHostJoin(env.Data)
	case "join":
		c.handlePlayerJoin(env.Data)
	case "start":
		c.handleStart()
	case "next":
		c.handleNext()
	case "answer":
		c.handleAnswer(env.Data)
	case "end":
		c.handleEnd()
	default:
		c.sendError("unknown event: " + env.Type)
	}
}

func (c *Client) handleHostJoin(data json.RawMessage) {
	var req struct {
		PIN       string `json:"pin"`
		TeacherID string `json:"teacherId"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		c.sendError("invalid host_join")
		return
	}
	room := c.hub.GetRoom(req.PIN)
	if room == nil {
		c.sendError("room not found")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.TeacherID != req.TeacherID {
		c.sendError("forbidden")
		return
	}
	c.role = RoleHost
	c.roomPIN = room.PIN
	c.hostID = req.TeacherID
	room.Host = c
	payload := map[string]any{
		"pin":           room.PIN,
		"quizTitle":     room.Quiz.Title,
		"phase":         room.Phase,
		"players":       publicPlayers(room),
		"questionCount": len(room.Questions),
	}
	if room.Phase == PhaseQuestion || room.Phase == PhaseReveal {
		payload["question"] = publicQuestion(room)
		payload["endsAt"] = room.Deadline.UTC().Format(time.RFC3339Nano)
		if room.Phase == PhaseReveal {
			payload["questionResult"] = room.buildResultPayloadLocked()
		}
	}
	c.emit("hosted", payload)
}

func (c *Client) handlePlayerJoin(data json.RawMessage) {
	var req struct {
		PIN      string `json:"pin"`
		Nickname string `json:"nickname"`
		Avatar   int    `json:"avatar"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		c.sendError("invalid join")
		return
	}
	if len(req.Nickname) < 2 || len(req.Nickname) > 20 {
		c.sendError("nickname must be 2-20 characters")
		return
	}
	const avatarCount = 12
	avatar := req.Avatar
	if avatar < 0 || avatar >= avatarCount {
		avatar = rand.Intn(avatarCount)
	}
	room := c.hub.GetRoom(req.PIN)
	if room == nil {
		c.sendError("PIN inválido")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.Phase != PhaseLobby {
		c.sendError("a partida já começou")
		return
	}
	for _, p := range room.Players {
		if p.Nickname == req.Nickname {
			c.sendError("nickname já em uso")
			return
		}
	}
	player := &Player{
		ID:       uuid.NewString(),
		Nickname: req.Nickname,
		Avatar:   avatar,
		conn:     c,
		choice:   -1,
	}
	room.Players[player.ID] = player
	c.role = RolePlayer
	c.roomPIN = room.PIN
	c.playerID = player.ID
	c.nickname = player.Nickname

	c.emit("joined", map[string]any{
		"playerId":  player.ID,
		"nickname":  player.Nickname,
		"avatar":    player.Avatar,
		"pin":       room.PIN,
		"quizTitle": room.Quiz.Title,
		"phase":     room.Phase,
	})
	room.broadcastLocked("player_joined", map[string]any{
		"players": publicPlayers(room),
	})
}

func (c *Client) handleStart() {
	room := c.hub.GetRoom(c.roomPIN)
	if room == nil || c.role != RoleHost || room.Host != c {
		c.sendError("você não é o host desta sala — reconectando...")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.Phase != PhaseLobby {
		c.sendError("already started")
		return
	}
	if len(room.Questions) == 0 {
		c.sendError("quiz has no questions")
		return
	}
	if len(room.Players) == 0 {
		c.sendError("aguarde pelo menos 1 jogador")
		return
	}
	room.CurrentIndex = 0
	room.startQuestionLocked()
}

func (c *Client) handleNext() {
	room := c.hub.GetRoom(c.roomPIN)
	if room == nil || c.role != RoleHost || room.Host != c {
		c.sendError("você não é o host desta sala — reconectando...")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.Phase != PhaseReveal {
		c.sendError("cannot advance now")
		return
	}
	room.cancelAdvanceLocked()
	room.CurrentIndex++
	if room.CurrentIndex >= len(room.Questions) {
		room.finishLocked()
		return
	}
	room.startQuestionLocked()
}

func (c *Client) handleEnd() {
	room := c.hub.GetRoom(c.roomPIN)
	if room == nil || c.role != RoleHost || room.Host != c {
		c.sendError("você não é o host desta sala — reconectando...")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	room.cancelAdvanceLocked()
	room.finishLocked()
}

func (c *Client) handleAnswer(data json.RawMessage) {
	var req struct {
		Choice int    `json:"choice"`
		Text   string `json:"text"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		c.sendError("invalid answer")
		return
	}
	room := c.hub.GetRoom(c.roomPIN)
	if room == nil || c.role != RolePlayer {
		c.sendError("not a player")
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if room.Phase != PhaseQuestion {
		c.sendError("not accepting answers")
		return
	}
	if time.Now().After(room.Deadline) {
		c.sendError("tempo esgotado")
		return
	}
	player, ok := room.Players[c.playerID]
	if !ok {
		c.sendError("player not found")
		return
	}
	if player.answered {
		c.sendError("já respondeu")
		return
	}
	q := room.Questions[room.CurrentIndex]
	qType := q.Type
	if qType == "" {
		qType = models.QuestionMultipleChoice
	}

	elapsed := time.Since(room.QuestionAt)
	correct := false
	points := 0
	similarity := 0.0
	choice := -1
	text := strings.TrimSpace(req.Text)

	switch qType {
	case models.QuestionEssay:
		if text == "" {
			c.sendError("escreva sua resposta")
			return
		}
		if len([]rune(text)) > 2000 {
			c.sendError("resposta muito longa")
			return
		}
		threshold := q.SimilarityThreshold
		if threshold <= 0 {
			threshold = 0.55
		}
		similarity = Similarity(text, q.ExpectedAnswer)
		correct = similarity >= threshold
		if similarity > 0 {
			base := calcScore(q.Weight, q.TimeLimitSec, elapsed)
			points = int(math.Round(float64(base) * similarity))
		}
	default:
		if req.Choice < 0 || req.Choice >= len(q.Options) {
			c.sendError("opção inválida")
			return
		}
		choice = req.Choice
		correct = choice == q.CorrectIndex
		if correct {
			similarity = 1
			points = calcScore(q.Weight, q.TimeLimitSec, elapsed)
		}
	}

	player.answered = true
	player.choice = choice
	player.Score += points
	room.Answers[player.ID] = answerRecord{
		playerID:   player.ID,
		choice:     choice,
		text:       text,
		elapsed:    elapsed,
		correct:    correct,
		points:     points,
		similarity: similarity,
	}
	c.emit("answer_ack", map[string]any{
		"received":   true,
		"points":     points,
		"similarity": similarity,
	})
	room.broadcastLocked("answer_count", map[string]any{
		"answered": countAnswered(room),
		"total":    len(room.Players),
	})
	if countAnswered(room) == len(room.Players) {
		room.revealLocked()
	}
}

func calcScore(weight float64, timeLimitSec int, elapsed time.Duration) int {
	if weight <= 0 {
		weight = 1
	}
	if timeLimitSec <= 0 {
		timeLimitSec = 20
	}
	base := weight * 1000
	ratio := 1 - elapsed.Seconds()/float64(timeLimitSec)
	if ratio < 0 {
		ratio = 0
	}
	// keep at least 20% of base for correct answers within time
	ratio = 0.2 + 0.8*ratio
	return int(math.Round(base * ratio))
}

func countAnswered(room *Room) int {
	n := 0
	for _, p := range room.Players {
		if p.answered {
			n++
		}
	}
	return n
}

func (r *Room) startQuestionLocked() {
	r.cancelTimersLocked()
	r.Answers = make(map[string]answerRecord)
	for _, p := range r.Players {
		p.answered = false
		p.choice = -1
	}
	q := r.Questions[r.CurrentIndex]
	limit := q.TimeLimitSec
	if limit <= 0 {
		limit = 60
	}
	r.Phase = PhaseQuestion
	r.QuestionAt = time.Now()
	r.Deadline = r.QuestionAt.Add(time.Duration(limit) * time.Second)

	r.broadcastLocked("question", map[string]any{
		"question":  publicQuestion(r),
		"startedAt": r.QuestionAt.UTC().Format(time.RFC3339Nano),
		"endsAt":    r.Deadline.UTC().Format(time.RFC3339Nano),
	})

	ctx, cancel := context.WithCancel(context.Background())
	r.timerCancel = cancel
	go r.runTimer(ctx, limit)
}

func publicQuestion(r *Room) models.PublicQuestion {
	q := r.Questions[r.CurrentIndex]
	limit := q.TimeLimitSec
	if limit <= 0 {
		limit = 60
	}
	pub := models.PublicQuestion{
		ID:           q.ID,
		Type:         q.Type,
		Text:         q.Text,
		Options:      q.Options,
		Weight:       q.Weight,
		TimeLimitSec: limit,
		Index:        r.CurrentIndex,
		Total:        len(r.Questions),
	}
	if pub.Type == "" {
		pub.Type = models.QuestionMultipleChoice
	}
	if pub.Type == models.QuestionEssay {
		pub.Options = nil
	}
	return pub
}

func (r *Room) runTimer(ctx context.Context, limit int) {
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.Phase == PhaseQuestion {
			r.revealLocked()
		}
	}
}

func (r *Room) revealLocked() {
	if r.timerCancel != nil {
		r.timerCancel()
		r.timerCancel = nil
	}
	r.Phase = PhaseReveal
	payload := r.buildResultPayloadLocked()
	payload["autoNextInSec"] = 5
	r.broadcastLocked("question_result", payload)
	r.scheduleAutoNextLocked()
}

func (r *Room) buildResultPayloadLocked() map[string]any {
	q := r.Questions[r.CurrentIndex]
	qType := q.Type
	if qType == "" {
		qType = models.QuestionMultipleChoice
	}

	results := make([]map[string]any, 0, len(r.Players))
	for _, p := range r.Players {
		ans, ok := r.Answers[p.ID]
		entry := map[string]any{
			"playerId":   p.ID,
			"nickname":   p.Nickname,
			"choice":     -1,
			"text":       "",
			"correct":    false,
			"points":     0,
			"similarity": 0.0,
			"score":      p.Score,
		}
		if ok {
			entry["choice"] = ans.choice
			entry["text"] = ans.text
			entry["correct"] = ans.correct
			entry["points"] = ans.points
			entry["similarity"] = ans.similarity
		}
		results = append(results, entry)
	}

	payload := map[string]any{
		"type":         qType,
		"correctIndex": q.CorrectIndex,
		"results":      results,
		"leaderboard":  leaderboard(r),
	}
	if qType == models.QuestionEssay {
		payload["expectedAnswer"] = q.ExpectedAnswer
		payload["correctIndex"] = -1
	}
	return payload
}

func (r *Room) scheduleAutoNextLocked() {
	r.cancelAdvanceLocked()
	ctx, cancel := context.WithCancel(context.Background())
	r.advanceCancel = cancel
	index := r.CurrentIndex
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			r.mu.Lock()
			defer r.mu.Unlock()
			if r.Phase != PhaseReveal || r.CurrentIndex != index {
				return
			}
			r.CurrentIndex++
			if r.CurrentIndex >= len(r.Questions) {
				r.finishLocked()
				return
			}
			r.startQuestionLocked()
		}
	}()
}

func (r *Room) cancelAdvanceLocked() {
	if r.advanceCancel != nil {
		r.advanceCancel()
		r.advanceCancel = nil
	}
}

func (r *Room) cancelTimersLocked() {
	if r.timerCancel != nil {
		r.timerCancel()
		r.timerCancel = nil
	}
	r.cancelAdvanceLocked()
}

func (r *Room) finishLocked() {
	r.cancelTimersLocked()
	r.Phase = PhaseFinished
	board := leaderboard(r)
	r.broadcastLocked("finished", map[string]any{
		"leaderboard": board,
		"quizTitle":   r.Quiz.Title,
	})

	ranking := make([]models.RankingEntry, 0, len(board))
	for i, e := range board {
		ranking = append(ranking, models.RankingEntry{
			PlayerID: e["playerId"].(string),
			Nickname: e["nickname"].(string),
			Score:    e["score"].(int),
			Rank:     i + 1,
		})
	}
	rec := &models.SessionRecord{
		ID:         r.SessionID,
		QuizID:     r.Quiz.ID,
		QuizTitle:  r.Quiz.Title,
		TeacherID:  r.TeacherID,
		PIN:        r.PIN,
		Ranking:    ranking,
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
	}
	go func() {
		_ = r.hub.store.SaveSession(context.Background(), rec)
	}()
}

func leaderboard(r *Room) []map[string]any {
	list := make([]*Player, 0, len(r.Players))
	for _, p := range r.Players {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Score == list[j].Score {
			return list[i].Nickname < list[j].Nickname
		}
		return list[i].Score > list[j].Score
	})
	out := make([]map[string]any, 0, len(list))
	for i, p := range list {
		out = append(out, map[string]any{
			"rank":     i + 1,
			"playerId": p.ID,
			"nickname": p.Nickname,
			"avatar":   p.Avatar,
			"score":    p.Score,
		})
	}
	return out
}

func publicPlayers(r *Room) []map[string]any {
	out := make([]map[string]any, 0, len(r.Players))
	for _, p := range r.Players {
		out = append(out, map[string]any{
			"id":       p.ID,
			"nickname": p.Nickname,
			"avatar":   p.Avatar,
			"score":    p.Score,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i]["nickname"].(string) < out[j]["nickname"].(string)
	})
	return out
}

func (r *Room) broadcastLocked(typ string, data any) {
	payload, _ := json.Marshal(data)
	env, _ := json.Marshal(Envelope{Type: typ, Data: payload})
	if r.Host != nil {
		select {
		case r.Host.send <- env:
		default:
		}
	}
	for _, p := range r.Players {
		if p.conn != nil {
			select {
			case p.conn.send <- env:
			default:
			}
		}
	}
}

func (c *Client) cleanup() {
	if c.roomPIN == "" {
		close(c.send)
		return
	}
	room := c.hub.GetRoom(c.roomPIN)
	if room == nil {
		close(c.send)
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if c.role == RoleHost && room.Host == c {
		room.Host = nil
	}
	if c.role == RolePlayer {
		if p, ok := room.Players[c.playerID]; ok && p.conn == c {
			delete(room.Players, c.playerID)
			room.broadcastLocked("player_left", map[string]any{
				"players": publicPlayers(room),
			})
		}
	}
	close(c.send)
}

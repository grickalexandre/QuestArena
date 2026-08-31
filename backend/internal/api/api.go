package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/questarena/questarena/internal/auth"
	"github.com/questarena/questarena/internal/game"
	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/seed"
	"github.com/questarena/questarena/internal/store"
)

type Server struct {
	Store     store.Store
	Auth      auth.Verifier
	Hub       *game.Hub
	AuthMode  string
	StaticDir string
}

func NewServer(s store.Store, v auth.Verifier, hub *game.Hub) *Server {
	return &Server{Store: s, Auth: v, Hub: hub, AuthMode: v.Mode()}
}

func (s *Server) Routes() http.Handler {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/api/health", s.handleHealth).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/auth/mode", s.handleAuthMode).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/auth/dev-login", s.handleDevLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/session", s.handleSession).Methods("POST", "OPTIONS")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.Middleware(s.Auth))

	api.HandleFunc("/quizzes", s.handleListQuizzes).Methods("GET", "OPTIONS")
	api.HandleFunc("/quizzes", s.handleCreateQuiz).Methods("POST", "OPTIONS")
	api.HandleFunc("/quizzes/{id}", s.handleGetQuiz).Methods("GET", "OPTIONS")
	api.HandleFunc("/quizzes/{id}", s.handleUpdateQuiz).Methods("PUT", "OPTIONS")
	api.HandleFunc("/quizzes/{id}", s.handleDeleteQuiz).Methods("DELETE", "OPTIONS")

	api.HandleFunc("/quizzes/{id}/questions", s.handleListQuestions).Methods("GET", "OPTIONS")
	api.HandleFunc("/quizzes/{id}/questions", s.handleCreateQuestion).Methods("POST", "OPTIONS")
	api.HandleFunc("/quizzes/{id}/questions/{qid}", s.handleUpdateQuestion).Methods("PUT", "OPTIONS")
	api.HandleFunc("/quizzes/{id}/questions/{qid}", s.handleDeleteQuestion).Methods("DELETE", "OPTIONS")

	api.HandleFunc("/sessions", s.handleCreateSession).Methods("POST", "OPTIONS")
	api.HandleFunc("/sessions", s.handleListSessions).Methods("GET", "OPTIONS")
	api.HandleFunc("/sessions/{id}", s.handleGetSession).Methods("GET", "OPTIONS")
	api.HandleFunc("/similarity", s.handleSimilarity).Methods("POST", "OPTIONS")

	r.HandleFunc("/ws", s.Hub.ServeWS)

	if s.StaticDir != "" {
		spa := spaHandler{staticDir: s.StaticDir}
		r.PathPrefix("/").Handler(spa)
	}

	return r
}

type spaHandler struct {
	staticDir string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticDir, filepath.Clean(r.URL.Path))
	if !strings.HasPrefix(path, filepath.Clean(h.staticDir)) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		http.ServeFile(w, r, path)
		return
	}
	http.ServeFile(w, r, filepath.Join(h.staticDir, "index.html"))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"authMode": s.AuthMode,
		"lanIP":    lanIPv4(),
	})
}

func lanIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		ip := ipnet.IP.To4()
		if ip == nil || (ip[0] == 169 && ip[1] == 254) {
			continue
		}
		return ip.String()
	}
	return ""
}

func (s *Server) handleAuthMode(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"mode": s.AuthMode})
}

func (s *Server) handleDevLogin(w http.ResponseWriter, r *http.Request) {
	if s.AuthMode != "dev" {
		writeErr(w, http.StatusBadRequest, "dev login only available in DEV_MODE")
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	token, claims, err := s.Auth.DevLogin(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"teacher": map[string]string{
			"id":    claims.ID,
			"email": claims.Email,
			"name":  claims.Name,
		},
	})
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		h := r.Header.Get("Authorization")
		if strings.HasPrefix(h, "Bearer ") {
			req.Token = strings.TrimPrefix(h, "Bearer ")
		}
	}
	if req.Token == "" {
		writeErr(w, http.StatusBadRequest, "token required")
		return
	}
	claims, err := s.Auth.Verify(r.Context(), req.Token)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"teacher": map[string]string{
			"id":    claims.ID,
			"email": claims.Email,
			"name":  claims.Name,
		},
		"authMode": s.AuthMode,
	})
}

func (s *Server) handleListQuizzes(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	if err := seed.EnsureHerancaQuiz(r.Context(), s.Store, teacher.ID); err != nil {
		log.Printf("seed heranca quiz: %v", err)
	}
	if err := seed.EnsureNodeSupabaseQuiz(r.Context(), s.Store, teacher.ID); err != nil {
		log.Printf("seed node+supabase quiz: %v", err)
	}
	if err := seed.EnsureMerQuiz(r.Context(), s.Store, teacher.ID); err != nil {
		log.Printf("seed mer quiz: %v", err)
	}
	list, err := s.Store.ListQuizzes(r.Context(), teacher.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleCreateQuiz(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeErr(w, http.StatusBadRequest, "title required")
		return
	}
	q := &models.Quiz{
		TeacherID:   teacher.ID,
		Title:       req.Title,
		Description: req.Description,
	}
	if err := s.Store.CreateQuiz(r.Context(), q); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, q)
}

func (s *Server) handleGetQuiz(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	id := mux.Vars(r)["id"]
	q, err := s.Store.GetQuiz(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "quiz not found")
		return
	}
	if q.TeacherID != teacher.ID {
		writeErr(w, http.StatusForbidden, "forbidden")
		return
	}
	writeJSON(w, http.StatusOK, q)
}

func (s *Server) handleUpdateQuiz(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	id := mux.Vars(r)["id"]
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	q := &models.Quiz{
		ID:          id,
		TeacherID:   teacher.ID,
		Title:       strings.TrimSpace(req.Title),
		Description: req.Description,
	}
	if q.Title == "" {
		writeErr(w, http.StatusBadRequest, "title required")
		return
	}
	if err := s.Store.UpdateQuiz(r.Context(), q); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := s.Store.GetQuiz(r.Context(), id)
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteQuiz(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	id := mux.Vars(r)["id"]
	if err := s.Store.DeleteQuiz(r.Context(), teacher.ID, id); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleListQuestions(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	quizID := mux.Vars(r)["id"]
	if err := s.ensureQuizOwner(r, teacher.ID, quizID); err != nil {
		writeErr(w, http.StatusForbidden, err.Error())
		return
	}
	list, err := s.Store.ListQuestions(r.Context(), quizID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleCreateQuestion(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	quizID := mux.Vars(r)["id"]
	if err := s.ensureQuizOwner(r, teacher.ID, quizID); err != nil {
		writeErr(w, http.StatusForbidden, err.Error())
		return
	}
	var req models.Question
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := validateQuestion(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	req.QuizID = quizID
	existing, _ := s.Store.ListQuestions(r.Context(), quizID)
	if req.Order == 0 {
		req.Order = len(existing)
	}
	if err := s.Store.CreateQuestion(r.Context(), &req); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, req)
}

func (s *Server) handleUpdateQuestion(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	vars := mux.Vars(r)
	quizID := vars["id"]
	qid := vars["qid"]
	if err := s.ensureQuizOwner(r, teacher.ID, quizID); err != nil {
		writeErr(w, http.StatusForbidden, err.Error())
		return
	}
	var req models.Question
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := validateQuestion(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = qid
	req.QuizID = quizID
	if err := s.Store.UpdateQuestion(r.Context(), &req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func (s *Server) handleDeleteQuestion(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	vars := mux.Vars(r)
	quizID := vars["id"]
	qid := vars["qid"]
	if err := s.ensureQuizOwner(r, teacher.ID, quizID); err != nil {
		writeErr(w, http.StatusForbidden, err.Error())
		return
	}
	if err := s.Store.DeleteQuestion(r.Context(), quizID, qid); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	var req struct {
		QuizID string `json:"quizId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.QuizID == "" {
		writeErr(w, http.StatusBadRequest, "quizId required")
		return
	}
	quiz, err := s.Store.GetQuiz(r.Context(), req.QuizID)
	if err != nil {
		writeErr(w, http.StatusNotFound, "quiz not found")
		return
	}
	if quiz.TeacherID != teacher.ID {
		writeErr(w, http.StatusForbidden, "forbidden")
		return
	}
	questions, err := s.Store.ListQuestions(r.Context(), quiz.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(questions) == 0 {
		writeErr(w, http.StatusBadRequest, "adicione pelo menos uma questão")
		return
	}
	room, err := s.Hub.CreateRoom(*quiz, questions, teacher.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"pin":       room.PIN,
		"sessionId": room.SessionID,
		"quizId":    quiz.ID,
		"quizTitle": quiz.Title,
	})
}

func (s *Server) handleListSessions(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	list, err := s.Store.ListSessions(r.Context(), teacher.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleGetSession(w http.ResponseWriter, r *http.Request) {
	teacher := auth.FromContext(r.Context())
	id := mux.Vars(r)["id"]
	rec, err := s.Store.GetSession(r.Context(), teacher.ID, id)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			writeErr(w, http.StatusForbidden, "forbidden")
			return
		}
		writeErr(w, http.StatusNotFound, "session not found")
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (s *Server) handleSimilarity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text            string   `json:"text"`
		ExpectedAnswer  string   `json:"expectedAnswer"`
		ExpectedAnswers []string `json:"expectedAnswers"`
		Threshold       float64  `json:"threshold"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	q := models.Question{
		ExpectedAnswer:  req.ExpectedAnswer,
		ExpectedAnswers: req.ExpectedAnswers,
	}
	refs := q.EssayReferences()
	if len(refs) == 0 {
		writeErr(w, http.StatusBadRequest, "expectedAnswer required")
		return
	}
	th := req.Threshold
	if th <= 0 {
		th = 0.55
	}
	if th > 1 {
		th = 1
	}
	matches := make([]map[string]any, 0, len(refs))
	best := 0.0
	for _, ref := range refs {
		sim := game.Similarity(req.Text, ref)
		if sim > best {
			best = sim
		}
		matches = append(matches, map[string]any{
			"answer":     ref,
			"similarity": sim,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"similarity": best,
		"passed":     best >= th,
		"threshold":  th,
		"matches":    matches,
	})
}

func (s *Server) ensureQuizOwner(r *http.Request, teacherID, quizID string) error {
	q, err := s.Store.GetQuiz(r.Context(), quizID)
	if err != nil {
		return err
	}
	if q.TeacherID != teacherID {
		return errForbidden
	}
	return nil
}

var errForbidden = &simpleError{"forbidden"}

type simpleError struct{ msg string }

func (e *simpleError) Error() string { return e.msg }

func validateQuestion(q *models.Question) error {
	q.Text = strings.TrimSpace(q.Text)
	if q.Text == "" {
		return &simpleError{"text required"}
	}
	if q.Type == "" {
		q.Type = models.QuestionMultipleChoice
	}
	q.CodeSnippet = strings.TrimRight(strings.Trim(q.CodeSnippet, "\r\n"), " \t\r\n")
	if q.CodeSnippet == "" {
		q.CodeLanguage = ""
	} else {
		if len([]rune(q.CodeSnippet)) > 6000 {
			return &simpleError{"código muito longo (máximo 6000 caracteres)"}
		}
		q.CodeLanguage = models.NormalizeCodeLanguage(q.CodeLanguage)
	}
	switch q.Type {
	case models.QuestionEssay:
		q.NormalizeEssayRefs()
		if q.ExpectedAnswer == "" {
			return &simpleError{"expectedAnswer required for essay questions"}
		}
		if len(q.EssayReferences()) > 8 {
			return &simpleError{"no máximo 8 respostas aceitas"}
		}
		q.Options = nil
		q.CorrectIndex = -1
		if q.SimilarityThreshold <= 0 {
			q.SimilarityThreshold = 0.55
		}
		if q.SimilarityThreshold > 1 {
			q.SimilarityThreshold = 1
		}
		if q.TimeLimitSec <= 0 {
			q.TimeLimitSec = 120
		}
	case models.QuestionMultipleChoice:
		if len(q.Options) < 2 || len(q.Options) > 4 {
			return &simpleError{"provide 2 to 4 options"}
		}
		for i := range q.Options {
			q.Options[i] = strings.TrimSpace(q.Options[i])
			if q.Options[i] == "" {
				return &simpleError{"options cannot be empty"}
			}
		}
		if q.CorrectIndex < 0 || q.CorrectIndex >= len(q.Options) {
			return &simpleError{"correctIndex out of range"}
		}
		q.ExpectedAnswer = ""
		q.ExpectedAnswers = nil
		if q.TimeLimitSec <= 0 {
			q.TimeLimitSec = 60
		}
	default:
		return &simpleError{"type must be multiple_choice or essay"}
	}
	if q.Weight <= 0 {
		q.Weight = 1
	}
	if q.TimeLimitSec > 600 {
		q.TimeLimitSec = 600
	}
	return nil
}

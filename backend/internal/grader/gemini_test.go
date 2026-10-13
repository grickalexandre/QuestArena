package grader

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func fakeGemini(t *testing.T, status *int32, reply string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-goog-api-key") != "k" {
			t.Errorf("missing api key header")
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "RESPOSTA-MODELO") {
			t.Errorf("prompt without reference answer")
		}
		if s := atomic.AddInt32(status, -1); s >= 0 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		out, _ := json.Marshal(map[string]any{
			"candidates": []any{map[string]any{
				"content": map[string]any{"parts": []any{map[string]string{"text": reply}}},
			}},
		})
		_, _ = w.Write(out)
	}))
}

func TestGradeParsesVerdicts(t *testing.T) {
	var fails int32 = 1
	srv := fakeGemini(t, &fails, `{"notas":[{"id":"a1","nota":85,"comentario":"Boa."},{"id":"zz","nota":100,"comentario":"x"},{"id":"a2","nota":140,"comentario":"Ok."}]}`)
	defer srv.Close()

	g := &Gemini{APIKey: "k", Model: "m", BaseURL: srv.URL, Client: srv.Client()}
	got, err := g.Grade(context.Background(), Batch{
		Question: "Diferencie especialização de autorrelacionamento.",
		Expected: []string{"Especialização é tipo de; autorrelacionamento liga a entidade a ela mesma."},
		Answers:  []Answer{{ID: "a1", Text: "..."}, {ID: "a2", Text: "..."}},
	})
	if err != nil {
		t.Fatalf("grade: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 verdicts (unknown id dropped), got %+v", got)
	}
	if got[0].ID != "a1" || got[0].Score != 0.85 || got[0].Feedback != "Boa." {
		t.Errorf("a1: %+v", got[0])
	}
	if got[1].Score != 1 {
		t.Errorf("score must be capped at 1, got %v", got[1].Score)
	}
}

func TestGradeFallsBackWhenPrimaryOverloaded(t *testing.T) {
	var primaryCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/models/primary:") {
			atomic.AddInt32(&primaryCalls, 1)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		out, _ := json.Marshal(map[string]any{
			"candidates": []any{map[string]any{
				"content": map[string]any{"parts": []any{map[string]string{"text": `{"notas":[{"id":"a1","nota":70,"comentario":"Ok."}]}`}}},
			}},
		})
		_, _ = w.Write(out)
	}))
	defer srv.Close()

	g := &Gemini{APIKey: "k", Model: "primary", Fallback: "backup", BaseURL: srv.URL, Client: srv.Client()}
	got, err := g.Grade(context.Background(), Batch{
		Question: "Q", Expected: []string{"E"}, Answers: []Answer{{ID: "a1", Text: "x"}},
	})
	if err != nil || len(got) != 1 || got[0].Score != 0.7 {
		t.Fatalf("fallback: %+v %v", got, err)
	}
	if primaryCalls != attemptsPerModel {
		t.Fatalf("primary tried %d times, want %d", primaryCalls, attemptsPerModel)
	}
}

func TestGradeEmptyBatch(t *testing.T) {
	g := &Gemini{APIKey: "k"}
	got, err := g.Grade(context.Background(), Batch{})
	if err != nil || got != nil {
		t.Fatalf("empty batch: %v %v", got, err)
	}
}

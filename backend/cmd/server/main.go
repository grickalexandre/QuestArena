package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"github.com/questarena/questarena/internal/api"
	"github.com/questarena/questarena/internal/auth"
	"github.com/questarena/questarena/internal/game"
	"github.com/questarena/questarena/internal/store"
	"google.golang.org/api/option"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	port := envOr("PORT", "8080")
	devMode := os.Getenv("DEV_MODE") == "true" || os.Getenv("FIREBASE_CREDENTIALS") == ""

	ctx := context.Background()
	var (
		st store.Store
		v  auth.Verifier
	)

	if devMode {
		dataDir := envOr("DATA_DIR", "data")
		storePath := filepath.Join(dataDir, "store.json")
		log.Printf("starting in DEV_MODE (file store %s + local teacher auth)", storePath)
		mem := store.NewPersistentMemoryStore(storePath)
		st = mem
		v = auth.NewDevVerifier(mem)
	} else {
		credPath := os.Getenv("FIREBASE_CREDENTIALS")
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		opt := option.WithCredentialsFile(credPath)
		conf := &firebase.Config{ProjectID: projectID}
		app, err := firebase.NewApp(ctx, conf, opt)
		if err != nil {
			log.Fatalf("firebase init: %v", err)
		}
		fs, err := store.NewFirestoreStore(ctx, app)
		if err != nil {
			log.Fatalf("firestore: %v", err)
		}
		st = fs
		fv, err := auth.NewFirebaseVerifier(ctx, app, fs)
		if err != nil {
			log.Fatalf("firebase auth: %v", err)
		}
		v = fv
		log.Println("starting with Firebase Auth + Firestore")
	}

	hub := game.NewHub(st)
	srv := api.NewServer(st, v, hub)

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		for _, candidate := range []string{"../frontend/dist", "frontend/dist", "./static"} {
			if st, err := os.Stat(candidate); err == nil && st.IsDir() {
				staticDir = candidate
				break
			}
		}
	}
	if staticDir != "" {
		srv.StaticDir = staticDir
		log.Printf("serving frontend from %s", staticDir)
	}

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("QuestArena listening on :%s (auth=%s)", port, v.Mode())
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

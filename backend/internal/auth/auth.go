package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

type ContextKey string

const UserContextKey ContextKey = "teacher"

type TeacherClaims struct {
	ID    string
	Email string
	Name  string
}

type Verifier interface {
	Verify(ctx context.Context, token string) (*TeacherClaims, error)
	DevLogin(ctx context.Context, email, password, name string) (token string, claims *TeacherClaims, err error)
	Mode() string
}

type FirebaseVerifier struct {
	client *auth.Client
	store  store.Store
}

func NewFirebaseVerifier(ctx context.Context, app *firebase.App, s store.Store) (*FirebaseVerifier, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &FirebaseVerifier{client: client, store: s}, nil
}

func (f *FirebaseVerifier) Mode() string { return "firebase" }

func (f *FirebaseVerifier) Verify(ctx context.Context, token string) (*TeacherClaims, error) {
	tok, err := f.client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}
	email, _ := tok.Claims["email"].(string)
	name, _ := tok.Claims["name"].(string)
	if name == "" {
		name = strings.Split(email, "@")[0]
	}
	claims := &TeacherClaims{ID: tok.UID, Email: email, Name: name}
	_ = f.store.UpsertTeacher(ctx, &models.Teacher{
		ID:        claims.ID,
		Email:     claims.Email,
		Name:      claims.Name,
		CreatedAt: time.Now().UTC(),
	})
	return claims, nil
}

func (f *FirebaseVerifier) DevLogin(context.Context, string, string, string) (string, *TeacherClaims, error) {
	return "", nil, fmt.Errorf("dev login disabled in firebase mode")
}

// DevVerifier provides email/password auth for local development without Firebase.
type DevVerifier struct {
	mu      sync.RWMutex
	users   map[string]devUser // email -> user
	tokens  map[string]string  // token -> email
	store   store.Store
}

type devUser struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
}

func NewDevVerifier(s store.Store) *DevVerifier {
	return &DevVerifier{
		users:  make(map[string]devUser),
		tokens: make(map[string]string),
		store:  s,
	}
}

func (d *DevVerifier) Mode() string { return "dev" }

func hashPassword(pw string) string {
	sum := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(sum[:])
}

func stableDevTeacherID(email string) string {
	sum := sha256.Sum256([]byte("questarena-dev-teacher:" + email))
	return hex.EncodeToString(sum[:16])
}

func (d *DevVerifier) DevLogin(ctx context.Context, email, password, name string) (string, *TeacherClaims, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return "", nil, fmt.Errorf("email and password required")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	u, ok := d.users[email]
	if !ok {
		if name == "" {
			name = strings.Split(email, "@")[0]
		}
		u = devUser{
			ID:           stableDevTeacherID(email),
			Email:        email,
			Name:         name,
			PasswordHash: hashPassword(password),
		}
		d.users[email] = u
	} else if u.PasswordHash != hashPassword(password) {
		return "", nil, fmt.Errorf("invalid credentials")
	}
	token := "dev_" + uuid.NewString()
	d.tokens[token] = email
	claims := &TeacherClaims{ID: u.ID, Email: u.Email, Name: u.Name}
	_ = d.store.UpsertTeacher(ctx, &models.Teacher{
		ID:        claims.ID,
		Email:     claims.Email,
		Name:      claims.Name,
		CreatedAt: time.Now().UTC(),
	})
	return token, claims, nil
}

func (d *DevVerifier) Verify(ctx context.Context, token string) (*TeacherClaims, error) {
	d.mu.RLock()
	email, ok := d.tokens[token]
	if !ok {
		d.mu.RUnlock()
		return nil, fmt.Errorf("invalid token")
	}
	u := d.users[email]
	d.mu.RUnlock()
	return &TeacherClaims{ID: u.ID, Email: u.Email, Name: u.Name}, nil
}

func Middleware(v Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
				return
			}
			token := strings.TrimPrefix(h, "Bearer ")
			claims, err := v.Verify(r.Context(), token)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) *TeacherClaims {
	c, _ := ctx.Value(UserContextKey).(*TeacherClaims)
	return c
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

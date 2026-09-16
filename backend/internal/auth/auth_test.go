package auth

import (
	"context"
	"strings"
	"testing"

	"github.com/questarena/questarena/internal/store"
)

func TestIsAllowedTeacherEmail(t *testing.T) {
	if !IsAllowedTeacherEmail("oliveiraalexandre1972@gmail.com") {
		t.Fatal("canonical email must be allowed")
	}
	if !IsAllowedTeacherEmail("  OliveiraAlexandre1972@Gmail.com  ") {
		t.Fatal("email match must be case-insensitive")
	}
	if IsAllowedTeacherEmail("prof@escola.com") {
		t.Fatal("other emails must be rejected")
	}
	if IsAllowedTeacherEmail("") {
		t.Fatal("empty email must be rejected")
	}
}

func TestDevLoginRejectsOtherEmails(t *testing.T) {
	v := NewDevVerifier(store.NewMemoryStore())
	_, _, err := v.DevLogin(context.Background(), "aluno@escola.com", "12345678", "Aluno")
	if err == nil {
		t.Fatal("expected rejection")
	}
	if !strings.Contains(err.Error(), "autorizado") {
		t.Fatalf("got %v", err)
	}
}

func TestDevLoginAllowsAuthorizedTeacher(t *testing.T) {
	v := NewDevVerifier(store.NewMemoryStore())
	token, claims, err := v.DevLogin(context.Background(), AllowedTeacherEmail, "12345678", "Alexandre")
	if err != nil {
		t.Fatalf("authorized login: %v", err)
	}
	if token == "" || claims == nil || claims.Email != AllowedTeacherEmail {
		t.Fatalf("unexpected claims %+v token %q", claims, token)
	}
	got, err := v.Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got.Email != AllowedTeacherEmail {
		t.Fatalf("verify email %q", got.Email)
	}
}

package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndParse(t *testing.T) {
	mgr := NewJWTManager("test-secret-key-123", 1)
	userID := uint64(42)
	role := "user"

	tok, err := mgr.GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if tok == "" {
		t.Fatal("empty token")
	}

	gotID, gotRole, err := mgr.ParseToken(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if gotID != userID {
		t.Errorf("userID: want %d, got %d", userID, gotID)
	}
	if gotRole != role {
		t.Errorf("role: want %q, got %q", role, gotRole)
	}
}

func TestJWTManager_ParseInvalidToken(t *testing.T) {
	mgr := NewJWTManager("secret-1", 1)
	_, _, err := mgr.ParseToken("not.a.real.token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	mgr := NewJWTManager("secret-2", 0) // 0 hours
	tok, err := mgr.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	time.Sleep(1100 * time.Millisecond)
	_, _, err = mgr.ParseToken(tok)
	if err == nil {
		t.Fatal("expected expired token error")
	}
}

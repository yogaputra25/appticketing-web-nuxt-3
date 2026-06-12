package model

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("supersecret123")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if hash == "" {
		t.Fatal("empty hash")
	}

	u := &User{PasswordHash: hash}
	if !u.CheckPassword("supersecret123") {
		t.Error("expected password match")
	}
	if u.CheckPassword("wrongpassword") {
		t.Error("expected mismatch")
	}
}

func TestHashPassword_TooShort(t *testing.T) {
	_, err := HashPassword("short")
	if err == nil {
		t.Fatal("expected error for short password")
	}
}

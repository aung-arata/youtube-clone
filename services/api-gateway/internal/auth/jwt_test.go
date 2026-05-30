package auth

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1, "testuser", "user")
	if err != nil {
		t.Fatalf("GenerateToken returned unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("expected JWT with 3 parts, got %d", len(parts))
	}
}

func TestValidateToken_Valid(t *testing.T) {
	token, err := GenerateToken(42, "testuser", "user")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned unexpected error: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID: got %d, want 42", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("Username: got %s, want testuser", claims.Username)
	}
	if claims.Role != "user" {
		t.Errorf("Role: got %s, want user", claims.Role)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"empty string", ""},
		{"random string", "notavalidtoken"},
		{"malformed jwt", "header.payload.badsignature"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateToken(tt.token)
			if err == nil {
				t.Error("expected error for invalid token, got nil")
			}
		})
	}
}

func TestValidateToken_TamperedSignature(t *testing.T) {
	token, err := GenerateToken(1, "user", "user")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	tampered := token[:len(token)-1] + "X"
	_, err = ValidateToken(tampered)
	if err == nil {
		t.Error("expected error for tampered token, got nil")
	}
}

func TestValidateToken_ExpiryInFuture(t *testing.T) {
	token, err := GenerateToken(5, "future", "user")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt should not be nil")
	}
	if !claims.ExpiresAt.Time.After(time.Now()) {
		t.Error("token should expire in the future")
	}
}

func TestRefreshToken_Valid(t *testing.T) {
	original, err := GenerateToken(7, "refreshme", "admin")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	refreshed, err := RefreshToken(original)
	if err != nil {
		t.Fatalf("RefreshToken error: %v", err)
	}
	if refreshed == "" {
		t.Fatal("RefreshToken returned empty token")
	}

	claims, err := ValidateToken(refreshed)
	if err != nil {
		t.Fatalf("ValidateToken on refreshed token error: %v", err)
	}
	if claims.UserID != 7 {
		t.Errorf("UserID: got %d, want 7", claims.UserID)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	_, err := RefreshToken("not.a.valid.token")
	if err == nil {
		t.Error("expected error when refreshing invalid token, got nil")
	}
}

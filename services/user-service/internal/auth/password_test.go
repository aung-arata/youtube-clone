package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned unexpected error: %v", err)
	}
	if hashed == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hashed == password {
		t.Error("hashed password should not equal plain text password")
	}
}

func TestHashPassword_DifferentHashesForSameInput(t *testing.T) {
	password := "samepassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first HashPassword error: %v", err)
	}
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second HashPassword error: %v", err)
	}
	// bcrypt uses random salt so hashes should differ
	if hash1 == hash2 {
		t.Error("two hashes of the same password should differ due to random salt")
	}
}

func TestComparePasswords_Correct(t *testing.T) {
	password := "correctpassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	if err := ComparePasswords(hashed, password); err != nil {
		t.Errorf("ComparePasswords failed for correct password: %v", err)
	}
}

func TestComparePasswords_Wrong(t *testing.T) {
	password := "correctpassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	err = ComparePasswords(hashed, "wrongpassword")
	if err == nil {
		t.Error("expected error when comparing with wrong password, got nil")
	}
}

func TestComparePasswords_EmptyPassword(t *testing.T) {
	hashed, err := HashPassword("somepassword")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	err = ComparePasswords(hashed, "")
	if err == nil {
		t.Error("expected error when comparing with empty string, got nil")
	}
}

func TestHashPassword_EmptyString(t *testing.T) {
	hashed, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword of empty string returned error: %v", err)
	}
	// Empty string hashes should still verify correctly
	if err := ComparePasswords(hashed, ""); err != nil {
		t.Errorf("ComparePasswords failed for empty password hash: %v", err)
	}
}

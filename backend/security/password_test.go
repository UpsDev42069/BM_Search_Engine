package security_test

import (
	"os"
	"testing"

	"github.com/UpsDev42069/BM_Search_Engine/backend/security"
)

func TestMain(m *testing.M) {
	// Set test env to true
	os.Setenv("TEST_ENV", "true")

	// Set session secret for testing	purposes
	os.Setenv("SESSION_SECRET", "testsecret123")

	code := m.Run()

	// Clean up environment variables
	os.Unsetenv("TEST_ENV")
	os.Unsetenv("SESSION_SECRET")

	os.Exit(code)
}

func TestHashPassword(t *testing.T) {
	password := "SecurePassword123!"
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("HashPassword() = %v, want a hashed password", hashedPassword)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "SecurePassword123!"
	wrongPassword := "WrongPassword123!"
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}

	// Test for correct password
	if !security.CheckPasswordHash(hashedPassword, password) {
		t.Errorf("CheckPasswordHash() = false, want true")
	}

	// Test for wrong password
	if security.CheckPasswordHash(hashedPassword, wrongPassword) {
		t.Errorf("CheckPasswordHash() = true, want false")
	}
}
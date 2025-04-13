package hashx_test

import (
	"github.com/kerim-dauren/user-service/pkg/hashx"
	"testing"
)

func TestArgon2PasswordManager(t *testing.T) {
	hasher := hashx.NewArgon2Hasher()
	checker := hashx.NewArgon2HashChecker()

	t.Run("Hash and Compare success", func(t *testing.T) {
		pass := "securepassword123"

		// Hash the password
		hashedPassword, err := hasher.Hash(pass)
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}

		// Compare the hashed password with the original password
		err = checker.CompareHashAndPassword(hashedPassword, pass)
		if err != nil {
			t.Fatalf("failed to compare password: %v", err)
		}
	})

	t.Run("Compare with wrong password", func(t *testing.T) {
		pass := "securepassword123"
		wrongPassword := "wrongpassword"

		// Hash the password
		hashedPassword, err := hasher.Hash(pass)
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}

		// Compare the hashed password with the wrong password
		err = checker.CompareHashAndPassword(hashedPassword, wrongPassword)
		if err == nil {
			t.Fatalf("failed to compare password: %v", err)
		}
	})

	t.Run("Hash with empty password", func(t *testing.T) {
		emptyPassword := ""

		_, err := hasher.Hash(emptyPassword)
		if err == nil {
			t.Error("expected an error for empty password, but got none")
		}
	})

	t.Run("Invalid hash format", func(t *testing.T) {
		invalidHash := "invalid$hash$format"

		err := checker.CompareHashAndPassword(invalidHash, "password")
		if err == nil {
			t.Error("expected an error for invalid hash format, but got none")
		}
	})
}

func BenchmarkHash(b *testing.B) {
	hasher := hashx.NewArgon2Hasher()
	pass := "securepassword123"

	for i := 0; i < b.N; i++ {
		_, err := hasher.Hash(pass)
		if err != nil {
			b.Fatalf("failed to hash password: %v", err)
		}
	}
}

func BenchmarkCompareHashAndPassword(b *testing.B) {
	hasher := hashx.NewArgon2Hasher()
	checker := hashx.NewArgon2HashChecker()
	pass := "securepassword123"
	hashedPassword, err := hasher.Hash(pass)
	if err != nil {
		b.Fatalf("failed to hash password: %v", err)
	}

	for i := 0; i < b.N; i++ {
		err := checker.CompareHashAndPassword(hashedPassword, pass)
		if err != nil {
			b.Fatalf("failed to compare password: %v", err)
		}
	}
}

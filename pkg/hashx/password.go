package hashx

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"runtime"
	"strings"
)

type argon2Hasher struct{}

func NewArgon2Hasher() Hasher {
	return argon2Hasher{}
}

func (a argon2Hasher) Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Argon2id params
	threads := getMaxCPUCores(8) // Limit the number of threads
	memory := uint32(32 * 1024)  // 32 MB of memory
	time := uint32(1)            // 1 iteration
	keyLen := uint32(32)         // Key length of 32 bytes
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	result := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, time, threads, b64Salt, b64Hash)
	return result, nil
}

type argon2HashChecker struct{}

func NewArgon2HashChecker() Checker {
	return argon2HashChecker{}
}

func (a argon2HashChecker) CompareHashAndPassword(hashedPassword, password string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return fmt.Errorf("invalid hash format")
	}

	// read the params
	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Decode the salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	// Generate a hash with the same parameters
	keyLen := uint32(len(expectedHash))
	calculatedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	if isValid := subtleCompare(expectedHash, calculatedHash); !isValid {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

func getMaxCPUCores(max int) uint8 {
	cores := runtime.NumCPU()
	if cores > max {
		return uint8(max)
	}
	return uint8(cores)
}

package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPasswordFromUser hashes the user's password.
func HashPasswordFromUser(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

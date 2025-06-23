package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString returns a SHA-256 hash of the given string in hexadecimal format.
func HashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

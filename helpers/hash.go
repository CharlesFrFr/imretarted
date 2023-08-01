package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(s string) string {
	shaBytes := sha256.Sum256([]byte(s))
	return hex.EncodeToString(shaBytes[:])
}
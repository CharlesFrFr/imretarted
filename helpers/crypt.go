package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func HashString(s string) string {
	shaBytes := sha256.Sum256([]byte(s))
	return hex.EncodeToString(shaBytes[:])
}

func DecodeBase64(s string) string {
	decoded, _ := base64.StdEncoding.DecodeString(s)
	return string(decoded)
}
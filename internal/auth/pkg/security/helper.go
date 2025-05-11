package security

import (
	"crypto/rand"
	"encoding/base64"
)

// Helper function to generate secure random strings
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

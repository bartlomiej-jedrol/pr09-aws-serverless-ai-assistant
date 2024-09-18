package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateToken generates a random token of the specified length
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	token := base64.URLEncoding.EncodeToString(bytes)[:length]
	return token, nil
}

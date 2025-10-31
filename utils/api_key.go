package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateAPIKey(size int) (string, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func NewApiKey() (string, error) {
	return GenerateAPIKey(32)
}

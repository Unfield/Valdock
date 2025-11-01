package utils

import (
	"crypto/sha256"
	"encoding/hex"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const ID_LENGTH = 24

func GenerateID() (string, error) {
	var alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	id, err := gonanoid.Generate(alphabet, ID_LENGTH)
	return id, err
}

func GenerateValkeyACLHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

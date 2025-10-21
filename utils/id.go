package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

const ID_LENGTH = 24

func GenerateID() (string, error) {
	var alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	id, err := gonanoid.Generate(alphabet, ID_LENGTH)
	return id, err
}

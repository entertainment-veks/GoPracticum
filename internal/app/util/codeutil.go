package util

import (
	"crypto/rand"
	"encoding/hex"
)

const CODE_LENGTH int = 5

func GenerateCode() (string, error) {
	b := make([]byte, CODE_LENGTH)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

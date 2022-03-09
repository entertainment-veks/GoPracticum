package util

import (
	"crypto/rand"
	"encoding/hex"
)

const codeLength int = 5

func GenerateCode() (string, error) {
	b := make([]byte, codeLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

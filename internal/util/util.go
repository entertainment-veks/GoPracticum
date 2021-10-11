package util

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
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

func IsURL(token string) bool {
	if len(token) == 0 {
		return false
	}
	_, err := url.ParseRequestURI(token)
	return err == nil
}

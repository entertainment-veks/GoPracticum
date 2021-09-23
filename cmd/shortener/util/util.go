package util

import (
	"go_practicum/cmd/shortener/repository"
	"math/rand"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateCode() (string, error) {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	generatedCode := string(b)
	//checking is generated code unic
	//this thing does not tested because it's not stable
	repo, err := repository.NewRepository()
	if err != nil {
		return "", err
	}

	url, err := repo.Get(generatedCode)
	if err != nil {
		return "", err
	}

	if len(url) != 0 {
		return GenerateCode()
	}
	return generatedCode, nil
}

func IsURL(token string) bool {
	if len(token) == 0 {
		return false
	}
	_, err := url.ParseRequestURI(token)
	return err == nil
}

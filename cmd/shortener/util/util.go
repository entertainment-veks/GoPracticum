package util

import (
	"GoPracticum/cmd/shortener/repository"
	"math/rand"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateCode() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	generatedCode := string(b)
	//checking is generated code unic
	//this thing does not tested because it's not stable
	url, _ := repository.NewRepository().Get(generatedCode)
	if len(url) != 0 {
		return GenerateCode()
	}
	return generatedCode
}

func IsURL(token string) bool {
	if len(token) == 0 {
		return false
	}
	_, err := url.ParseRequestURI(token)
	return err == nil
}

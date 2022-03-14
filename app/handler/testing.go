package handler

import (
	"go_practicum/app/config"
)

func newConfig() config.Config {
	return config.Config{
		ServerAddress:   ":8080",
		BaseURL:         "http://127.0.0.1:8080",
		FileStoragePath: "file",
		DatabaseURL:     "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable",
	}
}

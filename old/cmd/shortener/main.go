package main

import (
	"net/http"
	"os"

	"go_practicum/internal/app/shortener"
	"go_practicum/internal/repository"

	"github.com/gorilla/mux"
)

type Service struct {
	Repository *repository.Repository
	BaseURL    string
}

func SetupServer() mux.Router {
	config := shortener.NewConfig()
	file, _ := os.OpenFile(config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	repo := repository.NewRepository(file)
	service := Service{
		Repository: repo,
		BaseURL:    config.BaseURL,
	}

	router := mux.NewRouter()

	router.Handle("/{key}", AuthMiddleware(GzipMiddleware(GetHandler(&service)))).Methods(http.MethodGet)
	router.Handle("/", AuthMiddleware(GunzipMiddleware(PostHandler(&service)))).Methods(http.MethodPost)
	router.Handle("/api/shorten", AuthMiddleware(GunzipMiddleware(PostJSONHandler(&service)))).Methods(http.MethodPost)
	//router.Handle("/user/urls", UserURLsHandler(&service)).Methods(http.MethodGet)

	return *router
}

package main

import (
	"net/http"
	"os"

	"GoPracticum/cmd/shortener/handler"
	"GoPracticum/cmd/shortener/repository"

	"github.com/gorilla/mux"
)

func SetupServer() mux.Router {

	service := repository.Service{
		Repository: repository.NewRepository(),
		BaseURL:    os.Getenv("BASE_URL"),
	}

	router := mux.NewRouter()

	router.HandleFunc("/{key}", handler.GetHandler(&service))
	router.HandleFunc("/", handler.PostHandler(&service))
	router.HandleFunc("/api/shorten", handler.PostJSONHandler(&service))

	return *router
}

func main() {
	if len(os.Getenv("SERVER_ADDRESS")) == 0 && len(os.Getenv("BASE_URL")) == 0 {
		os.Setenv("SERVER_ADDRESS", ":8080")
		os.Setenv("BASE_URL", "http://localhost:8080")
	}

	router := SetupServer()
	http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), &router)
}

package main

import (
	"flag"
	"net/http"
	"os"

	"go_practicum/internal/repository"

	"github.com/gorilla/mux"
)

type Service struct {
	Repository *repository.Repository
	BaseURL    string
}

func init() {
	flag.Func("a", "Server address", func(s string) error {
		os.Setenv("SERVER_ADDRESS", s)
		return nil
	})

	flag.Func("b", "Base url", func(s string) error {
		os.Setenv("BASE_URL", s)
		return nil
	})

	flag.Func("f", "File storage path", func(s string) error {
		os.Setenv("FILE_STORAGE_PATH", s)
		return nil
	})
}

func SetupServer() mux.Router {
	file, _ := os.OpenFile(os.Getenv("FILE_STORAGE_PATH"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	repo := repository.NewRepository(file)
	service := Service{
		Repository: repo,
		BaseURL:    os.Getenv("BASE_URL"),
	}

	router := mux.NewRouter()

	router.Handle("/{key}", GzipMiddleware(GetHandler(&service))).Methods(http.MethodGet)
	router.Handle("/", GunzipMiddleware(PostHandler(&service))).Methods(http.MethodPost)
	router.Handle("/api/shorten", GunzipMiddleware(PostJSONHandler(&service))).Methods(http.MethodPost)

	return *router
}

func main() {

	flag.Parse()

	if len(os.Getenv("SERVER_ADDRESS")) == 0 {
		os.Setenv("SERVER_ADDRESS", ":8080")
	}

	if len(os.Getenv("BASE_URL")) == 0 {
		os.Setenv("BASE_URL", "http://localhost:8080")
	}

	if len(os.Getenv("FILE_STORAGE_PATH")) == 0 {
		os.Setenv("FILE_STORAGE_PATH", "file")
	}

	router := SetupServer()
	http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), &router)
}

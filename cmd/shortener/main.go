package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"GoPracticum/cmd/shortener/repository"

	"github.com/gorilla/mux"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateCode() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isURL(token string) bool {
	if len(token) == 0 {
		return false
	}
	_, err := url.ParseRequestURI(token)
	return err == nil
}

type Service struct {
	repository *repository.Repository
}

func SetupServer() mux.Router {
	service := Service{
		repository: repository.NewRepository(),
	}

	router := mux.NewRouter()

	router.HandleFunc("/{key}", service.getHandler)
	router.HandleFunc("/", service.postHandler)

	return *router
}

func (s *Service) getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := s.repository.Get(vars["key"])

	w.Header().Set("Location", code)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
	}

	link := string(body)

	if !isURL(link) {
		http.Error(w, "Invalid link", http.StatusBadRequest)
		return
	}

	code := generateCode()

	s.repository.Set(code, link)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + code))
}

func main() {
	router := SetupServer()
	http.ListenAndServe(":8080", &router)
}

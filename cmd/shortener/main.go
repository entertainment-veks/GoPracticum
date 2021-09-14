package main

import (
	"encoding/json"
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
		repository.NewRepository(),
	}

	router := mux.NewRouter()

	router.HandleFunc("/{key}", service.getHandler)
	router.HandleFunc("/", service.postHandler)
	router.HandleFunc("/api/shorten", service.postJsonHandler)

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

func (s *Service) postJsonHandler(w http.ResponseWriter, r *http.Request) {
	type Url struct {
		url string
	}

	type Result struct {
		result string
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
	}

	link := Url{}

	json.Unmarshal(body, &link)

	if !isURL(link.url) {
		http.Error(w, "Invalid link", http.StatusBadRequest)
		return
	}

	code := generateCode()

	s.repository.Set(code, link.url)

	rawResult := Result{
		"http://localhost:8080/" + code,
	}

	jsonResult, _ := json.Marshal(rawResult)

	w.Header().Add("Content-Type", "json")
	w.Write(jsonResult)
}

func main() {
	router := SetupServer()
	http.ListenAndServe(":8080", &router)
}

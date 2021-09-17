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
	mainURL    string
}

type Url struct {
	Url string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

func SetupServer(mainURL string) mux.Router {
	service := Service{
		repository.NewRepository(),
		mainURL,
	}

	router := mux.NewRouter()

	router.HandleFunc("/{key}", service.getHandler)
	router.HandleFunc("/", service.postHandler)
	router.HandleFunc("/api/shorten", service.postJSONHandler)

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
	w.Write([]byte(s.mainURL + code))
}

func (s *Service) postJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
	}

	link := Url{}

	json.Unmarshal(body, &link)

	if !isURL(link.Url) {
		http.Error(w, "Invalid link", http.StatusBadRequest)
		return
	}

	code := generateCode()

	s.repository.Set(code, link.Url)

	rawResult := Result{
		s.mainURL + code,
	}

	w.Header().Add("Content-Type", "json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rawResult)
}

func main() {
	mainURL := "https://localhost:8080/"
	router := SetupServer(mainURL)
	http.ListenAndServe(":8080", &router)
}

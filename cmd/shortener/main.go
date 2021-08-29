package main

import (
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var database map[string]string

func ChooseEndpoint(w http.ResponseWriter, r *http.Request) {
	if !isCorrectURL(r.FormValue("s")) {
		w.Write([]byte("400"))
	}

	if r.Method == http.MethodPost {
		PostMethod(w, r)
	}
	if r.Method == http.MethodGet {
		GetMethod(w, r)
	}
}

func PostMethod(w http.ResponseWriter, r *http.Request) {
	var link = r.FormValue("s")
	var code = generateCode()

	database[code] = link

}

func GetMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := database[vars["key"]]

	http.Redirect(w, r, link, r.Response.StatusCode)
}

func generateCode() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isCorrectURL(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

func main() {
	http.HandleFunc("/", ChooseEndpoint)
	http.HandleFunc("/{key}", ChooseEndpoint)

	http.ListenAndServe(":8080", nil)
}

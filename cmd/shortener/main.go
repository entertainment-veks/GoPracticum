package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var database map[string]string

func PostMethod(w http.ResponseWriter, r *http.Request) {
	if !isCorrectURL(r.FormValue("s")) {
		w.Write([]byte("400"))
		return
	}

	var link = r.FormValue("s")
	var code = generateCode()

	database[code] = link

}

func GetMethod(w http.ResponseWriter, r *http.Request) {
	if !isCorrectURL(r.FormValue("s")) {
		w.Write([]byte("400"))
		return
	}

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
	router := mux.NewRouter()
	router.HandleFunc("/", PostMethod)
	router.HandleFunc("/{key}", GetMethod)
	fmt.Println(http.ListenAndServe(":8000", router))
}

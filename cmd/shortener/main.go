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
	if !isURL(r.FormValue("s")) {
		w.WriteHeader(400)
		return
	}

	var link = r.FormValue("s")
	var code = generateCode()

	database[code] = link

	w.Write([]byte("http://localhost:8080/" + code))
}

func GetMethod(w http.ResponseWriter, r *http.Request) {
	if !isURL(r.FormValue("s")) {
		w.WriteHeader(400)
		return
	}

	vars := mux.Vars(r)
	link := database[vars["key"]]

	w.Header().Set("Location", link)
	w.WriteHeader(307)
}

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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", PostMethod)
	router.HandleFunc("/{key}", GetMethod)
	fmt.Println(http.ListenAndServe(":8080", router))
}

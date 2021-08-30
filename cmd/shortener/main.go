package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var database map[string]string

func PostMethod(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	link := string(body)

	if !isURL(link) {
		w.WriteHeader(400)
		return
	}

	var code = generateCode()

	database[code] = link

	w.Write([]byte("http://localhost:8080/" + code))
	w.WriteHeader(201)
}

func GetMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := database[vars["key"]]

	w.Header().Set("Location", code)
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

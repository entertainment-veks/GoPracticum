package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

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

func main() {
	database := make(map[string]string)

	router := mux.NewRouter()

	router.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		code := database[vars["key"]]

		w.Header().Set("Location", code)
		w.WriteHeader(307)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp := r.Body

		body, err := ioutil.ReadAll(temp)

		if err != nil {
			panic(err)
		}

		link := string(body)

		if !isURL(link) {
			w.WriteHeader(400)
			return
		}

		var code = generateCode()

		database[code] = link

		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + code))
	})

	http.ListenAndServe(":8080", router)
}

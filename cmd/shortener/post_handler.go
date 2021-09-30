package main

import (
	"go_practicum/internal/util"
	"io/ioutil"
	"net/http"
)

func PostHandler(s *Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		link := string(body)

		if !util.IsURL(link) {
			http.Error(w, "Invalid link", http.StatusBadRequest)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			http.Error(w, "Unable to generate unic-code", http.StatusInternalServerError)
			return
		}

		s.Repository.Set(code, link)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(s.BaseURL + "/" + code))
	}
}

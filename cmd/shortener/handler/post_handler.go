package handler

import (
	"go_practicum/cmd/shortener/repository"
	"go_practicum/cmd/shortener/util"
	"io/ioutil"
	"net/http"
)

func PostHandler(s *repository.Service) func(w http.ResponseWriter, r *http.Request) {
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

		code := util.GenerateCode()

		s.Repository.Set(code, link)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(s.BaseURL + "/" + code))
	}
}

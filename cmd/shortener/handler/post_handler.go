package handler

import (
	"GoPracticum/cmd/shortener/repository"
	"GoPracticum/cmd/shortener/util"
	"io/ioutil"
	"net/http"
)

func PostHandler(s *repository.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
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

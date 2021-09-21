package handler

import (
	"GoPracticum/cmd/shortener/repository"
	"GoPracticum/cmd/shortener/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type URL struct {
	URL string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

func PostJSONHandler(s *repository.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
		}

		link := URL{}

		json.Unmarshal(body, &link)

		if !util.IsURL(link.URL) {
			http.Error(w, "Invalid link", http.StatusBadRequest)
			return
		}

		code := util.GenerateCode()

		s.Repository.Set(code, link.URL)

		rawResult := Result{
			s.BaseURL + "/" + code,
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rawResult)
	}
}
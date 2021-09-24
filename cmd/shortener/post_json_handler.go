package main

import (
	"encoding/json"
	"go_practicum/cmd/shortener/pkg"
	"io/ioutil"
	"net/http"
)

type URL struct {
	URL string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

func PostJSONHandler(s *Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		link := URL{}

		if err := json.Unmarshal(body, &link); err != nil {
			http.Error(w, "Unable to unmarshal link", http.StatusInternalServerError)
			return
		}

		if !pkg.IsURL(link.URL) {
			http.Error(w, "Invalid link", http.StatusBadRequest)
			return
		}

		code, err := pkg.GenerateCode()
		if err != nil {
			http.Error(w, "Unable to generate unic-code", http.StatusInternalServerError)
			return
		}

		s.Repository.Set(code, link.URL)

		rawResult := Result{
			s.BaseURL + "/" + code,
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(rawResult); err != nil {
			http.Error(w, "Unable to encode result", http.StatusInternalServerError)
		} //if it's not end, need to add 'return'
	}
}

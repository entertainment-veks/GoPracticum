package handler

import (
	"net/http"

	"go_practicum/cmd/shortener/repository"

	"github.com/gorilla/mux"
)

func GetHandler(s *repository.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		link, err := s.Repository.Get(vars["key"])
		if err != nil {
			http.Error(w, "Error till finding url", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

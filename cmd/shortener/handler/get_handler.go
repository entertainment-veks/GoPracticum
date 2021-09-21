package handler

import (
	"net/http"

	"GoPracticum/cmd/shortener/repository"

	"github.com/gorilla/mux"
)

func GetHandler(s *repository.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}
		vars := mux.Vars(r)
		code := s.Repository.Get(vars["key"])

		w.Header().Set("Location", code)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

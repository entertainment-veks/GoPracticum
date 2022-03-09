package handler

import (
	"go_practicum/app/store"
	"net/http"
)

func HandlePing(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s != nil {
			respond(w, http.StatusOK, "")
		} else {
			respond(w, http.StatusInternalServerError, "")
		}
	}
}

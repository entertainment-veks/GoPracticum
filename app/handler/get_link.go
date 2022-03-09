package handler

import (
	"github.com/gorilla/mux"
	"go_practicum/app/config"
	"go_practicum/app/store"
	"net/http"
)

func HandleLinkGet(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["key"]
		l, err := s.Link().GetByCode(code)
		if err != nil {
			if err == store.ErrURLDeleted {
				respond(w, http.StatusGone, "Gone")
				return
			}
			respondError(w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Location", l.Link)
		respond(w, http.StatusTemporaryRedirect, "")
	}
}

func HandleGetUserLinks(s store.Store, cfg config.Config) http.HandlerFunc {
	type userLink struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := s.Link().GetAllByUserID(r.Context().Value(userIDContextKey).(string))
		if err != nil {
			if err == store.ErrRecordNotFound {
				respondError(w, http.StatusNoContent, err)
				return
			}
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		userLinks := []userLink{}
		for _, current := range links {
			userLinks = append(userLinks, userLink{
				ShortURL:    cfg.BaseURL + "/" + current.Code,
				OriginalURL: current.Link,
			})
		}

		respondJSON(w, http.StatusCreated, userLinks)
	}
}

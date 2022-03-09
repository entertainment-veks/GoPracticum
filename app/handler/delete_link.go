package handler

import (
	"go_practicum/app/store"
	"io/ioutil"
	"net/http"
	"strings"
)

func HandleLinkDelete(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytedBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		replacer := strings.NewReplacer(
			"[", "",
			"]", "",
			" ", "",
			`"`, "",
		)
		body := replacer.Replace(string(bytedBody))
		req := strings.Split(body, ",")

		go s.Link().DeleteAllByCode(req)

		respond(w, http.StatusAccepted, "Accepted")
	}
}

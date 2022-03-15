package handler

import (
	"go_practicum/app/store"
	"io/ioutil"
	"net/http"
	"strings"
)

// HandleLinkDelete uses for deleting links from system.
//
// Links which needs to be deleted should be passed as array, ex:
// ["my1code", "code1234", "4321code"]
func HandleLinkDelete(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
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
		stringedBody := replacer.Replace(string(body))
		req := strings.Split(stringedBody, ",")

		go s.Link().DeleteAllByCode(req)

		respond(w, http.StatusAccepted, "Accepted")
	}
}

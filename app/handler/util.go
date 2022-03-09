package handler

import (
	"encoding/json"
	"net/http"
)

type ctxKey int8

const (
	userIDCookieKey         = "shortener-userid"
	userIDContextKey ctxKey = iota
)

func respond(w http.ResponseWriter, code int, data string) {
	w.WriteHeader(code)
	if len(data) != 0 {
		w.Write([]byte(data))
	}
}

func respondError(w http.ResponseWriter, code int, err error) {
	respond(w, code, err.Error())
}

func respondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	respond(w, statusCode, "")
}

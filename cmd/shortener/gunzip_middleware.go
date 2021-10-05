package main

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GunzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			result, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Unable to unzip link", http.StatusInternalServerError)
			}
			r.Body = result
		}

		next.ServeHTTP(w, r)
	})
}

package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware uses when client supports gzip compress
//
// To activate it you need to provide "Content-Encoding" header with "application/gzip" value
func GzipMiddleware(next http.Handler) http.Handler {
	type gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "application/gzip") {
			result, err := gzip.NewReader(r.Body)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			r.Body = result
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "application/gzip") {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "application/gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package main

import "net/http"

func UserURLsHandler(s *Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("user-id")
		// if err == http.ErrNoCookie {
		// 	w.WriteHeader(http.StatusNoContent)
		// }
		// cookie = cookie
	})
}

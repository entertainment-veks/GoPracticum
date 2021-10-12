package main

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// key := []byte{102, 25, 178, 191, 42, 4, 28, 4, 214, 239, 191, 229, 133, 53, 59, 100}

		// h := hmac.New(sha256.New, key)
		// cookie, err := r.Cookie("user-id")
		// if err == http.ErrNoCookie {
		// 	//todo r.SetCookie
		// }
		// value := cookie.Value

	})
}

package handler

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

// AuthMiddleware needs for auth users.
//
// It puts cookie with user id
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(userIDCookieKey)

		var newUserID string
		if err != nil {
			if err == http.ErrNoCookie {
				newUserID = uuid.New().String()

				cookie = &http.Cookie{
					Name:  userIDCookieKey,
					Value: newUserID,
				}
				http.SetCookie(w, cookie)
			} else {
				respondError(w, http.StatusInternalServerError, err)
				return
			}
		}

		newUserID = cookie.Value

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(),
				userIDContextKey,
				newUserID,
			),
		))
	})
}

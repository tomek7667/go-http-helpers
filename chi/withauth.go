package chi

import (
	"context"
	"net/http"
	"strings"

	"github.com/tomek7667/go-http-helpers/h"
)

type ContextKeyType string

const UserContextKey ContextKeyType = "user-go-http-helpers"

func GetUser[User any](r *http.Request) *User {
	u, ok := r.Context().Value(UserContextKey).(*User)
	if !ok {
		panic("GetUser errored - Probably trying to get a user in a request context without the WithAuth middleware")
	}
	return u
}

func WithAuth[User any](auther Auther[User]) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authValue := strings.TrimSpace(r.Header.Get("Authorization"))
			if strings.HasPrefix(authValue, "Bearer ") {
				token, _ := strings.CutPrefix(authValue, "Bearer ")
				user, err := auther.GetUserFromToken(r.Context(), token)
				if err != nil {
					h.ResUnauthorized(w)
					return
				}
				r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))
			} else {
				h.ResUnauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

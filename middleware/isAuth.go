package middleware

import (
	"context"
	fk "github.com/andrei-nita/goframework/framework"
	"github.com/andrei-nita/goframework/framework/auth"
	"net/http"
)

func IsAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if cookie := fk.CookieGetSession("_cookie", w, r); cookie != "" {
			session := auth.Session{}
			session.Uuid = cookie
			check, err := session.Check()
			if err != nil || !check {
				handler(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), "cookie", check)
			r = r.WithContext(ctx)
		}
		handler(w, r)
	}
}

func IsAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if cookie := fk.CookieGetSession("_cookie", w, r); cookie != "" {
			session := auth.Session{}
			session.Uuid = cookie
			check, err := session.Check()
			if err != nil || !check {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), "cookie", check)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

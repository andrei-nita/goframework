package middleware

import (
	"fmt"
	fk "github.com/andrei-nita/goframework/framework"
	"net/http"
)

// Secure fallows this guides:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
// https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000>; includeSubDomains; preload")
		if fk.Server.Mode == fk.ModeProd {
			w.Header().Set("Content-Security-Policy",
				fmt.Sprintf("default-src 'self' %[1]s *.%[1]s", fk.Server.Domain))
		}
		//w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

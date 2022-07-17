package framework

import (
	"encoding/base64"
	"net/http"
)

func CookieSetSession(cookieName string, w http.ResponseWriter, uuid string) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    uuid,
		Domain:   Server.Domain,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if Server.UseSSL {
		cookie.Secure = true
	}

	http.SetCookie(w, &cookie)
}

func CookieGetSession(cookieName string, w http.ResponseWriter, r *http.Request) string {
	// get cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func CookieDeleteSession(cookieName string, w http.ResponseWriter, r *http.Request) error {
	// get cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}

	// delete cookie
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	return nil
}

func CookieSetFlash(cookieName string, w http.ResponseWriter, msg string) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    base64.StdEncoding.EncodeToString([]byte(msg)),
		Domain:   Server.Domain,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if Server.UseSSL {
		cookie.Secure = true
	}

	http.SetCookie(w, &cookie)
}

func CookieGetFlash(cookieName string, w http.ResponseWriter, r *http.Request) string {
	// get cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	// decode cookie
	bytes, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return ""
	}
	// delete cookie
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	return string(bytes)
}

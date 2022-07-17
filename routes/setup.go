package routes

import (
	fk "github.com/andrei-nita/goframework/framework"
	"log"
	"net/http"
	"time"
)

func Setup(mux *http.ServeMux) {

	if fk.Server.CacheTempls {
		if err := parseTemplates(); err != nil {
			log.Fatalln(err)
		}
	}

	if fk.Server.Mode == fk.ModeDev {
		fk.BuildSitemap(
			fk.CreateURL(fk.Domain("/"), time.Now().Format("2006-01-02"), fk.Monthly, 1.0),
			fk.CreateURL(fk.Domain("/csrf-safe"), "2022-07-15", fk.Monthly, 0.6),
			fk.CreateURL(fk.Domain("/api"), "2022-07-16", fk.Monthly, 0.5),
		)
	}

	mux.HandleFunc("/robots.txt", robotsText)
	mux.HandleFunc("/sitemap.xml", sitemap)

	mux.HandleFunc("/", home)
	mux.Handle("/csrf-safe", http.HandlerFunc(csrfSafe))
	mux.HandleFunc("/api", api)
	mux.HandleFunc("/__livereload", livereload)

	// authentication
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)

	mux.HandleFunc("/signup", signup)
	//mux.HandleFunc("/signup_account", signupAccount)

	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/users", users)
}

package routes

import "net/http"

// GET

// robotsText is the route to robots.txt file
func robotsText(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "robots.txt")
}

// sitemap is the route to robots.txt file
func sitemap(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "sitemap.xml")
}

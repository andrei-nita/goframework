package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	data, _, next := methodCSRFAuth(w, r, http.MethodGet)
	if next {
		tmplExecute(w, r, "index", data)
	}
}

func csrfSafe(w http.ResponseWriter, r *http.Request) {
	data, _, next := methodAuth(w, r, http.MethodPost)
	if next {
		tmplExecute(w, r, "csrf-safe", data)
	}
}

type post struct {
	User    string
	Threads []string
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := &post{
		User:    "No Name",
		Threads: []string{"first", "second", "third"},
	}
	JSON, err := json.Marshal(p)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = w.Write(JSON)
	if err != nil {
		log.Fatalln(err)
	}
}

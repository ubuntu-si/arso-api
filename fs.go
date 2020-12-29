package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rakyll/statik/fs"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServerHandler(r chi.Router) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticHandler := http.FileServer(statikFS)
	path := "/"

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		staticHandler.ServeHTTP(w, r)
	})
}

package main

import (
	"net/http"
	"os"

	_ "arso/statik" // UI

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// @title ARSO API
// @version 1.0
// @description This is a ARSO JSON API.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServerHandler(r)
	setupRoutes(r)
	http.ListenAndServe(":"+port, r)
}

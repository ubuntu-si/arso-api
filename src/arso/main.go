package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "arso/statik" // UI

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/docgen/raml"
	"github.com/go-chi/render"
	"github.com/rakyll/statik/fs"
	yaml "gopkg.in/yaml.v2"
)

var build = ""

// Potresi returs list of locations with tremor info.
func Potresi(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, r, ARSOPotresi())
}

// Postaje returs list of locations with weather info.
func Postaje(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, ARSOVreme())
}

func setupRoutes(router chi.Router) {

	router.Get(`/potresi.json`, Potresi)
	router.Get(`/postaje.json`, Postaje)
}

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
	setupRoutes(r)
	fileServerHandler(r)

	http.ListenAndServe(":"+port, r)
}
func docs(router chi.Router) {
	ramlDocs := &raml.RAML{
		Title:     "ARSO API",
		BaseUri:   "https://arso.herokuapp.com",
		Version:   "v1.0",
		MediaType: "application/json",
	}

	chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		handlerInfo := docgen.GetFuncInfo(handler)
		resource := &raml.Resource{
			Description: handlerInfo.Comment,
		}

		return ramlDocs.Add(method, route, resource)
	})

	dr, _ := yaml.Marshal(ramlDocs)
	header := []byte("#%RAML 1.0\n---\n")
	doc := append(header, dr...)
	ioutil.WriteFile("api.yaml", doc, 0644)
}

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

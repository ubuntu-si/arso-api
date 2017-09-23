package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/docgen/raml"
	"github.com/go-chi/render"
	yaml "gopkg.in/yaml.v2"
)

var build = ""

// Potresi returs list of locations with tremor info.
func Potresi(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	for _, potres := range ARSOPotresi() {
		list = append(list, &potres)
	}

	render.RenderList(w, r, list)
}

// Postaje returs list of locations with weather info.
func Postaje(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	for _, potres := range ARSOVreme() {
		list = append(list, &potres)
	}

	render.RenderList(w, r, list)
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
	docs(r)
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	FileServer(r, "/static", http.Dir(filesDir))
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
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

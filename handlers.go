package main

import (
	"net/http"

	_ "arso/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/swaggo/http-swagger"
)

// Potresi godoc
// @Summary Show tremor data
// @Produce  json
// @Success 200 {object} model.Potres
// @Router /potresi.json [get]
func Potresi(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, r, ARSOPotresi())
}

// Postaje godoc
// @Summary Show weather info
// @Produce  json
// @Success 200 {object} model.Postaja
// @Router /postaje.json [get]
func Postaje(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, ARSOVreme())
}

func setupRoutes(router chi.Router) {

	router.Get(`/potresi.json`, Potresi)
	router.Get(`/postaje.json`, Postaje)
	router.Get("/swagger/*", httpSwagger.WrapHandler)
}

package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/staticbin"
	"github.com/martini-contrib/throttle"
	"github.com/stretchr/hoard"
	"log"
	"strconv"
	"time"
)

type Potres struct {
	Magnituda float64
	Lat       float64
	Lon       float64
	Datum     string
	Lokacija  string
}

var m *martini.Martini

func ScrapeARSOPotresi() []Potres {

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument("http://www.arso.gov.si/potresi/obvestila%20o%20potresih/aip/"); e != nil {
		log.Fatal(e)
	}
	var potresi []Potres
	doc.Find("#glavna td.vsebina table tr").Each(func(i int, s *goquery.Selection) {
		magnituda := s.Find("td:nth-child(4)").Text()
		if magnituda != "" {
			potres := Potres{}
			potres.Magnituda, _ = strconv.ParseFloat(magnituda, 2)
			potres.Lat, _ = strconv.ParseFloat(s.Find("td:nth-child(2)").Text(), 3)
			potres.Lon, _ = strconv.ParseFloat(s.Find("td:nth-child(3)").Text(), 3)
			potres.Lokacija = s.Find("td:nth-child(6)").Text()
			potres.Datum = s.Find("td:nth-child(1)").Text()
			potresi = append(potresi, potres)
		}
	})

	return potresi
}

func GetArsoPotresi() []Potres {
	return hoard.Get("potresi", func() (interface{}, *hoard.Expiration) {
		obj := ScrapeARSOPotresi()
		return obj, hoard.Expires().AfterMinutes(2)
	}).([]Potres)
}

func main() {
	m := martini.Classic()
	if martini.Env == "production" {
		// run folowing before deploy
		// go get github.com/jteeuwen/go-bindata/...
		// bin/go-bindata static/
		m.Use(staticbin.Static("static", Asset))
	} else {
		m.Use(martini.Static("static"))
	}

	m.Use(render.Renderer())
	limits := throttle.Policy(&throttle.Quota{
		Limit:  100,
		Within: time.Hour,
	})

	// Setup routes

	m.Get(`/potresi.json`, limits, func(r render.Render) {
		r.JSON(200, GetArsoPotresi())
	})

	m.Run()
}

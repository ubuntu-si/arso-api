package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/staticbin"
	"github.com/martini-contrib/throttle"
	"github.com/stretchr/hoard"
	"log"
	"time"
)

type Potres struct {
	Magnituda string
	Lat       string
	Lon       string
	Datum     string
	Lokacija  string
}

var (
	zadnji_request = 0
	m              *martini.Martini
)

func ScrapeARSO() []Potres {

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
			potres.Magnituda = magnituda
			potres.Lokacija = s.Find("td:nth-child(6)").Text()
			potres.Lat = s.Find("td:nth-child(2)").Text()
			potres.Lon = s.Find("td:nth-child(3)").Text()
			potres.Datum = s.Find("td:nth-child(1)").Text()
			potresi = append(potresi, potres)
		}
	})

	return potresi
}

func GetArso() []Potres {
	return hoard.Get("my-key", func() (interface{}, *hoard.Expiration) {
		obj := ScrapeARSO()
		return obj, hoard.Expires().AfterMinutes(2)
	}).([]Potres)
}

func main() {
	m := martini.Classic()
	m.Use(staticbin.Static("static", Asset))
	m.Use(render.Renderer())
	limits := throttle.Policy(&throttle.Quota{
		Limit:  100,
		Within: time.Hour,
	})

	// Setup routes

	m.Get(`/potresi.json`, limits, func(r render.Render) {
		r.JSON(200, GetArso())
	})

	m.Run()
}

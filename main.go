package main

import (
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gorelic"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/staticbin"
	"github.com/martini-contrib/throttle"
	"github.com/stretchr/hoard"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Potres struct {
	Magnituda float64
	Lat       float64
	Lon       float64
	Datum     string
	Lokacija  string
}

type Postaja struct {
	XMLName  *xml.Name `xml:"data" json:",omitempty"`
	Title    string    `xml:"metData>domain_longTitle"`
	Lat      string    `xml:"metData>domain_lat"`
	Lon      string    `xml:"metData>domain_lon"`
	Altitude string    `xml:"metData>domain_altitude"`
	Issued   string    `xml:"metData>tsUpdated_RFC822"`
	Temp     string    `xml:"metData>t"`
	RH       string    `xml:"metData>rh" json:",omitempty"`
	Pressure string    `xml:"metData>p" json:",omitempty"`
	Valid    string    `xml:"metData>tsValid_issued_UTC"`
	URL      string
	Auto     bool
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
func ScrapeARSOVreme() []Postaja {

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument("http://meteo.arso.gov.si/uploads/probase/www/observ/surface/text/sl/observation_si/index.html"); e != nil {
		log.Fatal(e)
	}

	var vreme []Postaja
	doc.Find("td:nth-child(2) > a").Each(func(i int, s *goquery.Selection) {
		url, found := s.Attr("href")
		if found {

			if strings.Contains(url, ".xml") && !strings.Contains(url, "media") {
				url = "http://meteo.arso.gov.si/" + url
				log.Println("Fetch", url)
				response, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				} else {
					defer response.Body.Close()
					contents, _ := ioutil.ReadAll(response.Body)
					var q Postaja
					xml.Unmarshal(contents, &q)
					if q.Title != "" {
						q.URL = url
						q.Auto = strings.Contains(url, "observationAms")
						vreme = append(vreme, q)
					}
					log.Println(q)
				}
			} else {
				log.Println("Skip", url)
			}

		}
	})

	return vreme
}

func GetArsoPotresi() []Potres {
	return hoard.Get("GetArsoPotresi", func() (interface{}, *hoard.Expiration) {
		obj := ScrapeARSOPotresi()
		return obj, hoard.Expires().AfterMinutes(1)
	}).([]Potres)
}

func GetArsoPostaje() []Postaja {
	return hoard.Get("GetArsoPostaje", func() (interface{}, *hoard.Expiration) {
		obj := ScrapeARSOVreme()
		return obj, hoard.Expires().AfterMinutes(30)
	}).([]Postaja)
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
	nr := os.Getenv("NEWRELIC")
	if nr != "" {
		m.Use(gorelic.Handler)
		gorelic.InitNewrelicAgent(nr, "wpapi", true)
	}
	m.Use(render.Renderer())
	limits := throttle.Policy(&throttle.Quota{
		Limit:  120,
		Within: time.Hour,
	})

	// Setup routes

	m.Get(`/potresi.json`, limits, func(r render.Render) {
		r.JSON(200, GetArsoPotresi())
	})

	m.Get(`/postaje.json`, limits, func(r render.Render) {
		r.JSON(200, GetArsoPostaje())
	})

	m.Run()
}

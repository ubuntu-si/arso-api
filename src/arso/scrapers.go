package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	cache "github.com/patrickmn/go-cache"
)

var cacheArso = cache.New(5*time.Minute, 30*time.Second)

// Potres holds info about earthquake
type Potres struct {
	Magnituda float64
	Lat       float64
	Lon       float64
	Datum     string
	Lokacija  string
}

func (p *Potres) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Postaja holds info about weather
type Postaja struct {
	XMLName       *xml.Name `xml:"data" json:",omitempty"`
	ID            string    `xml:"metData>domain_meteosiId"`
	Title         string    `xml:"metData>domain_longTitle"`
	Lat           float64   `xml:"metData>domain_lat"`
	Lon           float64   `xml:"metData>domain_lon"`
	Altitude      float64   `xml:"metData>domain_altitude"`
	Issued        string    `xml:"metData>tsUpdated_RFC822"`
	Temp          float64   `xml:"metData>t"`
	Wind          float64   `xml:"metData>ff_val" json:",omitempty"`
	WindDirection string    `xml:"metData>dd_icon" json:",omitempty"`
	RH            float64   `xml:"metData>rh" json:",omitempty"`
	Pressure      float64   `xml:"metData>p" json:",omitempty"`
	Sky           string    `xml:"metData>nn_shortText" json:",omitempty"`
	Valid         string    `xml:"metData>tsValid_issued_UTC"`
	URL           string
	Auto          bool
}

func (p *Postaja) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ARSOPotresi returs slice of Potres struct
func ARSOPotresi() []Potres {
	var potresi []Potres
	var doc *goquery.Document
	var e error

	if res, found := cacheArso.Get("potresi"); found {
		return res.([]Potres)
	}

	if doc, e = goquery.NewDocument("http://www.arso.gov.si/potresi/obvestila%20o%20potresih/aip/"); e != nil {
		return potresi
	}

	doc.Find("#glavna td.vsebina table tr").Each(func(i int, s *goquery.Selection) {
		magnituda, err := strconv.ParseFloat(s.Find("td:nth-child(4)").Text(), 2)
		if magnituda > 0 && err == nil {
			potres := Potres{}
			potres.Magnituda = magnituda
			potres.Lat, _ = strconv.ParseFloat(s.Find("td:nth-child(2)").Text(), 3)
			potres.Lon, _ = strconv.ParseFloat(s.Find("td:nth-child(3)").Text(), 3)
			potres.Lokacija = s.Find("td:nth-child(6)").Text()
			potres.Datum = s.Find("td:nth-child(1)").Text()
			potresi = append(potresi, potres)
		}
	})
	cacheArso.Set("potresi", potresi, cache.DefaultExpiration)
	return potresi
}

// ARSOVreme returs slice of Postaje struct
func ARSOVreme() []Postaja {
	var vreme []Postaja
	var doc *goquery.Document
	var e error

	if res, found := cacheArso.Get("vreme"); found {
		return res.([]Postaja)
	}

	if doc, e = goquery.NewDocument("http://meteo.arso.gov.si/uploads/probase/www/observ/surface/text/sl/observation_si/index.html"); e != nil {
		return vreme
	}

	doc.Find("td:nth-child(2) > a").Each(func(i int, s *goquery.Selection) {
		url, found := s.Attr("href")
		if found {

			if strings.Contains(url, ".xml") && !strings.Contains(url, "media") && !strings.Contains(url, "_si_") {

				url = "http://meteo.arso.gov.si/" + url
				response, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				} else {
					defer response.Body.Close()
					contents, _ := ioutil.ReadAll(response.Body)
					var q Postaja
					xml.Unmarshal(contents, &q)
					if q.Title != "" {
						q.ID = getMD5Hash(q.ID)
						q.URL = fmt.Sprintf("/vreme/%s", q.ID)
						q.Auto = strings.Contains(url, "observationAms")
						vreme = append(vreme, q)
					}
				}
			} else {
				log.Println("Skip", url)
			}

		}
	})
	cacheArso.Set("vreme", vreme, cache.DefaultExpiration)
	return vreme
}

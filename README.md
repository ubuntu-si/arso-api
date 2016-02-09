[![Build Status](https://travis-ci.org/dz0ny/arso-api.svg?branch=master)](https://travis-ci.org/dz0ny/arso-api)
arso-potresi
============

Arso potresi na zemljevidi in z API-jem. Vir http://www.arso.gov.si/potresi/obvestila%20o%20potresih/aip/

# Gradnja

    env GOPATH=$(pwd) go build arso

#### JSON

    GET http://potresi.herokuapp.com/potresi.json

    [
      {
        "Magnituda": "2.3",
        "Lat": "46.12",
        "Lon": "14.71",
        "Datum": "11/06/2014 11:39:21",
        "Lokacija": "8.km  V od DOMŽAL"
      },
      ...
    ]

    GET http://potresi.herokuapp.com/postaje.json

    [
      {
        "Title": "Letališče Cerklje ob Krki",
        "Lat": 45.8936,
        "Lon": 15.525,
        "Altitude": 154,
        "Issued": "11 Jun 2014 20:22:00 +0000",
        "Temp": 10,
        "RH": 56,
        "Pressure": 0,
        "URL": "/uploads/probase/www/observ/surface/text/sl/observation_si_latest.xml",
        "Auto": false
      },
      ...
    ]

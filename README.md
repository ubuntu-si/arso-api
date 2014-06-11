[![Build Status](https://travis-ci.org/dz0ny/arso-potresi.svg?branch=master)](https://travis-ci.org/dz0ny/arso-potresi)
arso-potresi
============

Arso potresi na zemljevidi in z API-jem. Vir http://www.arso.gov.si/potresi/obvestila%20o%20potresih/aip/

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

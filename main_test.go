package main

import (
	"testing"
)

func TestArso(t *testing.T) {
	potresi := ScrapeARSOPotresi()
	for _, potres := range potresi {
		t.Log("Potres %s - %s %d %d\n", potres.Magnituda, potres.Lokacija, potres.Lat, potres.Lon)
	}

}

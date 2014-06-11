package main

import (
	"testing"
)

func TestArso(t *testing.T) {
	potresi := ScrapeARSO()
	for _, potres := range potresi {
		t.Log("Potres %s - %s\n", potres.Magnituda, potres.Lokacija)
	}

}

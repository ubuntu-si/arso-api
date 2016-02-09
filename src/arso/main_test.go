package main

import (
	"testing"
)

func TestArsoPotresi(t *testing.T) {
	potresi := ARSOPotresi()
	for _, potres := range potresi {
		t.Log(potres.Magnituda, potres.Lokacija)
	}

}
func TestArsoVreme(t *testing.T) {
	vreme := ARSOVreme()
	for _, postaja := range vreme {
		t.Log(postaja.Title, postaja.Temp, postaja.URL)
	}
}
